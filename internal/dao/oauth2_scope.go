package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
)

// OAuth2Scoper ...
type OAuth2Scoper interface {
	Insert(ctx context.Context, m *model.OAuth2Scope) (err error)
	Delete(ctx context.Context, code string) (err error)
	Select(ctx context.Context, code string) (m *model.OAuth2Scope, err error)
	SelectAll(ctx context.Context) (m []*model.OAuth2Scope, err error)
	SelectByAllBasic(ctx context.Context) (m []*model.OAuth2Scope, err error)
	Update(ctx context.Context, m *model.OAuth2Scope) (err error)
}

type oauth2Scope struct {
}

func (s *oauth2Scope) Insert(ctx context.Context, m *model.OAuth2Scope) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(m).Error
	return
}

func (s *oauth2Scope) Delete(ctx context.Context, code string) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.OAuth2Scope{}, "code = ?", code).Error
	return
}

func (s *oauth2Scope) SelectAll(ctx context.Context) (m []*model.OAuth2Scope, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.OAuth2Scope{}).Find(&m).Error
	return
}

func (s *oauth2Scope) SelectByAllBasic(ctx context.Context) (scopes []*model.OAuth2Scope, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.OAuth2Scope{}).Where("type = ?", model.OAuth2ScopeTypeBasic).Find(&scopes).Error
	return
}

func (s *oauth2Scope) Select(ctx context.Context, code string) (m *model.OAuth2Scope, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.OAuth2Scope
	err = gdb.First(&dbResult, "code = ?", code).Error
	if err != nil {
		return
	}
	m = &dbResult
	return
}

func (s *oauth2Scope) Update(ctx context.Context, m *model.OAuth2Scope) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(m).Update(m).Error
	if err != nil {
		return
	}
	return
}
