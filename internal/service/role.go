package service

import (
	"context"
	"fmt"

	"github.com/nilorg/naas/pkg/errors"

	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/casbin"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type role struct {
}

func (r *role) Recursive(ctx context.Context, organizationID model.ID) []*model.Role {
	var (
		rootRoles []*model.Role
		err       error
	)
	if organizationID == 0 {
		rootRoles, err = dao.Role.SelectByRoot(ctx)
		if err != nil {
			logrus.Errorf("dao.Role.SelectByRoot: %s", err)
		}
	} else {
		rootRoles, err = dao.Role.SelectByRootAndOrganizationID(ctx, organizationID)
		if err != nil {
			logrus.Errorf("dao.Role.SelectByRootAndOrganizationID(%v): %s", organizationID, err)
		}
	}
	r.recursive(ctx, rootRoles)
	return rootRoles
}

func (r *role) recursive(ctx context.Context, roles []*model.Role) {
	if len(roles) == 0 {
		return
	}
	var (
		childRoles []*model.Role
		err        error
	)
	for _, role := range roles {
		childRoles, err = dao.Role.SelectByParentCode(ctx, role.Code)
		if err != nil {
			logrus.Errorf("dao.Role.SelectByParentCode(%v): %s", role.Code, err)
		}
		r.recursive(ctx, childRoles)
		role.ChildRoles = childRoles
	}
}

// GetAllRoleByUserID ...
func (r *role) GetAllRoleByUserID(ctx context.Context, userID model.ID) (roles []*model.UserRole, err error) {
	roles, err = dao.UserRole.SelectAllByUserID(ctx, userID)
	return
}

// ListResourceRouteByRoleCode 根据RoleCode获取资源服务器路由
func (r *role) ListResourceRouteByRoleCode(ctx context.Context, roleCode model.Code, limit int) (list []*model.ResourceRoute, err error) {
	var rrwrList []*model.RoleResourceRelation
	rrwrList, err = dao.RoleResourceRelation.ListByRelationTypeAndRoleCode(ctx, model.RoleResourceRelationTypeRoute, roleCode, limit)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	var ids []model.ID
	for _, rrwr := range rrwrList {
		ids = append(ids, rrwr.RelationID)
	}
	list, err = dao.ResourceRoute.ListInIDs(ctx, ids...)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
	}
	return
}

// GetOneByCode 根据code获取角色
func (r *role) GetOneByCode(ctx context.Context, code model.Code) (role *model.Role, err error) {
	role, err = dao.Role.SelectByCode(ctx, code)
	return
}

// ListByName 根据名称查询列表
func (r *role) ListByName(ctx context.Context, name string, limit int) (list []*model.Role, err error) {
	return dao.Role.ListByNameAndOrganizationID(ctx, 0, name, limit)
}

// ListByNameForOrganization 根据名称查询列表
func (r *role) ListByNameForOrganization(ctx context.Context, organizationID model.ID, name string, limit int) (list []*model.Role, err error) {
	return dao.Role.ListByNameAndOrganizationID(ctx, organizationID, name, limit)
}

func (r *role) ListPaged(ctx context.Context, start, limit int) (result []*model.ResultRole, total int64, err error) {
	var (
		roleList []*model.Role
	)
	roleList, total, err = dao.Role.ListPaged(ctx, start, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, role := range roleList {
		resultRole := &model.ResultRole{
			Role: role,
		}
		if role.ParentCode != "" {
			resultRole.ParentRole, _ = dao.Role.SelectByCode(ctx, role.ParentCode)
		}
		if role.OrganizationID > 0 {
			resultRole.Organization, _ = dao.Organization.Select(ctx, role.OrganizationID)
		}
		result = append(result, resultRole)
	}
	return
}

// RoleEditModel ...
type RoleEditModel struct {
	Code           model.Code `json:"code"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	ParentCode     model.Code `json:"parent_code"`
	OrganizationID model.ID   `json:"organization_id"`
}

// Create 创建权限
func (r *role) Create(ctx context.Context, create *RoleEditModel) (err error) {
	if create.ParentCode != "" {
		var role *model.Role
		role, err = dao.Role.SelectByCode(ctx, create.ParentCode)
		if err != nil {
			return
		}
		if role.OrganizationID != create.OrganizationID {
			err = errors.New("上级权限组织和创建组织不相同")
			return
		}
	}

	var organization *model.Organization
	organization, err = dao.Organization.Select(ctx, create.OrganizationID)
	if err != nil {
		return
	}
	createCode := model.Code(fmt.Sprintf("%s_%s", organization.Code, create.Code))
	var exist bool
	exist, err = dao.Role.ExistByCode(ctx, createCode)
	if err != nil {
		return
	}
	if exist {
		err = errors.ErrRoleCodeExist
		return
	}
	role := &model.Role{
		Name:           create.Name,
		Description:    create.Description,
		ParentCode:     create.ParentCode,
		OrganizationID: create.OrganizationID,
	}
	role.Code = createCode
	err = dao.Role.Insert(ctx, role)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
	}
	return
}

// Update 修改权限
func (r *role) Update(ctx context.Context, update *RoleEditModel) (err error) {
	var (
		role *model.Role
	)
	role, err = dao.Role.SelectByCode(ctx, update.Code)
	if err != nil {
		return
	}
	if update.ParentCode != "" {
		var parentRole *model.Role
		parentRole, err = dao.Role.SelectByCode(ctx, update.ParentCode)
		if err != nil {
			return
		}
		if role.OrganizationID != parentRole.OrganizationID {
			err = errors.New("上级权限组织和修改组织不相同")
			return
		}
	}

	role.Name = update.Name
	role.Description = update.Description
	role.ParentCode = update.ParentCode

	err = dao.Role.Update(ctx, role)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		err = errors.ErrRoleUpdate
	}
	return
}

// DeleteByCodes 根据Code删除角色
func (r *role) DeleteByCodes(ctx context.Context, codes ...model.Code) (err error) {
	err = dao.Role.DeleteInCodes(ctx, codes...)
	return
}

// ListByUserIDAndOrganizationID ...
func (r *role) ListByUserIDAndOrganizationID(ctx context.Context, userID model.ID, organizationID model.ID) (list []*model.Role, err error) {
	var userRoles []*model.UserRole
	userRoles, err = dao.UserRole.ListByUserIDAndOrganizationID(ctx, userID, organizationID)
	if err != nil {
		return
	}
	for _, userRole := range userRoles {
		var role *model.Role
		role, err = dao.Role.SelectByCode(ctx, userRole.RoleCode)
		if err != nil {
			return
		}
		list = append(list, role)
	}
	return
}

// AddRoleResourceRelation 添加角色资源关系
func (r *role) AddRoleResourceRelation(
	ctx context.Context,
	roleCode model.Code,
	relationType model.RoleResourceRelationType,
	resourceServerID model.ID,
	relationIDs ...model.ID,
) (err error) {
	tran := store.DB.Begin()
	ctx = contexts.NewGormTranContext(ctx, tran)
	defer func() {
		if err != nil {
			tran.Rollback()
		}
	}()

	// 查询原有的
	var ids []model.ID
	ids, err = dao.RoleResourceRelation.ListForRelationIDByRelationTypeAndRoleCodeAndResourceServerID(ctx, relationType, roleCode, resourceServerID)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	added, deleted := model.DiffIDSlice(ids, relationIDs)
	for _, relationID := range deleted {
		err = dao.RoleResourceRelation.DeleteByRelationTypeAndRoleCodeAndRelationID(ctx, relationType, roleCode, relationID)
		if err != nil {
			logrus.WithContext(ctx).Errorln(err)
			return
		}

		if relationType == model.RoleResourceRelationTypeRoute {
			var resourceRoute *model.ResourceRoute
			resourceRoute, err = dao.ResourceRoute.Select(ctx, relationID)
			if err != nil {
				logrus.WithContext(ctx).Errorln(err)
				return
			}
			sub, dom, obj, act := formatRoutePolicyForRole(roleCode, resourceRoute)
			_, casbinErr := casbin.Enforcer.DeletePermission(sub, dom, obj, act)
			if casbinErr != nil {
				logrus.Errorf("casbin.Enforcer.DeletePermission: %s", casbinErr)
			}
		}
	}
	for _, relationID := range added {
		err = dao.RoleResourceRelation.Insert(ctx, &model.RoleResourceRelation{
			RoleCode:         roleCode,
			RelationID:       relationID,
			ResourceServerID: resourceServerID,
			RelationType:     relationType,
		})
		if err != nil {
			return
		}
		if relationType == model.RoleResourceRelationTypeRoute {
			var resourceRoute *model.ResourceRoute
			resourceRoute, err = dao.ResourceRoute.Select(ctx, relationID)
			if err != nil {
				logrus.WithContext(ctx).Errorln(err)
				return
			}
			if resourceRoute.ResourceServerID != resourceServerID {
				err = errors.New("路由的资源服务器ID不匹配")
				return
			}
			sub, dom, obj, act := formatRoutePolicyForRole(roleCode, resourceRoute)
			_, casbinErr := casbin.Enforcer.AddPolicy(sub, dom, obj, act)
			if casbinErr != nil {
				logrus.Errorf("casbin.Enforcer.AddPolicy: %s", casbinErr)
			}
		}
	}
	err = tran.Commit().Error
	return
}
