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

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetOne gets one User by ID
func (r *UserRepo) GetOne(ctx context.Context, id string) (*domain.User, error) {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return nil, err
	}

	// Gets from db
	db := r.db.WithContext(ctx)

	srcModel := new(src.User)
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

// GetList gets Users list with optional filters
func (r *UserRepo) GetList(ctx context.Context, f []*filters.Filter) ([]*domain.User, error) {
	// string filters -> src filters
	srcFilters, err := filters.TransformWebToSrc(f, &domain.User{}, &src.User{})
	if err != nil {
		return nil, err
	}

	// Gets from db
	srcModels := make(src.Users, 0, 1)
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

// Create creates new User
func (r *UserRepo) Create(ctx context.Context, modelIn *domain.User) (*domain.User, error) {
	// domain Model -> src model
	srcModel, err := new(src.User).ScanFromWeb(modelIn)
	if err != nil {
		return nil, err
	}

	// db operations
	db := r.db.WithContext(ctx)

	// check existence
	tmpRec := new(src.User)
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

// Update updates User
func (r *UserRepo) Update(ctx context.Context, modelIn *domain.User) (*domain.User, error) {
	// domain Model -> src model
	srcModel, err := new(src.User).ScanFromWeb(modelIn)
	if err != nil {
		return nil, err
	}

	// db operations
	db := r.db.WithContext(ctx)

	// check existence
	tmpRec := new(src.User)
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

// Delete deletes User
func (r *UserRepo) Delete(ctx context.Context, id string) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	// db operations
	db := r.db.WithContext(ctx)
	dbResult := db.Delete(&src.User{}, idSrc)
	if dbResult.Error != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}
	if dbResult.RowsAffected == 0 {
		return omniErr.ErrResourceNotFound
	}

	return nil
}

// AppendCompany appends Company new relation to User by id
func (r *UserRepo) AppendCompany(ctx context.Context, id string, relationData *domain.Company) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	srcModelNew, err := new(src.Company).ScanFromWeb(relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.User{ID: idSrc}).Association("Company").Append(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// ReplaceCompany replaces Company old relation in User by id with new Company
func (r *UserRepo) ReplaceCompany(ctx context.Context, id string, relationData *domain.Company) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	srcModelNew, err := new(src.Company).ScanFromWeb(relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.User{ID: idSrc}).Association("Company").Replace(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// DeleteCompany deletes Company relation in User by id
func (r *UserRepo) DeleteCompany(ctx context.Context, id string, relationData *domain.Company) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	srcModelNew, err := new(src.Company).ScanFromWeb(relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.User{ID: idSrc}).Association("Company").Delete(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// AppendLocation appends Location new relation to User by id
func (r *UserRepo) AppendLocation(ctx context.Context, id string, relationData *domain.Location) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	srcModelNew, err := new(src.Location).ScanFromWeb(relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.User{ID: idSrc}).Association("Location").Append(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// ReplaceLocation replaces Location old relation in User by id with new Location
func (r *UserRepo) ReplaceLocation(ctx context.Context, id string, relationData *domain.Location) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	srcModelNew, err := new(src.Location).ScanFromWeb(relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.User{ID: idSrc}).Association("Location").Replace(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// DeleteLocation deletes Location relation in User by id
func (r *UserRepo) DeleteLocation(ctx context.Context, id string, relationData *domain.Location) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	srcModelNew, err := new(src.Location).ScanFromWeb(relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.User{ID: idSrc}).Association("Location").Delete(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// AppendRoles appends Roles new relation to User by id
func (r *UserRepo) AppendRoles(ctx context.Context, id string, relationData []*domain.Role) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	srcModelNew, err := src.Roles.ScanFromWeb(nil, relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.User{ID: idSrc}).Association("Roles").Append(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// ReplaceRoles replaces Roles old relation in User by id with new Roles
func (r *UserRepo) ReplaceRoles(ctx context.Context, id string, relationData []*domain.Role) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	srcModelNew, err := src.Roles.ScanFromWeb(nil, relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.User{ID: idSrc}).Association("Roles").Replace(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// DeleteRoles deletes Roles relation in User by id
func (r *UserRepo) DeleteRoles(ctx context.Context, id string, relationData []*domain.Role) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	srcModelNew, err := src.Roles.ScanFromWeb(nil, relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.User{ID: idSrc}).Association("Roles").Delete(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// AppendSubscriptions appends Subscriptions new relation to User by id
func (r *UserRepo) AppendSubscriptions(ctx context.Context, id string, relationData []*domain.Subscription) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	srcModelNew, err := src.Subscriptions.ScanFromWeb(nil, relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.User{ID: idSrc}).Association("Subscriptions").Append(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// ReplaceSubscriptions replaces Subscriptions old relation in User by id with new Subscriptions
func (r *UserRepo) ReplaceSubscriptions(ctx context.Context, id string, relationData []*domain.Subscription) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	srcModelNew, err := src.Subscriptions.ScanFromWeb(nil, relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.User{ID: idSrc}).Association("Subscriptions").Replace(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// DeleteSubscriptions deletes Subscriptions relation in User by id
func (r *UserRepo) DeleteSubscriptions(ctx context.Context, id string, relationData []*domain.Subscription) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	srcModelNew, err := src.Subscriptions.ScanFromWeb(nil, relationData)

	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.User{ID: idSrc}).Association("Subscriptions").Delete(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}
