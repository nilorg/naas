package dao

import (
	"context"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"gorm.io/gorm"
)

// ScopeResourceRelationer ...
type ScopeResourceRelationer interface {
	Insert(ctx context.Context, roleResourceRoute *model.ScopeResourceRelation) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	DeleteByRelationTypeAndScopeCode(ctx context.Context, relationType model.ScopeResourceRelationType, scopeCode model.Code) (err error)
	DeleteByRelationTypeAndScopeCodeAndRelationID(ctx context.Context, relationType model.ScopeResourceRelationType, scopeCode model.Code, relationID model.ID) (err error)
	Select(ctx context.Context, id model.ID) (roleResourceRoute *model.ScopeResourceRelation, err error)
	SelectAll(ctx context.Context, relationType model.ScopeResourceRelationType) (roleResourceRoutes []*model.ScopeResourceRelation, err error)
	Update(ctx context.Context, roleResourceRoute *model.ScopeResourceRelation) (err error)
	ExistByRelationTypeAndScopeCodeAndRelationID(ctx context.Context, relationType model.ScopeResourceRelationType, scopeCode model.Code, relationID model.ID) (exist bool, err error)
	ExistByRelationTypeAndRelationID(ctx context.Context, relationType model.ScopeResourceRelationType, relationID model.ID) (exist bool, err error)
	ListByRelationTypeAndScopeCode(ctx context.Context, relationType model.ScopeResourceRelationType, scopeCode model.Code, limit int) (list []*model.ScopeResourceRelation, err error)
	ListPagedByRelationTypeAndScopeCode(ctx context.Context, start, limit int, relationType model.ScopeResourceRelationType, scopeCode model.Code) (list []*model.ScopeResourceRelation, total int64, err error)
	ListPagedByRelationTypeAndScopeCodeAndResourceServerID(ctx context.Context, start, limit int, relationType model.ScopeResourceRelationType, scopeCode model.Code, resourceServerID model.ID) (list []*model.ScopeResourceRelation, total int64, err error)
	ListByRelationTypeAndScopeCodeAndResourceServerID(ctx context.Context, relationType model.ScopeResourceRelationType, scopeCode model.Code, resourceServerID model.ID) (list []*model.ScopeResourceRelation, err error)
	ListForRelationIDByRelationTypeAndScopeCodeAndResourceServerID(ctx context.Context, relationType model.ScopeResourceRelationType, scopeCode model.Code, resourceServerID model.ID) (ids []model.ID, err error)
}

type scopeResourceRelation struct {
}

func (*scopeResourceRelation) Insert(ctx context.Context, roleResourceRoute *model.ScopeResourceRelation) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(roleResourceRoute).Error
	return
}

func (*scopeResourceRelation) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.ScopeResourceRelation{}, id).Error
	return
}

func (*scopeResourceRelation) Select(ctx context.Context, id model.ID) (roleResourceRoute *model.ScopeResourceRelation, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	roleResourceRoute = new(model.ScopeResourceRelation)
	err = gdb.Model(roleResourceRoute).Where("id = ?", id).Take(roleResourceRoute).Error
	if err != nil {
		roleResourceRoute = nil
		return
	}
	return
}
func (*scopeResourceRelation) Update(ctx context.Context, roleResourceRoute *model.ScopeResourceRelation) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(roleResourceRoute).Save(roleResourceRoute).Error
	return
}

func (*scopeResourceRelation) SelectAll(ctx context.Context, relationType model.ScopeResourceRelationType) (roleResourceRoutes []*model.ScopeResourceRelation, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.ScopeResourceRelation{}).Where("relation_type = ?", relationType).Find(&roleResourceRoutes).Error
	return
}

func (*scopeResourceRelation) exist(ctx context.Context, query interface{}, args ...interface{}) (exist bool, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var count int64
	err = gdb.Model(&model.ScopeResourceRelation{}).Where(query, args...).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		exist = true
	}
	return
}

// ExistByRelationTypeAndScopeCodeAndRelationID 判断根据ScopeCode和资源关系ID
func (r *scopeResourceRelation) ExistByRelationTypeAndScopeCodeAndRelationID(ctx context.Context, relationType model.ScopeResourceRelationType, scopeCode model.Code, relationID model.ID) (exist bool, err error) {
	return r.exist(ctx, "relation_type = ? and scope_code = ? and relation_id = ?", relationType, scopeCode, relationID)
}

// ExistByRelationTypeAndRelationID 判断根据资源关系ID
func (r *scopeResourceRelation) ExistByRelationTypeAndRelationID(ctx context.Context, relationType model.ScopeResourceRelationType, relationID model.ID) (exist bool, err error) {
	return r.exist(ctx, "relation_type = ? and relation_id = ?", relationType, relationID)
}

func (r *scopeResourceRelation) ListByRelationTypeAndScopeCode(ctx context.Context, relationType model.ScopeResourceRelationType, scopeCode model.Code, limit int) (list []*model.ScopeResourceRelation, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	exp := gdb.Model(&model.ScopeResourceRelation{}).Where("relation_type = ? and scope_code = ?", relationType, scopeCode)
	if limit > 0 {
		exp = exp.Offset(0).Limit(limit)
	}
	err = exp.Find(&list).Error
	return
}

func (r *scopeResourceRelation) ListPagedByRelationTypeAndScopeCode(ctx context.Context, start, limit int, relationType model.ScopeResourceRelationType, scopeCode model.Code) (list []*model.ScopeResourceRelation, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.ScopeResourceRelation{}).Where("relation_type = ? and scope_code = ?", relationType, scopeCode)
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (r *scopeResourceRelation) ListPagedByRelationTypeAndScopeCodeAndResourceServerID(ctx context.Context, start, limit int, relationType model.ScopeResourceRelationType, scopeCode model.Code, resourceServerID model.ID) (list []*model.ScopeResourceRelation, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.ScopeResourceRelation{}).Where("relation_type = ? and scope_code = ? and resource_server_id = ?", relationType, scopeCode, resourceServerID)
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (r *scopeResourceRelation) ListByRelationTypeAndScopeCodeAndResourceServerID(ctx context.Context, relationType model.ScopeResourceRelationType, scopeCode model.Code, resourceServerID model.ID) (list []*model.ScopeResourceRelation, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(&model.ScopeResourceRelation{}).Where("relation_type = ? and scope_code = ? and resource_server_id = ?", relationType, scopeCode, resourceServerID).Find(&list).Error
	return
}

func (r *scopeResourceRelation) ListForRelationIDByRelationTypeAndScopeCodeAndResourceServerID(ctx context.Context, relationType model.ScopeResourceRelationType, scopeCode model.Code, resourceServerID model.ID) (ids []model.ID, err error) {
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
	err = gdb.Model(&model.ScopeResourceRelation{}).Where("relation_type = ? and scope_code = ? and resource_server_id = ?", relationType, scopeCode, resourceServerID).Find(&items).Error
	if err != nil {
		return
	}
	for _, i := range items {
		ids = append(ids, i.RelationID)
	}
	return
}

func (r *scopeResourceRelation) delete(ctx context.Context, query interface{}, args ...interface{}) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where(query, args...).Delete(model.ScopeResourceRelation{}).Error
	return
}

func (r *scopeResourceRelation) DeleteByRelationTypeAndScopeCode(ctx context.Context, relationType model.ScopeResourceRelationType, scopeCode model.Code) (err error) {
	err = r.delete(ctx, "relation_type = ? and scope_code = ?", relationType, scopeCode)
	if err != nil {
		return
	}
	return
}

func (r *scopeResourceRelation) DeleteByRelationTypeAndScopeCodeAndRelationID(ctx context.Context, relationType model.ScopeResourceRelationType, scopeCode model.Code, relationID model.ID) (err error) {
	err = r.delete(ctx, "relation_type = ? and scope_code = ? and relation_id = ?", relationType, scopeCode, relationID)
	if err != nil {
		return
	}
	return
}
