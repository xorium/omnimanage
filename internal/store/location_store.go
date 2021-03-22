package store

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	omniErr "omnimanage/pkg/error"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/mapper"
	"omnimanage/pkg/model/domain"
	"omnimanage/pkg/model/src"
)

type LocationRepo struct {
	db *gorm.DB
}

func NewLocationRepo(db *gorm.DB) *LocationRepo {
	return &LocationRepo{db: db}
}

// GetOne gets one Location by ID
func (r *LocationRepo) GetOne(ctx context.Context, id string) (*domain.Location, error) {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Location{})
	if err != nil {
		return nil, err
	}

	// Gets from db
	db := r.db.WithContext(ctx)

	srcModel := new(src.Location)
	dbResult := db.Where("id = ?", idSrc).Preload(clause.Associations).First(srcModel)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, omniErr.ErrResourceNotFound
	} else if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	// src Model -> domain
	dom, err := srcModel.ToWeb()
	if err != nil {
		return nil, err
	}

	return dom, nil
}

// GetList gets Locations list with optional filters
func (r *LocationRepo) GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Location, error) {
	// string filters -> src filters
	srcFilters, err := filters.TransformWebToSrc(f, &domain.Location{}, &src.Location{})
	if err != nil {
		return nil, err
	}

	// Gets from db
	srcModels := make(src.Locations, 0, 1)
	db := r.db.WithContext(ctx)
	db, err = filters.SetGormFilters(db, &srcModels, srcFilters)
	if err != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	dbResult := db.Preload(clause.Associations).Find(&srcModels)
	if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	if dbResult.RowsAffected == 0 {
		return nil, omniErr.ErrResourceNotFound
	}

	// src Model -> domain
	domModels, err := srcModels.ToWeb()
	if err != nil {
		return nil, err
	}

	return domModels, nil
}

// Create creates new Location
func (r *LocationRepo) Create(ctx context.Context, modelIn *domain.Location) (*domain.Location, error) {
	// domain Model -> src model
	srcModel, err := new(src.Location).ScanFromWeb(modelIn)
	if err != nil {
		return nil, err
	}

	// db operations
	db := r.db.WithContext(ctx)

	// check existence
	tmpRec := new(src.Location)
	dbResult := db.Where("id = ?", srcModel.ID).First(tmpRec)
	if dbResult.Error != nil && !errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}
	if dbResult.RowsAffected > 0 {
		return nil, fmt.Errorf("%w", omniErr.ErrResourceExists)
	}

	// create data
	dbResult = db.Preload(clause.Associations).Create(&srcModel)
	if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	// src Model -> domain model
	domModel, err := srcModel.ToWeb()
	if err != nil {
		return nil, err
	}

	return domModel, nil
}

// Update updates Location
func (r *LocationRepo) Update(ctx context.Context, modelIn *domain.Location) (*domain.Location, error) {
	// domain Model -> src model
	srcModel, err := new(src.Location).ScanFromWeb(modelIn)
	if err != nil {
		return nil, err
	}

	// db operations
	db := r.db.WithContext(ctx)

	// check existence
	tmpRec := new(src.Location)
	dbResult := db.Where("id = ?", srcModel.ID).First(tmpRec)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, omniErr.ErrResourceNotFound
	} else if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	// update data
	dbResult = db.Preload(clause.Associations).Save(&srcModel)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, omniErr.ErrResourceNotFound
	} else if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	// src Model -> domain model
	domModel, err := srcModel.ToWeb()
	if err != nil {
		return nil, err
	}

	return domModel, nil
}

// Delete deletes Location
func (r *LocationRepo) Delete(ctx context.Context, id string) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Location{})
	if err != nil {
		return err
	}

	// db operations
	db := r.db.WithContext(ctx)
	dbResult := db.Delete(&src.Location{}, idSrc)
	if dbResult.Error != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}
	if dbResult.RowsAffected == 0 {
		return omniErr.ErrResourceNotFound
	}

	return nil
}

// AppendCompany appends Company new relation to Location by id
func (r *LocationRepo) AppendCompany(ctx context.Context, id string, relationData *domain.Company) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Location{})
	if err != nil {
		return err
	}

	srcModelNew, err := new(src.Company).ScanFromWeb(relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.Location{ID: idSrc}).Association("Company").Append(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// ReplaceCompany replaces Company old relation in Location by id with new Company
func (r *LocationRepo) ReplaceCompany(ctx context.Context, id string, relationData *domain.Company) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Location{})
	if err != nil {
		return err
	}

	srcModelNew, err := new(src.Company).ScanFromWeb(relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.Location{ID: idSrc}).Association("Company").Replace(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// DeleteCompany deletes Company relation in Location by id
func (r *LocationRepo) DeleteCompany(ctx context.Context, id string, relationData *domain.Company) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Location{})
	if err != nil {
		return err
	}

	srcModelNew, err := new(src.Company).ScanFromWeb(relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.Location{ID: idSrc}).Association("Company").Delete(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// AppendChildren appends Children new relation to Location by id
func (r *LocationRepo) AppendChildren(ctx context.Context, id string, relationData []*domain.Location) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Location{})
	if err != nil {
		return err
	}

	srcModelNew, err := src.Locations.ScanFromWeb(nil, relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.Location{ID: idSrc}).Association("Children").Append(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// ReplaceChildren replaces Children old relation in Location by id with new Children
func (r *LocationRepo) ReplaceChildren(ctx context.Context, id string, relationData []*domain.Location) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Location{})
	if err != nil {
		return err
	}

	srcModelNew, err := src.Locations.ScanFromWeb(nil, relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.Location{ID: idSrc}).Association("Children").Replace(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// DeleteChildren deletes Children relation in Location by id
func (r *LocationRepo) DeleteChildren(ctx context.Context, id string, relationData []*domain.Location) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Location{})
	if err != nil {
		return err
	}

	srcModelNew, err := src.Locations.ScanFromWeb(nil, relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.Location{ID: idSrc}).Association("Children").Delete(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// AppendUsers appends Users new relation to Location by id
func (r *LocationRepo) AppendUsers(ctx context.Context, id string, relationData []*domain.User) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Location{})
	if err != nil {
		return err
	}

	srcModelNew, err := src.Users.ScanFromWeb(nil, relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.Location{ID: idSrc}).Association("Users").Append(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// ReplaceUsers replaces Users old relation in Location by id with new Users
func (r *LocationRepo) ReplaceUsers(ctx context.Context, id string, relationData []*domain.User) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Location{})
	if err != nil {
		return err
	}

	srcModelNew, err := src.Users.ScanFromWeb(nil, relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.Location{ID: idSrc}).Association("Users").Replace(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// DeleteUsers deletes Users relation in Location by id
func (r *LocationRepo) DeleteUsers(ctx context.Context, id string, relationData []*domain.User) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Location{})
	if err != nil {
		return err
	}

	srcModelNew, err := src.Users.ScanFromWeb(nil, relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.Location{ID: idSrc}).Association("Users").Delete(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}
