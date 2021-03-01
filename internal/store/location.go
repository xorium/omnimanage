package store

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/model/src"
)

type LocationRepo struct {
	db *gorm.DB
}

func NewLocationRepo(db *gorm.DB) *LocationRepo {
	return &LocationRepo{db: db}
}

func (r *LocationRepo) GetOne(ctx context.Context, id int) (*src.Location, error) {
	rec := new(src.Location)
	result := r.db.Debug().WithContext(ctx).Where("id = ?", id).Preload(clause.Associations).Find(rec)
	if result.Error != nil {
		return nil, result.Error
	}

	return rec, nil
}

func (r *LocationRepo) GetList(ctx context.Context, f []*filters.Filter) (src.Locations, error) {
	records := make([]*src.Location, 5)

	//db := filt.SetGormFilters(r.db, filters)
	result := r.db.Debug().WithContext(ctx).Preload(clause.Associations).Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}

func (r *LocationRepo) Create(ctx context.Context, modelIn *src.Location) (*src.Location, error) {

	return modelIn, nil
}

func (r *LocationRepo) Update(ctx context.Context, modelIn *src.Location) (*src.Location, error) {

	return modelIn, nil
}

func (r *LocationRepo) Delete(ctx context.Context, id int) error {

	return nil
}
