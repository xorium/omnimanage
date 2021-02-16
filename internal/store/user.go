package store

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"omnimanage/internal/model"
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
	if result.Error != nil {
		return nil, result.Error
	}

	return rec, nil
}

func (r *UserRepo) GetList(ctx context.Context, f []*filters.Filter) (model.Users, error) {
	users := make([]*model.User, 0, 1)

	db := r.db.Debug().WithContext(ctx)
	db, err := filters.SetGormFilters(db, &users, f)
	if err != nil {
		return nil, err
	}

	//result := r.db.Debug().WithContext(ctx).Joins(""+
	//	"JOIN locations on locations.id = users.location_id and locations.id=? ", "1").Joins(""+
	//	"JOIN companies on companies.id = users.company_id AND companies.id=?", 4).Preload(clause.Associations).Find(&users)
	result := db.Preload(clause.Associations).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	//db.SetupJoinTable()

	return users, nil
}
