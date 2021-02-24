package store

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"omnimanage/internal/model"
	omnierror "omnimanage/pkg/error"
	"omnimanage/pkg/filters"
)

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) GetOne(ctx context.Context, id int) (*model.Role, error) {
	db := r.db.Debug().WithContext(ctx)

	rec := new(model.Role)
	dbResult := db.Where("id = ?", id).Preload(clause.Associations).First(rec)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, omnierror.ErrResourceNotFound
	} else if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, dbResult.Error)
	}

	return rec, nil
}

func (r *RoleRepo) GetList(ctx context.Context, f []*filters.Filter) (model.Roles, error) {
	res := make([]*model.Role, 0, 1)

	db := r.db.Debug().WithContext(ctx)
	db, err := filters.SetGormFilters(db, &res, f)
	if err != nil {
		return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, err)
	}

	dbResult := db.Preload(clause.Associations).Find(&res)
	if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, dbResult.Error)
	}

	if dbResult.RowsAffected == 0 {
		return nil, omnierror.ErrResourceNotFound
	}

	return res, nil
}

func (r *RoleRepo) Create(ctx context.Context, modelIn *model.Role) (*model.Role, error) {

	db := r.db.Debug().WithContext(ctx)

	tmpRec := new(model.Role)
	dbResult := db.Where("id = ?", modelIn.ID).First(tmpRec)
	if dbResult.Error != nil && !errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, dbResult.Error)
	}
	if dbResult.RowsAffected > 0 {
		return nil, fmt.Errorf("%w", omnierror.ErrResourceExists)
	}

	dbResult = db.Preload(clause.Associations).Create(&modelIn)
	if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, dbResult.Error)
	}

	return modelIn, nil
}

func (r *RoleRepo) Update(ctx context.Context, modelIn *model.Role) (*model.Role, error) {
	db := r.db.Debug().WithContext(ctx)

	tmpRec := new(model.Role)
	dbResult := db.Where("id = ?", modelIn.ID).First(tmpRec)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, omnierror.ErrResourceNotFound
	} else if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, dbResult.Error)
	}

	dbResult = db.Preload(clause.Associations).Save(&modelIn)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, omnierror.ErrResourceNotFound
	} else if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, dbResult.Error)
	}

	return modelIn, nil
}
