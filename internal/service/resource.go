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

// Get resource
func (*resource) Get(ctx context.Context, id model.ID) (resource *model.Resource, err error) {
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

// ListPaged ...
func (r *resource) ListPaged(ctx context.Context, start, limit int) (result []*model.ResultResource, total int64, err error) {
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
		resultRes := &model.ResultResource{
			Resource: res,
		}
		if res.OrganizationID != 0 {
			resultRes.Organization, _ = dao.Organization.Select(ctx, res.OrganizationID)
		}
		result = append(result, resultRes)
	}
	return
}

// ResourceAddWebRouteRequest ...
type ResourceAddWebRouteRequest struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

// AddWebRoute 添加web路由
func (*resource) AddWebRoute(ctx context.Context, resourceID model.ID, req *ResourceAddWebRouteRequest) (err error) {
	v := &model.ResourceWebRoute{
		Name:       req.Name,
		Path:       req.Path,
		Method:     req.Method,
		ResourceID: resourceID,
	}
	err = dao.ResourceWebRoute.Insert(ctx, v)
	return
}
