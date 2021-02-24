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

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetOne(ctx context.Context, id int) (*model.User, error) {
	rec := new(model.User)
	result := r.db.Debug().WithContext(ctx).Where("id = ?", id).Preload(clause.Associations).Find(rec)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, omnierror.ErrResourceNotFound
	} else if result.Error != nil {
		return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, result.Error)
	}

	return rec, nil
}

func (r *UserRepo) GetList(ctx context.Context, f []*filters.Filter) (model.Users, error) {
	users := make([]*model.User, 0, 1)

	db := r.db.Debug().WithContext(ctx)
	db, err := filters.SetGormFilters(db, &users, f)
	if err != nil {
		return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, err)
	}

	result := db.Preload(clause.Associations).Find(&users)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, omnierror.ErrResourceNotFound
	} else if result.Error != nil {
		return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, result.Error)
	}

	return users, nil
}
