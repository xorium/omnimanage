package store

import (
	"context"
	"gorm.io/gorm"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/model/domain"
)

type LocationRepo struct {
	db *gorm.DB
}

func NewLocationRepo(db *gorm.DB) *LocationRepo {
	return &LocationRepo{db: db}
}

func (r *LocationRepo) GetOne(ctx context.Context, id string) (*domain.Location, error) {
	//db := r.db.WithContext(ctx)
	//
	//rec := new(src.Location)
	//dbResult := db.Where("id = ?", id).Preload(clause.Associations).Find(rec)
	//if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
	//	return nil, omnierror.ErrResourceNotFound
	//} else if dbResult.Error != nil {
	//	return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, dbResult.Error)
	//}
	//
	//return rec, nil
	return nil, nil
}

func (r *LocationRepo) GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Location, error) {
	//res := make([]*src.Location, 0, 1)
	//
	//db := r.db.WithContext(ctx)
	//db, err := filters.SetGormFilters(db, &res, f)
	//if err != nil {
	//	return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, err)
	//}
	//
	//dbResult := db.Preload(clause.Associations).Find(&res)
	//if dbResult.Error != nil {
	//	return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, dbResult.Error)
	//}
	//
	//if dbResult.RowsAffected == 0 {
	//	return nil, omnierror.ErrResourceNotFound
	//}
	//
	//return res, nil
	return nil, nil
}

func (r *LocationRepo) Create(ctx context.Context, modelIn *domain.Location) (*domain.Location, error) {
	//db := r.db.WithContext(ctx)
	//
	//tmpRec := new(src.Location)
	//dbResult := db.Where("id = ?", modelIn.ID).First(tmpRec)
	//if dbResult.Error != nil && !errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
	//	return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, dbResult.Error)
	//}
	//if dbResult.RowsAffected > 0 {
	//	return nil, fmt.Errorf("%w", omnierror.ErrResourceExists)
	//}
	//
	//dbResult = db.Preload(clause.Associations).Create(&modelIn)
	//if dbResult.Error != nil {
	//	return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, dbResult.Error)
	//}
	//
	//return modelIn, nil
	return nil, nil
}

func (r *LocationRepo) Update(ctx context.Context, modelIn *domain.Location) (*domain.Location, error) {
	//db := r.db.WithContext(ctx)
	//
	//tmpRec := new(src.Location)
	//dbResult := db.Where("id = ?", modelIn.ID).First(tmpRec)
	//if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
	//	return nil, omnierror.ErrResourceNotFound
	//} else if dbResult.Error != nil {
	//	return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, dbResult.Error)
	//}
	//
	//dbResult = db.Preload(clause.Associations).Save(&modelIn)
	//if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
	//	return nil, omnierror.ErrResourceNotFound
	//} else if dbResult.Error != nil {
	//	return nil, fmt.Errorf("%w %v", omnierror.ErrInternal, dbResult.Error)
	//}
	//
	//return modelIn, nil
	return nil, nil
}

func (r *LocationRepo) Delete(ctx context.Context, id string) error {
	//db := r.db.WithContext(ctx)
	//
	//dbResult := db.Delete(&src.Location{}, id)
	//if dbResult.Error != nil {
	//	return fmt.Errorf("%w %v", omnierror.ErrInternal, dbResult.Error)
	//}
	//if dbResult.RowsAffected == 0 {
	//	return omnierror.ErrResourceNotFound
	//}
	return nil
}

func (r *LocationRepo) ReplaceRelation(ctx context.Context, id string, relationName string, relationData interface{}) error {
	//db := r.db.WithContext(ctx)
	//
	//err := db.Model(&src.Location{ID: id}).Association(relationName).Replace(relationData)
	//if err != nil {
	//	return fmt.Errorf("%w %v", omnierror.ErrInternal, err)
	//}

	return nil
}

func (r *LocationRepo) AppendRelation(ctx context.Context, id string, relationName string, relationData interface{}) error {
	//db := r.db.WithContext(ctx)
	//
	//err := db.Model(&src.Location{ID: id}).Association(relationName).Append(relationData)
	//if err != nil {
	//	return fmt.Errorf("%w %v", omnierror.ErrInternal, err)
	//}

	return nil
}

func (r *LocationRepo) DeleteRelation(ctx context.Context, id string, relationName string, relationData interface{}) error {
	//db := r.db.WithContext(ctx)
	//
	//err := db.Model(&src.Location{ID: id}).Association(relationName).Delete(relationData)
	//if err != nil {
	//	return fmt.Errorf("%w %v", omnierror.ErrInternal, err)
	//}

	return nil
}
