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

// Create 创建资源
func (r *resource) Create(ctx context.Context, create *ResourceEditModel) (err error) {
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

// Update 修改资源
func (r *resource) Update(ctx context.Context, id model.ID, update *ResourceEditModel) (err error) {
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
		exist, err = dao.RoleResourceWebRoute.ExistByResourceWebRouteID(ctx, resourceWebRouteID)
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
