package store

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	omniErr "omnimanage/pkg/error"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/model/src"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetOne(ctx context.Context, id int) (*src.User, error) {
	db := r.db.WithContext(ctx)

	rec := new(src.User)
	dbResult := db.Where("id = ?", id).Preload(clause.Associations).First(rec)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, omniErr.ErrResourceNotFound
	} else if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	return rec, nil
}

func (r *UserRepo) GetList(ctx context.Context, f []*filters.Filter) (src.Users, error) {

	srcUsers := make(src.Users, 0, 1)

	db := r.db.WithContext(ctx)
	db, err := filters.SetGormFilters(db, &srcUsers, f)
	if err != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	dbResult := db.Preload(clause.Associations).Find(&srcUsers)
	if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	if dbResult.RowsAffected == 0 {
		return nil, omniErr.ErrResourceNotFound
	}

	return srcUsers, nil
}

func (r *UserRepo) Create(ctx context.Context, modelIn *src.User) (*src.User, error) {

	db := r.db.WithContext(ctx)

	tmpRec := new(src.User)
	dbResult := db.Where("id = ?", modelIn.ID).First(tmpRec)
	if dbResult.Error != nil && !errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}
	if dbResult.RowsAffected > 0 {
		return nil, fmt.Errorf("%w", omniErr.ErrResourceExists)
	}

	dbResult = db.Preload(clause.Associations).Create(&modelIn)
	if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	return modelIn, nil
}

func (r *UserRepo) Update(ctx context.Context, modelIn *src.User) (*src.User, error) {
	db := r.db.WithContext(ctx)

	tmpRec := new(src.User)
	dbResult := db.Where("id = ?", modelIn.ID).First(tmpRec)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, omniErr.ErrResourceNotFound
	} else if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	dbResult = db.Preload(clause.Associations).Save(&modelIn)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, omniErr.ErrResourceNotFound
	} else if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	return modelIn, nil
}

func (r *UserRepo) Delete(ctx context.Context, id int) error {
	db := r.db.WithContext(ctx)

	dbResult := db.Delete(&src.User{}, id)
	if dbResult.Error != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}
	if dbResult.RowsAffected == 0 {
		return omniErr.ErrResourceNotFound
	}
	return nil
}

func (r *UserRepo) ReplaceRelation(ctx context.Context, id int, relationName string, relationData interface{}) error {
	db := r.db.WithContext(ctx)

	err := db.Model(&src.User{ID: id}).Association(relationName).Replace(relationData)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

func (r *UserRepo) AppendRelation(ctx context.Context, id int, relationName string, relationData interface{}) error {
	db := r.db.WithContext(ctx)

	err := db.Model(&src.User{ID: id}).Association(relationName).Append(relationData)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

func (r *UserRepo) DeleteRelation(ctx context.Context, id int, relationName string, relationData interface{}) error {
	db := r.db.WithContext(ctx)

	err := db.Model(&src.User{ID: id}).Association(relationName).Delete(relationData)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}
