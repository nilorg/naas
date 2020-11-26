package service

import (
	"context"

	"github.com/nilorg/naas/pkg/errors"
	"github.com/sirupsen/logrus"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"gorm.io/gorm"
)

type resource struct {
}

// ResourceEditModel ...
type ResourceEditModel struct {
	Name           string   `json:"name"`
	Secret         string   `json:"secret"`
	Description    string   `json:"description"`
	OrganizationID model.ID `json:"organization_id"`
}

// CreateServer 创建资源
func (r *resource) CreateServer(ctx context.Context, create *ResourceEditModel) (err error) {
	if create.OrganizationID > 0 {
		var existOrgID bool
		existOrgID, err = dao.Organization.ExistByID(ctx, create.OrganizationID)
		if err != nil {
			return
		}
		if !existOrgID {
			err = errors.ErrOrganizationNotFound
			return
		}
	}
	res := &model.Resource{
		Name:           create.Name,
		Secret:         create.Secret,
		Description:    create.Description,
		OrganizationID: create.OrganizationID,
	}
	err = dao.Resource.Insert(ctx, res)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		err = errors.ErrResourceCreate
	}
	return
}

// UpdateServer 修改资源
func (r *resource) UpdateServer(ctx context.Context, id model.ID, update *ResourceEditModel) (err error) {
	var (
		res *model.Resource
	)
	res, err = dao.Resource.Select(ctx, id)
	if err != nil {
		return
	}
	if update.OrganizationID > 0 {
		var existOrgID bool
		existOrgID, err = dao.Organization.ExistByID(ctx, update.OrganizationID)
		if err != nil {
			return
		}
		if !existOrgID {
			err = errors.ErrOrganizationNotFound
			return
		}
	}
	res.Name = update.Name
	res.Description = update.Description
	res.Secret = update.Secret
	res.OrganizationID = update.OrganizationID

	err = dao.Resource.Update(ctx, res)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		err = errors.ErrResourceUpdate
		return
	}
	return
}

// GetServer resource
func (*resource) GetServer(ctx context.Context, id model.ID) (resource *model.Resource, err error) {
	resource, err = dao.Resource.Select(ctx, id)
	return
}

// LoadPolicy 加载规则
func (*resource) LoadPolicy(ctx context.Context, resourceID model.ID) (results []*gormadapter.CasbinRule, err error) {
	return dao.Resource.LoadPolicy(ctx, resourceID)
}

// DeleteByIDs 根据ID删除资源服务器
func (r *resource) DeleteByIDs(ctx context.Context, ids ...model.ID) (err error) {
	err = dao.Resource.DeleteInIDs(ctx, ids)
	return
}

// ListByName 根据名称查询列表
func (r *resource) ListByName(ctx context.Context, name string, limit int) (list []*model.Resource, err error) {
	return dao.Resource.ListByName(ctx, name, limit)
}

// ListServerPaged ...
func (r *resource) ListServerPaged(ctx context.Context, start, limit int) (result []*model.ResultResourceServer, total int64, err error) {
	var (
		resList []*model.Resource
	)
	resList, total, err = dao.Resource.ListPaged(ctx, start, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, res := range resList {
		resultRes := &model.ResultResourceServer{
			ResourceServer: res,
		}
		if res.OrganizationID != 0 {
			resultRes.Organization, _ = dao.Organization.Select(ctx, res.OrganizationID)
		}
		result = append(result, resultRes)
	}
	return
}

// ResourceRouteEdit ...
type ResourceRouteEdit struct {
	Name             string   `json:"name"`
	Path             string   `json:"path"`
	Method           string   `json:"method"`
	ResourceServerID model.ID `json:"resource_server_id"`
}

// AddRoute 添加路由
func (*resource) AddRoute(ctx context.Context, req *ResourceRouteEdit) (err error) {
	var resourceExist bool
	resourceExist, err = dao.Resource.ExistByID(ctx, req.ResourceServerID)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	if !resourceExist {
		err = errors.ErrResourceNotFound
		return
	}
	v := &model.ResourceRoute{
		Name:             req.Name,
		Path:             req.Path,
		Method:           req.Method,
		ResourceServerID: req.ResourceServerID,
	}
	err = dao.ResourceRoute.Insert(ctx, v)
	return
}

// UpdateRoute 修改路由
func (*resource) UpdateRoute(ctx context.Context, resourceRouteID model.ID, req *ResourceRouteEdit) (err error) {
	var resourceExist bool
	resourceExist, err = dao.Resource.ExistByID(ctx, req.ResourceServerID)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	if !resourceExist {
		err = errors.ErrResourceNotFound
		return
	}

	var rwr *model.ResourceRoute
	rwr, err = dao.ResourceRoute.Select(ctx, resourceRouteID)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	rwr.Name = req.Name
	rwr.Path = req.Path
	rwr.Method = req.Method
	rwr.ResourceServerID = req.ResourceServerID
	err = dao.ResourceRoute.Update(ctx, rwr)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
	}
	return
}

// DeleteRoute 删除路由
func (*resource) DeleteRoute(ctx context.Context, resourceRouteIDs ...model.ID) (err error) {
	// 验证路由是否被使用，如果被使用不能删除
	for _, resourceRouteID := range resourceRouteIDs {
		var exist bool
		exist, err = dao.RoleResourceRelation.ExistByRelationTypeAndRelationID(ctx, model.RoleResourceRelationTypeRoute, resourceRouteID)
		if err != nil {
			logrus.WithContext(ctx).Errorln(err)
			return
		}
		if exist {
			err = errors.New("资源已被使用不能删除")
			return
		}
	}
	err = dao.ResourceRoute.DeleteInIDs(ctx, resourceRouteIDs...)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
	}
	return
}

// ListRouteByResourceID 根据资源服务器ID获取路由
func (r *resource) ListRouteByResourceID(ctx context.Context, resourceID model.ID, limit int) (list []*model.ResourceRoute, err error) {
	return dao.ResourceRoute.ListByResourceID(ctx, resourceID, limit)
}

// ListRoutePaged ...
func (r *resource) ListRoutePaged(ctx context.Context, start, limit int) (result []*model.ResultResourceRoute, total int64, err error) {
	var (
		rwrList []*model.ResourceRoute
	)
	rwrList, total, err = dao.ResourceRoute.ListPaged(ctx, start, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, rwr := range rwrList {
		resultRes := &model.ResultResourceRoute{
			ResourceRoute: rwr,
		}
		if rwr.ResourceServerID > 0 {
			resultRes.ResourceServer, _ = dao.Resource.Select(ctx, rwr.ResourceServerID)
		}
		result = append(result, resultRes)
	}
	return
}

// GetResourceRoute resource  route
func (*resource) GetResourceRoute(ctx context.Context, resourceRouteID model.ID) (resource *model.ResourceRoute, err error) {
	resource, err = dao.ResourceRoute.Select(ctx, resourceRouteID)
	return
}

// ListRoutePagedByRoleCode ...
func (r *resource) ListRoutePagedByRoleCode(ctx context.Context, start, limit int, roleCode model.Code) (result []*model.ResourceRoute, total int64, err error) {
	var (
		rrwrList []*model.RoleResourceRelation
	)
	rrwrList, total, err = dao.RoleResourceRelation.ListPagedByRelationTypeAndRoleCode(ctx, start, limit, model.RoleResourceRelationTypeRoute, roleCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, rrwr := range rrwrList {
		var rwr *model.ResourceRoute
		rwr, err = dao.ResourceRoute.Select(ctx, rrwr.RelationID)
		if err != nil {
			return
		}
		result = append(result, rwr)
	}
	return
}

// ================

// ResourceMenuEdit ...
type ResourceMenuEdit struct {
	Name             string   `json:"name"`
	URL              string   `json:"url"`
	Icon             string   `json:"icon"`
	SerialNumber     int      `json:"serial_number"`
	Leaf             int      `json:"leaf"` // 是：子组件，否：是父组件
	ParentID         model.ID `json:"parent_id"`
	ResourceServerID model.ID `json:"resource_server_id"`
}

// AddMenu 添加菜单
func (*resource) AddMenu(ctx context.Context, req *ResourceMenuEdit) (err error) {
	var resourceExist bool
	resourceExist, err = dao.Resource.ExistByID(ctx, req.ResourceServerID)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	if !resourceExist {
		err = errors.ErrResourceNotFound
		return
	}
	level := 0
	if req.ParentID > 0 {
		var menu *model.ResourceMenu
		menu, err = dao.ResourceMenu.Select(ctx, req.ParentID)
		if err != nil {
			return
		}
		if menu.Leaf {
			err = errors.New("叶子菜单，不可作为上级菜单")
			return
		}
		level = menu.Level + 1
	}

	v := &model.ResourceMenu{
		Name:             req.Name,
		URL:              req.URL,
		Icon:             req.Icon,
		Level:            level,
		SerialNumber:     req.SerialNumber,
		Leaf:             req.Leaf == 1,
		ParentID:         req.ParentID,
		ResourceServerID: req.ResourceServerID,
	}
	err = dao.ResourceMenu.Insert(ctx, v)
	return
}

// UpdateMenu 修改菜单
func (*resource) UpdateMenu(ctx context.Context, resourceMenuID model.ID, req *ResourceMenuEdit) (err error) {
	var resourceExist bool
	resourceExist, err = dao.Resource.ExistByID(ctx, req.ResourceServerID)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	if !resourceExist {
		err = errors.ErrResourceNotFound
		return
	}

	var rwr *model.ResourceMenu
	rwr, err = dao.ResourceMenu.Select(ctx, resourceMenuID)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	level := 0
	if req.ParentID > 0 {
		var menu *model.ResourceMenu
		menu, err = dao.ResourceMenu.Select(ctx, req.ParentID)
		if err != nil {
			return
		}
		if menu.Leaf {
			err = errors.New("叶子菜单，不可作为上级菜单")
			return
		}
		level = menu.Level + 1
	}

	leaf := req.Leaf == 1
	if leaf { // 检查是否有子菜单，如果有提示用户不能设置为叶子菜单
		var count int64
		count, err = dao.ResourceMenu.CountByParentID(ctx, rwr.ID)
		if err != nil {
			return
		}
		if count > 0 {
			err = errors.New("叶子菜单不能未做父级菜单")
			return
		}
	}

	rwr.Name = req.Name
	rwr.URL = req.URL
	rwr.Icon = req.Icon
	rwr.Level = level
	rwr.SerialNumber = req.SerialNumber
	rwr.Leaf = leaf
	rwr.ParentID = req.ParentID
	rwr.ResourceServerID = req.ResourceServerID
	err = dao.ResourceMenu.Update(ctx, rwr)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
	}
	return
}

// DeleteMenu 删除菜单
func (*resource) DeleteMenu(ctx context.Context, resourceMenuIDs ...model.ID) (err error) {
	// 验证路由是否被使用，如果被使用不能删除
	for _, resourceMenuID := range resourceMenuIDs {
		var exist bool
		exist, err = dao.RoleResourceRelation.ExistByRelationTypeAndRelationID(ctx, model.RoleResourceRelationTypeMenu, resourceMenuID)
		if err != nil {
			logrus.WithContext(ctx).Errorln(err)
			return
		}
		if exist {
			err = errors.New("资源已被使用不能删除")
			return
		}
	}
	err = dao.ResourceMenu.DeleteInIDs(ctx, resourceMenuIDs...)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
	}
	return
}

// ListMenuByResourceID 根据资源服务器ID获取路由
func (r *resource) ListMenuByResourceID(ctx context.Context, resourceID model.ID, limit int) (list []*model.ResourceMenu, err error) {
	return dao.ResourceMenu.ListByResourceServerID(ctx, resourceID, limit)
}

// ListMenuPaged ...
func (r *resource) ListMenuPaged(ctx context.Context, start, limit int) (result []*model.ResultResourceMenu, total int64, err error) {
	var (
		rwmList []*model.ResourceMenu
	)
	rwmList, total, err = dao.ResourceMenu.ListPaged(ctx, start, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, rwm := range rwmList {
		resultRes := &model.ResultResourceMenu{
			ResourceMenu: rwm,
		}
		if rwm.ResourceServerID > 0 {
			resultRes.ResourceServer, _ = dao.Resource.Select(ctx, rwm.ResourceServerID)
		}
		if rwm.ParentID > 0 {
			resultRes.ParentResourceMenu, _ = dao.ResourceMenu.Select(ctx, rwm.ParentID)
		}
		result = append(result, resultRes)
	}
	return
}

// GetResourceMenu resource  menu
func (*resource) GetResourceMenu(ctx context.Context, resourceMenuID model.ID) (resource *model.ResourceMenu, err error) {
	resource, err = dao.ResourceMenu.Select(ctx, resourceMenuID)
	return
}

// ListMenuPagedByRoleCode ...
func (r *resource) ListMenuPagedByRoleCode(ctx context.Context, start, limit int, roleCode model.Code) (result []*model.ResourceMenu, total int64, err error) {
	var (
		rrwrList []*model.RoleResourceRelation
	)
	rrwrList, total, err = dao.RoleResourceRelation.ListPagedByRelationTypeAndRoleCode(ctx, start, limit, model.RoleResourceRelationTypeMenu, roleCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, rrwr := range rrwrList {
		var rwr *model.ResourceMenu
		rwr, err = dao.ResourceMenu.Select(ctx, rrwr.RelationID)
		if err != nil {
			return
		}
		result = append(result, rwr)
	}
	return
}

func (r *resource) RecursiveResourceMenu(ctx context.Context, resourceServerID model.ID) []*model.ResourceMenu {
	if resourceServerID <= 0 {
		return nil
	}
	var (
		rootResourceMenus []*model.ResourceMenu
		err               error
	)
	rootResourceMenus, err = dao.ResourceMenu.ListRootByResourceServerID(ctx, resourceServerID)
	if err != nil {
		logrus.Errorf("dao.ResourceMenu.ListRootByResourceServerID(%v): %s", resourceServerID, err)
	}
	r.recursiveResourceMenu(ctx, rootResourceMenus)
	return rootResourceMenus
}

func (r *resource) recursiveResourceMenu(ctx context.Context, menus []*model.ResourceMenu) {
	if len(menus) == 0 {
		return
	}
	var (
		childMenus []*model.ResourceMenu
		err        error
	)
	for _, menu := range menus {
		childMenus, err = dao.ResourceMenu.ListByResourceServerIDAndParentID(ctx, menu.ResourceServerID, menu.ID, -1)
		if err != nil {
			logrus.Errorf("dao.ResourceMenu.ListByResourceServerIDAndParentID(%v, %v, -1): %s", menu.ResourceServerID, menu.ID, err)
		}
		r.recursiveResourceMenu(ctx, childMenus)
		menu.ChildResourceMenus = childMenus
	}
}

// ================

// ResourceActionEdit ...
type ResourceActionEdit struct {
	Code             model.Code `json:"code"`
	Name             string     `json:"name"`
	Group            string     `json:"group"`
	Description      string     `json:"description"`
	ResourceServerID model.ID   `json:"resource_server_id"`
}

// AddAction 添加动作
func (*resource) AddAction(ctx context.Context, req *ResourceActionEdit) (err error) {
	var resourceExist bool
	resourceExist, err = dao.Resource.ExistByID(ctx, req.ResourceServerID)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	if !resourceExist {
		err = errors.ErrResourceNotFound
		return
	}
	codeExist := false
	codeExist, err = dao.ResourceAction.ExistByCode(ctx, req.Code)
	if err != nil {
		return
	}
	if codeExist {
		err = errors.New("资源动作Code已存在")
		return
	}

	v := &model.ResourceAction{
		Code:             req.Code,
		Name:             req.Name,
		Group:            req.Group,
		Description:      req.Description,
		ResourceServerID: req.ResourceServerID,
	}
	err = dao.ResourceAction.Insert(ctx, v)
	return
}

// UpdateAction 修改动作
func (*resource) UpdateAction(ctx context.Context, resourceActionID model.ID, req *ResourceActionEdit) (err error) {
	var resourceExist bool
	resourceExist, err = dao.Resource.ExistByID(ctx, req.ResourceServerID)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	if !resourceExist {
		err = errors.ErrResourceNotFound
		return
	}

	var action *model.ResourceAction
	action, err = dao.ResourceAction.Select(ctx, resourceActionID)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	if action.Code != req.Code {
		a, aerr := dao.ResourceAction.SelectByCode(ctx, req.Code)
		if aerr != nil {
			if !errors.Is(aerr, gorm.ErrRecordNotFound) {
				err = aerr
				return
			}
		} else {
			if a.ID != action.ID {
				err = errors.New("动作Code已被使用")
				return
			}
		}
	}

	action.Code = req.Code
	action.Name = req.Name
	action.Group = req.Group
	action.Description = req.Description
	action.ResourceServerID = req.ResourceServerID
	err = dao.ResourceAction.Update(ctx, action)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
	}
	return
}

// DeleteAction 删除动作
func (*resource) DeleteAction(ctx context.Context, resourceActionIDs ...model.ID) (err error) {
	// 验证动作是否被使用，如果被使用不能删除
	for _, resourceActionID := range resourceActionIDs {
		var exist bool
		exist, err = dao.RoleResourceRelation.ExistByRelationTypeAndRelationID(ctx, model.RoleResourceRelationTypeAction, resourceActionID)
		if err != nil {
			logrus.WithContext(ctx).Errorln(err)
			return
		}
		if exist {
			err = errors.New("资源已被使用不能删除")
			return
		}
	}
	err = dao.ResourceAction.DeleteInIDs(ctx, resourceActionIDs...)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
	}
	return
}

// ListActionByResourceID 根据资源服务器ID获取动作
func (r *resource) ListActionByResourceID(ctx context.Context, resourceID model.ID, limit int) (list []*model.ResourceAction, err error) {
	return dao.ResourceAction.ListByResourceServerID(ctx, resourceID, limit)
}

// ListActionPaged ...
func (r *resource) ListActionPaged(ctx context.Context, start, limit int) (list []*model.ResultResourceAction, total int64, err error) {
	var actionList []*model.ResourceAction
	actionList, total, err = dao.ResourceAction.ListPaged(ctx, start, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
	}
	for _, action := range actionList {
		resultAction := &model.ResultResourceAction{
			ResourceAction: action,
		}
		if action.ResourceServerID > 0 {
			resultAction.ResourceServer, _ = dao.Resource.Select(ctx, action.ResourceServerID)
		}
		list = append(list, resultAction)
	}
	return
}

// GetResourceAction resource action
func (*resource) GetResourceAction(ctx context.Context, resourceActionID model.ID) (resource *model.ResourceAction, err error) {
	resource, err = dao.ResourceAction.Select(ctx, resourceActionID)
	return
}

// ListActionPagedByRoleCode ...
func (r *resource) ListActionPagedByRoleCode(ctx context.Context, start, limit int, roleCode model.Code) (result []*model.ResourceAction, total int64, err error) {
	var (
		relationList []*model.RoleResourceRelation
	)
	relationList, total, err = dao.RoleResourceRelation.ListPagedByRelationTypeAndRoleCode(ctx, start, limit, model.RoleResourceRelationTypeAction, roleCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, relation := range relationList {
		var action *model.ResourceAction
		action, err = dao.ResourceAction.Select(ctx, relation.RelationID)
		if err != nil {
			return
		}
		result = append(result, action)
	}
	return
}
