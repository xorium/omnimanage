package store

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"omnimanage/internal/model"
	"omnimanage/pkg/filters"
)

type LocationRepo struct {
	db *gorm.DB
}

func NewLocationRepo(db *gorm.DB) *LocationRepo {
	return &LocationRepo{db: db}
}

func (r *LocationRepo) GetOne(ctx context.Context, id int) (*model.Location, error) {
	rec := new(model.Location)
	result := r.db.Debug().WithContext(ctx).Where("id = ?", id).Preload(clause.Associations).Find(rec)
	if result.Error != nil {
		return nil, result.Error
	}

	return rec, nil
}

func (r *LocationRepo) GetList(ctx context.Context, f []*filters.Filter) (model.Locations, error) {
	records := make([]*model.Location, 5)

	//db := filt.SetGormFilters(r.db, filters)
	result := r.db.Debug().WithContext(ctx).Preload(clause.Associations).Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}

func (r *LocationRepo) Create(ctx context.Context, modelIn *model.Location) (*model.Location, error) {

	return modelIn, nil
}

func (r *LocationRepo) Update(ctx context.Context, modelIn *model.Location) (*model.Location, error) {

	return modelIn, nil
}

func (r *LocationRepo) Delete(ctx context.Context, id int) error {

	return nil
}
