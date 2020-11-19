package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"gorm.io/gorm"
)

// Roleer ...
type Roleer interface {
	Insert(ctx context.Context, m *model.Role) (err error)
	Delete(ctx context.Context, code model.Code) (err error)
	DeleteInCodes(ctx context.Context, codes ...model.Code) (err error)
	SelectByCode(ctx context.Context, code model.Code) (m *model.Role, err error)
	SelectByRoot(ctx context.Context) (results []*model.Role, err error)
	SelectByRootAndOrganizationID(ctx context.Context, organizationID model.ID) (results []*model.Role, err error)
	SelectByParentCode(ctx context.Context, parentCode model.Code) (results []*model.Role, err error)
	Update(ctx context.Context, m *model.Role) (err error)
	ListPaged(ctx context.Context, start, limit int) (list []*model.Role, total int64, err error)
	ListPagedByOrganizationID(ctx context.Context, organizationID model.ID, start, limit int) (list []*model.Role, total int64, err error)
	ListByNameAndOrganizationID(ctx context.Context, organizationID model.ID, name string, limit int) (list []*model.Role, err error)
	ExistByCode(ctx context.Context, code model.Code) (exist bool, err error)
}

type role struct {
}

func (r *role) Insert(ctx context.Context, m *model.Role) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(m).Error
	return
}

func (r *role) Delete(ctx context.Context, code model.Code) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.Role{}, "code = ?", code).Error
	return
}

func (r *role) SelectByCode(ctx context.Context, code model.Code) (m *model.Role, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.Role
	err = gdb.Where("code = ?", code).First(&dbResult).Error
	if err != nil {
		return
	}
	m = &dbResult
	return
}

func (r *role) SelectByRoot(ctx context.Context) (results []*model.Role, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("ISNULL(parent_code) OR parent_code = ''").Find(&results).Error
	return
}

func (r *role) SelectByRootAndOrganizationID(ctx context.Context, organizationID model.ID) (results []*model.Role, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("(ISNULL(parent_code) OR parent_code = '') and organization_id = ?", organizationID).Find(&results).Error
	return
}

func (r *role) SelectByParentCode(ctx context.Context, parentCode model.Code) (results []*model.Role, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("parent_code = ?", parentCode).Find(&results).Error
	return
}

func (r *role) Update(ctx context.Context, m *model.Role) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(m).Save(m).Error
	return
}

func (r *role) ListPaged(ctx context.Context, start, limit int) (list []*model.Role, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.Role{})
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (r *role) ListByNameAndOrganizationID(ctx context.Context, organizationID model.ID, name string, limit int) (list []*model.Role, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var whereSQL strings.Builder
	var whereValues []interface{}

	whereSQL.WriteString("name like ?")
	whereValues = append(whereValues, fmt.Sprintf("%%%s%%", name))

	if organizationID > 0 {
		whereSQL.WriteString("organization_id = ?")
		whereValues = append(whereValues, organizationID)
	}
	err = gdb.Model(&model.Role{}).Where(whereSQL.String(), whereValues...).Offset(0).Limit(limit).Find(&list).Error
	return
}

func (r *role) ExistByCode(ctx context.Context, code model.Code) (exist bool, err error) {
	return r.exist(ctx, "code = ?", code)
}

func (r *role) exist(ctx context.Context, query interface{}, args ...interface{}) (exist bool, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var count int64
	err = gdb.Model(&model.Role{}).Where(query, args...).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		exist = true
	}
	return
}

func (r *role) DeleteInCodes(ctx context.Context, codes ...model.Code) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("code in ?", codes).Delete(model.Role{}).Error
	return
}

func (r *role) ListPagedByOrganizationID(ctx context.Context, organizationID model.ID, start, limit int) (list []*model.Role, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.Role{}).Where("organization_id = ?", organizationID)
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}
