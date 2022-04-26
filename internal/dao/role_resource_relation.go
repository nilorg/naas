package dao

import (
	"context"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"gorm.io/gorm"
)

// RoleResourceRelationer ...
type RoleResourceRelationer interface {
	Insert(ctx context.Context, roleResourceRoute *model.RoleResourceRelation) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	DeleteByRelationTypeAndRoleCode(ctx context.Context, relationType model.RoleResourceRelationType, roleCode model.Code) (err error)
	DeleteByRelationTypeAndRoleCodeAndRelationID(ctx context.Context, relationType model.RoleResourceRelationType, roleCode model.Code, relationID model.ID) (err error)
	Select(ctx context.Context, id model.ID) (roleResourceRoute *model.RoleResourceRelation, err error)
	SelectAll(ctx context.Context, relationType model.RoleResourceRelationType) (roleResourceRoutes []*model.RoleResourceRelation, err error)
	Update(ctx context.Context, roleResourceRoute *model.RoleResourceRelation) (err error)
	ExistByRelationTypeAndRoleCodeAndRelationID(ctx context.Context, relationType model.RoleResourceRelationType, roleCode model.Code, relationID model.ID) (exist bool, err error)
	ExistByRelationTypeAndRelationID(ctx context.Context, relationType model.RoleResourceRelationType, relationID model.ID) (exist bool, err error)
	ListByRelationTypeAndRoleCode(ctx context.Context, relationType model.RoleResourceRelationType, roleCode model.Code, limit int) (list []*model.RoleResourceRelation, err error)
	ListPagedByRelationTypeAndRoleCode(ctx context.Context, start, limit int, relationType model.RoleResourceRelationType, roleCode model.Code) (list []*model.RoleResourceRelation, total int64, err error)
	ListPagedByRelationTypeAndRoleCodeAndResourceServerID(ctx context.Context, start, limit int, relationType model.RoleResourceRelationType, roleCode model.Code, resourceServerID model.ID) (list []*model.RoleResourceRelation, total int64, err error)
	ListByRelationTypeAndRoleCodeAndResourceServerID(ctx context.Context, relationType model.RoleResourceRelationType, roleCode model.Code, resourceServerID model.ID) (list []*model.RoleResourceRelation, err error)
	ListForRelationIDByRelationTypeAndRoleCodeAndResourceServerID(ctx context.Context, relationType model.RoleResourceRelationType, roleCode model.Code, resourceServerID model.ID) (ids []model.ID, err error)
}

type roleResourceRelation struct {
}

func (*roleResourceRelation) Insert(ctx context.Context, roleResourceRoute *model.RoleResourceRelation) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(roleResourceRoute).Error
	return
}

func (*roleResourceRelation) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.RoleResourceRelation{}, id).Error
	return
}

func (*roleResourceRelation) Select(ctx context.Context, id model.ID) (roleResourceRoute *model.RoleResourceRelation, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	roleResourceRoute = new(model.RoleResourceRelation)
	err = gdb.Model(roleResourceRoute).Where("id = ?", id).Take(roleResourceRoute).Error
	if err != nil {
		roleResourceRoute = nil
		return
	}
	return
}
func (*roleResourceRelation) Update(ctx context.Context, roleResourceRoute *model.RoleResourceRelation) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(roleResourceRoute).Save(roleResourceRoute).Error
	return
}

func (*roleResourceRelation) SelectAll(ctx context.Context, relationType model.RoleResourceRelationType) (roleResourceRoutes []*model.RoleResourceRelation, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.RoleResourceRelation{}).Where("relation_type = ?", relationType).Find(&roleResourceRoutes).Error
	return
}

func (*roleResourceRelation) exist(ctx context.Context, query interface{}, args ...interface{}) (exist bool, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var count int64
	err = gdb.Model(&model.RoleResourceRelation{}).Where(query, args...).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		exist = true
	}
	return
}

// ExistByRelationTypeAndRoleCodeAndRelationID 判断根据RoleCode和资源关系ID
func (r *roleResourceRelation) ExistByRelationTypeAndRoleCodeAndRelationID(ctx context.Context, relationType model.RoleResourceRelationType, roleCode model.Code, relationID model.ID) (exist bool, err error) {
	return r.exist(ctx, "relation_type = ? and role_code = ? and relation_id = ?", relationType, roleCode, relationID)
}

// ExistByRelationTypeAndRelationID 判断根据资源关系ID
func (r *roleResourceRelation) ExistByRelationTypeAndRelationID(ctx context.Context, relationType model.RoleResourceRelationType, relationID model.ID) (exist bool, err error) {
	return r.exist(ctx, "relation_type = ? and relation_id = ?", relationType, relationID)
}

func (r *roleResourceRelation) ListByRelationTypeAndRoleCode(ctx context.Context, relationType model.RoleResourceRelationType, roleCode model.Code, limit int) (list []*model.RoleResourceRelation, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	exp := gdb.Model(&model.RoleResourceRelation{}).Where("relation_type = ? and role_code = ?", relationType, roleCode)
	if limit > 0 {
		exp = exp.Offset(0).Limit(limit)
	}
	err = exp.Find(&list).Error
	return
}

func (r *roleResourceRelation) ListPagedByRelationTypeAndRoleCode(ctx context.Context, start, limit int, relationType model.RoleResourceRelationType, roleCode model.Code) (list []*model.RoleResourceRelation, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.RoleResourceRelation{}).Where("relation_type = ? and role_code = ?", relationType, roleCode)
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (r *roleResourceRelation) ListPagedByRelationTypeAndRoleCodeAndResourceServerID(ctx context.Context, start, limit int, relationType model.RoleResourceRelationType, roleCode model.Code, resourceServerID model.ID) (list []*model.RoleResourceRelation, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.RoleResourceRelation{}).Where("relation_type = ? and role_code = ? and resource_server_id = ?", relationType, roleCode, resourceServerID)
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (r *roleResourceRelation) ListByRelationTypeAndRoleCodeAndResourceServerID(ctx context.Context, relationType model.RoleResourceRelationType, roleCode model.Code, resourceServerID model.ID) (list []*model.RoleResourceRelation, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(&model.RoleResourceRelation{}).Where("relation_type = ? and role_code = ? and resource_server_id = ?", relationType, roleCode, resourceServerID).Find(&list).Error
	return
}

func (r *roleResourceRelation) ListForRelationIDByRelationTypeAndRoleCodeAndResourceServerID(ctx context.Context, relationType model.RoleResourceRelationType, roleCode model.Code, resourceServerID model.ID) (ids []model.ID, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	// result ...
	type result struct {
		RelationID model.ID `gorm:"column:relation_id"`
	}

	var items []*result
	err = gdb.Model(&model.RoleResourceRelation{}).Where("relation_type = ? and role_code = ? and resource_server_id = ?", relationType, roleCode, resourceServerID).Find(&items).Error
	if err != nil {
		return
	}
	for _, i := range items {
		ids = append(ids, i.RelationID)
	}
	return
}

func (r *roleResourceRelation) delete(ctx context.Context, query interface{}, args ...interface{}) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where(query, args...).Delete(model.RoleResourceRelation{}).Error
	return
}

func (r *roleResourceRelation) DeleteByRelationTypeAndRoleCode(ctx context.Context, relationType model.RoleResourceRelationType, roleCode model.Code) (err error) {
	err = r.delete(ctx, "relation_type = ? and role_code = ?", relationType, roleCode)
	if err != nil {
		return
	}
	return
}

func (r *roleResourceRelation) DeleteByRelationTypeAndRoleCodeAndRelationID(ctx context.Context, relationType model.RoleResourceRelationType, roleCode model.Code, relationID model.ID) (err error) {
	err = r.delete(ctx, "relation_type = ? and role_code = ? and relation_id = ?", relationType, roleCode, relationID)
	if err != nil {
		return
	}
	return
}
