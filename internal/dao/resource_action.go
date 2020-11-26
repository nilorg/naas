package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"gorm.io/gorm"
)

// ResourceActioner ...
type ResourceActioner interface {
	Insert(ctx context.Context, m *model.ResourceAction) (err error)
	Delete(ctx context.Context, code model.Code) (err error)
	DeleteInIDs(ctx context.Context, ids ...model.ID) (err error)
	Select(ctx context.Context, id model.ID) (m *model.ResourceAction, err error)
	SelectByCode(ctx context.Context, code model.Code) (m *model.ResourceAction, err error)
	Update(ctx context.Context, m *model.ResourceAction) (err error)
	ListByResourceServerID(ctx context.Context, resourceServerID model.ID, limit int) (list []*model.ResourceAction, err error)
	ListPaged(ctx context.Context, start, limit int) (list []*model.ResourceAction, total int64, err error)
	ListPagedByResourceServerID(ctx context.Context, resourceServerID model.ID, start, limit int) (list []*model.ResourceAction, total int64, err error)
	ListByNameAndResourceServerID(ctx context.Context, resourceServerID model.ID, name string, limit int) (list []*model.ResourceAction, err error)
	ExistByCode(ctx context.Context, code model.Code) (exist bool, err error)
}

type resourceAction struct {
}

func (r *resourceAction) Insert(ctx context.Context, m *model.ResourceAction) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(m).Error
	return
}

func (r *resourceAction) Delete(ctx context.Context, code model.Code) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.ResourceAction{}, "code = ?", code).Error
	return
}

func (r *resourceAction) Select(ctx context.Context, id model.ID) (m *model.ResourceAction, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.ResourceAction
	err = gdb.Where("id = ?", id).First(&dbResult).Error
	if err != nil {
		return
	}
	m = &dbResult
	return
}

func (r *resourceAction) SelectByCode(ctx context.Context, code model.Code) (m *model.ResourceAction, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.ResourceAction
	err = gdb.Where("code = ?", code).First(&dbResult).Error
	if err != nil {
		return
	}
	m = &dbResult
	return
}

func (r *resourceAction) Update(ctx context.Context, m *model.ResourceAction) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(m).Save(m).Error
	return
}

func (r *resourceAction) ListPaged(ctx context.Context, start, limit int) (list []*model.ResourceAction, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.ResourceAction{})
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (r *resourceAction) ListByNameAndResourceServerID(ctx context.Context, resourceServerID model.ID, name string, limit int) (list []*model.ResourceAction, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var whereSQL strings.Builder
	var whereValues []interface{}

	whereSQL.WriteString("name like ?")
	whereValues = append(whereValues, fmt.Sprintf("%%%s%%", name))

	if resourceServerID > 0 {
		whereSQL.WriteString("resource_server_id = ?")
		whereValues = append(whereValues, resourceServerID)
	}
	err = gdb.Model(&model.ResourceAction{}).Where(whereSQL.String(), whereValues...).Offset(0).Limit(limit).Find(&list).Error
	return
}

func (r *resourceAction) ExistByCode(ctx context.Context, code model.Code) (exist bool, err error) {
	return r.exist(ctx, "code = ?", code)
}

func (r *resourceAction) exist(ctx context.Context, query interface{}, args ...interface{}) (exist bool, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var count int64
	err = gdb.Model(&model.ResourceAction{}).Where(query, args...).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		exist = true
	}
	return
}

func (r *resourceAction) DeleteInIDs(ctx context.Context, ids ...model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("id in ?", ids).Delete(model.ResourceAction{}).Error
	return
}

func (r *resourceAction) ListPagedByResourceServerID(ctx context.Context, resourceServerID model.ID, start, limit int) (list []*model.ResourceAction, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.ResourceAction{}).Where("resource_server_id = ?", resourceServerID)
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (r *resourceAction) ListByResourceServerID(ctx context.Context, resourceID model.ID, limit int) (list []*model.ResourceAction, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	exp := gdb.Model(&model.ResourceAction{}).Where("resource_server_id = ?", resourceID)
	if limit > 0 {
		exp = exp.Offset(0).Limit(limit)
	}
	err = exp.Find(&list).Error
	return
}
