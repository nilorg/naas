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

// ResourceWebRouteEdit ...
type ResourceWebRouteEdit struct {
	Name             string   `json:"name"`
	Path             string   `json:"path"`
	Method           string   `json:"method"`
	ResourceServerID model.ID `json:"resource_server_id"`
}

// AddWebRoute 添加web路由
func (*resource) AddWebRoute(ctx context.Context, req *ResourceWebRouteEdit) (err error) {
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
	v := &model.ResourceWebRoute{
		Name:             req.Name,
		Path:             req.Path,
		Method:           req.Method,
		ResourceServerID: req.ResourceServerID,
	}
	err = dao.ResourceWebRoute.Insert(ctx, v)
	return
}

// UpdateWebRoute 修改web路由
func (*resource) UpdateWebRoute(ctx context.Context, resourceWebRouteID model.ID, req *ResourceWebRouteEdit) (err error) {
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

	var rwr *model.ResourceWebRoute
	rwr, err = dao.ResourceWebRoute.Select(ctx, resourceWebRouteID)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	rwr.Name = req.Name
	rwr.Path = req.Path
	rwr.Method = req.Method
	rwr.ResourceServerID = req.ResourceServerID
	err = dao.ResourceWebRoute.Update(ctx, rwr)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
	}
	return
}

// DeleteWebRoute 删除web路由
func (*resource) DeleteWebRoute(ctx context.Context, resourceWebRouteIDs ...model.ID) (err error) {
	// 验证Web路由是否被使用，如果被使用不能删除
	for _, resourceWebRouteID := range resourceWebRouteIDs {
		var exist bool
		exist, err = dao.RoleResourceRelation.ExistByRelationTypeAndRelationID(ctx, model.RoleResourceRelationTypeWebRoute, resourceWebRouteID)
		if err != nil {
			logrus.WithContext(ctx).Errorln(err)
			return
		}
		if exist {
			err = errors.New("资源已被使用不能删除")
			return
		}
	}
	err = dao.ResourceWebRoute.DeleteInIDs(ctx, resourceWebRouteIDs...)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
	}
	return
}

// ListWebRouteByResourceID 根据资源服务器ID获取Web路由
func (r *resource) ListWebRouteByResourceID(ctx context.Context, resourceID model.ID, limit int) (list []*model.ResourceWebRoute, err error) {
	return dao.ResourceWebRoute.ListByResourceID(ctx, resourceID, limit)
}

// ListWebRoutePaged ...
func (r *resource) ListWebRoutePaged(ctx context.Context, start, limit int) (result []*model.ResultResourceWebRoute, total int64, err error) {
	var (
		rwrList []*model.ResourceWebRoute
	)
	rwrList, total, err = dao.ResourceWebRoute.ListPaged(ctx, start, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, rwr := range rwrList {
		resultRes := &model.ResultResourceWebRoute{
			ResourceWebRoute: rwr,
		}
		if rwr.ResourceServerID > 0 {
			resultRes.ResourceServer, _ = dao.Resource.Select(ctx, rwr.ResourceServerID)
		}
		result = append(result, resultRes)
	}
	return
}

// GetResourceWebRoute resource web route
func (*resource) GetResourceWebRoute(ctx context.Context, resourceWebRouteID model.ID) (resource *model.ResourceWebRoute, err error) {
	resource, err = dao.ResourceWebRoute.Select(ctx, resourceWebRouteID)
	return
}

// ListWebRoutePagedByRoleCode ...
func (r *resource) ListWebRoutePagedByRoleCode(ctx context.Context, start, limit int, roleCode model.Code) (result []*model.ResourceWebRoute, total int64, err error) {
	var (
		rrwrList []*model.RoleResourceRelation
	)
	rrwrList, total, err = dao.RoleResourceRelation.ListPagedByRelationTypeAndRoleCode(ctx, start, limit, model.RoleResourceRelationTypeWebRoute, roleCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, rrwr := range rrwrList {
		var rwr *model.ResourceWebRoute
		rwr, err = dao.ResourceWebRoute.Select(ctx, rrwr.RelationID)
		if err != nil {
			return
		}
		result = append(result, rwr)
	}
	return
}

// ================

// ResourceWebMenuEdit ...
type ResourceWebMenuEdit struct {
	Name             string   `json:"name"`
	URL              string   `json:"url"`
	Icon             string   `json:"icon"`
	SerialNumber     int      `json:"serial_number"`
	Leaf             int      `json:"leaf"` // 是：子组件，否：是父组件
	ParentID         model.ID `json:"parent_id"`
	ResourceServerID model.ID `json:"resource_server_id"`
}

// AddWebMenu 添加web菜单
func (*resource) AddWebMenu(ctx context.Context, req *ResourceWebMenuEdit) (err error) {
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
		var menu *model.ResourceWebMenu
		menu, err = dao.ResourceWebMenu.Select(ctx, req.ParentID)
		if err != nil {
			return
		}
		if menu.Leaf {
			err = errors.New("叶子菜单，不可作为上级菜单")
			return
		}
		level = menu.Level + 1
	}

	v := &model.ResourceWebMenu{
		Name:             req.Name,
		URL:              req.URL,
		Icon:             req.Icon,
		Level:            level,
		SerialNumber:     req.SerialNumber,
		Leaf:             req.Leaf == 1,
		ParentID:         req.ParentID,
		ResourceServerID: req.ResourceServerID,
	}
	err = dao.ResourceWebMenu.Insert(ctx, v)
	return
}

// UpdateWebMenu 修改web菜单
func (*resource) UpdateWebMenu(ctx context.Context, resourceWebMenuID model.ID, req *ResourceWebMenuEdit) (err error) {
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

	var rwr *model.ResourceWebMenu
	rwr, err = dao.ResourceWebMenu.Select(ctx, resourceWebMenuID)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	level := 0
	if req.ParentID > 0 {
		var menu *model.ResourceWebMenu
		menu, err = dao.ResourceWebMenu.Select(ctx, req.ParentID)
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
		count, err = dao.ResourceWebMenu.CountByParentID(ctx, rwr.ID)
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
	err = dao.ResourceWebMenu.Update(ctx, rwr)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
	}
	return
}

// DeleteWebMenu 删除web菜单
func (*resource) DeleteWebMenu(ctx context.Context, resourceWebMenuIDs ...model.ID) (err error) {
	// 验证Web路由是否被使用，如果被使用不能删除
	for _, resourceWebMenuID := range resourceWebMenuIDs {
		var exist bool
		exist, err = dao.RoleResourceRelation.ExistByRelationTypeAndRelationID(ctx, model.RoleResourceRelationTypeWebMenu, resourceWebMenuID)
		if err != nil {
			logrus.WithContext(ctx).Errorln(err)
			return
		}
		if exist {
			err = errors.New("资源已被使用不能删除")
			return
		}
	}
	err = dao.ResourceWebMenu.DeleteInIDs(ctx, resourceWebMenuIDs...)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
	}
	return
}

// ListWebMenuByResourceID 根据资源服务器ID获取Web路由
func (r *resource) ListWebMenuByResourceID(ctx context.Context, resourceID model.ID, limit int) (list []*model.ResourceWebMenu, err error) {
	return dao.ResourceWebMenu.ListByResourceServerID(ctx, resourceID, limit)
}

// ListWebMenuPaged ...
func (r *resource) ListWebMenuPaged(ctx context.Context, start, limit int) (result []*model.ResultResourceWebMenu, total int64, err error) {
	var (
		rwmList []*model.ResourceWebMenu
	)
	rwmList, total, err = dao.ResourceWebMenu.ListPaged(ctx, start, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, rwm := range rwmList {
		resultRes := &model.ResultResourceWebMenu{
			ResourceWebMenu: rwm,
		}
		if rwm.ResourceServerID > 0 {
			resultRes.ResourceServer, _ = dao.Resource.Select(ctx, rwm.ResourceServerID)
		}
		if rwm.ParentID > 0 {
			resultRes.ParentResourceWebMenu, _ = dao.ResourceWebMenu.Select(ctx, rwm.ParentID)
		}
		result = append(result, resultRes)
	}
	return
}

// GetResourceWebMenu resource web menu
func (*resource) GetResourceWebMenu(ctx context.Context, resourceWebMenuID model.ID) (resource *model.ResourceWebMenu, err error) {
	resource, err = dao.ResourceWebMenu.Select(ctx, resourceWebMenuID)
	return
}

// ListWebMenuPagedByRoleCode ...
func (r *resource) ListWebMenuPagedByRoleCode(ctx context.Context, start, limit int, roleCode model.Code) (result []*model.ResourceWebMenu, total int64, err error) {
	var (
		rrwrList []*model.RoleResourceRelation
	)
	rrwrList, total, err = dao.RoleResourceRelation.ListPagedByRelationTypeAndRoleCode(ctx, start, limit, model.RoleResourceRelationTypeWebMenu, roleCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, rrwr := range rrwrList {
		var rwr *model.ResourceWebMenu
		rwr, err = dao.ResourceWebMenu.Select(ctx, rrwr.RelationID)
		if err != nil {
			return
		}
		result = append(result, rwr)
	}
	return
}

func (r *resource) RecursiveResourceWebMenu(ctx context.Context, resourceServerID model.ID) []*model.ResourceWebMenu {
	if resourceServerID <= 0 {
		return nil
	}
	var (
		rootResourceWebMenus []*model.ResourceWebMenu
		err                  error
	)
	rootResourceWebMenus, err = dao.ResourceWebMenu.ListRootByResourceServerID(ctx, resourceServerID)
	if err != nil {
		logrus.Errorf("dao.ResourceWebMenu.ListRootByResourceServerID(%v): %s", resourceServerID, err)
	}
	r.recursiveResourceWebMenu(ctx, rootResourceWebMenus)
	return rootResourceWebMenus
}

func (r *resource) recursiveResourceWebMenu(ctx context.Context, webMenus []*model.ResourceWebMenu) {
	if len(webMenus) == 0 {
		return
	}
	var (
		childWebMenus []*model.ResourceWebMenu
		err           error
	)
	for _, webMenu := range webMenus {
		childWebMenus, err = dao.ResourceWebMenu.ListByResourceServerIDAndParentID(ctx, webMenu.ResourceServerID, webMenu.ID, -1)
		if err != nil {
			logrus.Errorf("dao.ResourceWebMenu.ListByResourceServerIDAndParentID(%v, %v, -1): %s", webMenu.ResourceServerID, webMenu.ID, err)
		}
		r.recursiveResourceWebMenu(ctx, childWebMenus)
		webMenu.ChildResourceWebMenus = childWebMenus
	}
}
