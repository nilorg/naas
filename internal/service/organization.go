package service

import (
	"context"

	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type organization struct {
}

// OrganizationEditModel ...
type OrganizationEditModel struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Code        model.Code `json:"code"`
	ParentID    model.ID   `json:"parent_id"`
}

// Create 创建组织
func (o *organization) Create(ctx context.Context, create *OrganizationEditModel) (err error) {
	var exist bool
	exist, err = dao.Organization.ExistByCode(ctx, create.Code)
	if err != nil {
		return
	}
	if exist {
		err = errors.ErrOrganizationCodeExist
		return
	}
	if create.ParentID > 0 {
		var existParentID bool
		existParentID, err = dao.Organization.ExistByID(ctx, create.ParentID)
		if err != nil {
			return
		}
		if !existParentID {
			err = errors.ErrOrganizationParentNotExist
			return
		}
	}
	org := &model.Organization{
		Name:        create.Name,
		Description: create.Description,
		Code:        create.Code,
		ParentID:    create.ParentID,
	}
	err = dao.Organization.Insert(ctx, org)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		err = errors.ErrOrganizationCreate
	}
	return
}

// Update 修改组织
func (o *organization) Update(ctx context.Context, id model.ID, update *OrganizationEditModel) (err error) {
	var (
		org *model.Organization
	)
	org, err = dao.Organization.Select(ctx, id)
	if err != nil {
		return
	}
	if update.ParentID > 0 {
		if update.ParentID == id {
			err = errors.ErrOrganizationCurrentAndParentEqual
			return
		}
		var existParentID bool
		existParentID, err = dao.Organization.ExistByID(ctx, update.ParentID)
		if err != nil {
			return
		}
		if !existParentID {
			err = errors.ErrOrganizationParentNotExist
			return
		}
	}
	org.Name = update.Name
	org.Description = update.Description
	org.ParentID = update.ParentID

	err = dao.Organization.Update(ctx, org)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		err = errors.ErrOrganizationUpdate
		return
	}
	return
}

func (o *organization) ListPaged(ctx context.Context, start, limit int) (result []*model.ResultOrganization, total int64, err error) {
	var (
		orgList []*model.Organization
	)
	orgList, total, err = dao.Organization.ListPaged(ctx, start, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, org := range orgList {
		resultOrg := &model.ResultOrganization{
			Organization: org,
		}
		if org.ParentID != 0 {
			resultOrg.ParentOrganization, _ = dao.Organization.Select(ctx, org.ParentID)
		}
		result = append(result, resultOrg)
	}
	return
}

// ListByName 根据名称查询列表
func (o *organization) ListByName(ctx context.Context, name string, limit int) (list []*model.Organization, err error) {
	return dao.Organization.ListByName(ctx, name, limit)
}

// DeleteByID 根据ID删除用户
func (o *organization) DeleteByIDs(ctx context.Context, ids ...model.ID) (err error) {
	tran := store.DB.Begin()
	ctx = contexts.NewGormTranContext(ctx, tran)
	err = dao.Organization.DeleteInIDs(ctx, ids...)
	if err != nil {
		tran.Rollback()
		return
	}
	// TODO：删除组织关联的用户，权限
	err = tran.Commit().Error
	return
}

// GetOneByID 根据ID获取用户
func (o *organization) GetOneByID(ctx context.Context, id model.ID) (org *model.Organization, err error) {
	org, err = dao.Organization.Select(ctx, id)
	return
}

// ListByUserID 根据用户ID查询列表
func (o *organization) ListByUserID(ctx context.Context, userID model.ID) (list []*model.Organization, err error) {
	var userOrganization []*model.UserOrganization
	userOrganization, err = dao.UserOrganization.SelectAllByUserID(ctx, userID)
	if err != nil {
		return
	}
	for _, uo := range userOrganization {
		var org *model.Organization
		org, err = dao.Organization.Select(ctx, uo.OrganizationID)
		if err != nil {
			return
		}
		list = append(list, org)
	}
	return
}
