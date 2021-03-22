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

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

// GetOne gets one Role by ID
func (r *RoleRepo) GetOne(ctx context.Context, id string) (*domain.Role, error) {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Role{})
	if err != nil {
		return nil, err
	}

	// Gets from db
	db := r.db.WithContext(ctx)

	srcModel := new(src.Role)
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

// GetList gets Roles list with optional filters
func (r *RoleRepo) GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Role, error) {
	// string filters -> src filters
	srcFilters, err := filters.TransformWebToSrc(f, &domain.Role{}, &src.Role{})
	if err != nil {
		return nil, err
	}

	// Gets from db
	srcModels := make(src.Roles, 0, 1)
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

// Create creates new Role
func (r *RoleRepo) Create(ctx context.Context, modelIn *domain.Role) (*domain.Role, error) {
	// domain Model -> src model
	srcModel, err := new(src.Role).ScanFromWeb(modelIn)
	if err != nil {
		return nil, err
	}

	// db operations
	db := r.db.WithContext(ctx)

	// check existence
	tmpRec := new(src.Role)
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

// Update updates Role
func (r *RoleRepo) Update(ctx context.Context, modelIn *domain.Role) (*domain.Role, error) {
	// domain Model -> src model
	srcModel, err := new(src.Role).ScanFromWeb(modelIn)
	if err != nil {
		return nil, err
	}

	// db operations
	db := r.db.WithContext(ctx)

	// check existence
	tmpRec := new(src.Role)
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

// Delete deletes Role
func (r *RoleRepo) Delete(ctx context.Context, id string) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Role{})
	if err != nil {
		return err
	}

	// db operations
	db := r.db.WithContext(ctx)
	dbResult := db.Delete(&src.Role{}, idSrc)
	if dbResult.Error != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}
	if dbResult.RowsAffected == 0 {
		return omniErr.ErrResourceNotFound
	}

	return nil
}

// AppendCompany appends Company new relation to Role by id
func (r *RoleRepo) AppendCompany(ctx context.Context, id string, relationData *domain.Company) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Role{})
	if err != nil {
		return err
	}

	srcModelNew, err := new(src.Company).ScanFromWeb(relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.Role{ID: idSrc}).Association("Company").Append(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// ReplaceCompany replaces Company old relation in Role by id with new Company
func (r *RoleRepo) ReplaceCompany(ctx context.Context, id string, relationData *domain.Company) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Role{})
	if err != nil {
		return err
	}

	srcModelNew, err := new(src.Company).ScanFromWeb(relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.Role{ID: idSrc}).Association("Company").Replace(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// DeleteCompany deletes Company relation in Role by id
func (r *RoleRepo) DeleteCompany(ctx context.Context, id string, relationData *domain.Company) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.Role{})
	if err != nil {
		return err
	}

	srcModelNew, err := new(src.Company).ScanFromWeb(relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.Role{ID: idSrc}).Association("Company").Delete(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}
