package store

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"omnimanage/internal/service"
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

func (r *UserRepo) ModifyRelation(ctx context.Context, id string, relationName string, operation int, relationData interface{}) error {
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	var srcModelNew interface{}
	switch relationName {
	case "Location":
		domModel, ok := relationData.(*domain.Location)
		if !ok {
			return fmt.Errorf("%w wrong relation data", omniErr.ErrInternal)
		}

		srcModelNew, err = new(src.Location).ScanFromWeb(domModel)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		}

	case "Roles":
		domModel, ok := relationData.([]*domain.Role)
		if !ok {
			return fmt.Errorf("%w wrong relation data", omniErr.ErrInternal)
		}

		srcModelNew, err = src.Roles.ScanFromWeb(nil, domModel)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		}

	default:
		return fmt.Errorf("%w wrong relation name '%v'", omniErr.ErrResourceNotFound, relationName)

	}

	db := r.db.WithContext(ctx)

	switch operation {
	case service.OperationRelationReplace:
		err = db.Model(&src.User{ID: idSrc}).Association(relationName).Replace(srcModelNew)
	case service.OperationRelationAppend:
		err = db.Model(&src.User{ID: idSrc}).Association(relationName).Append(srcModelNew)
	case service.OperationRelationDelete:
		err = db.Model(&src.User{ID: idSrc}).Association(relationName).Delete(srcModelNew)
	}
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

func (r *UserRepo) ReplaceRelation(ctx context.Context, id string, relationName string, relationData interface{}) error {
	//idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	//if err != nil {
	//	return err
	//}
	//
	//var srcModelNew interface{}
	//switch relationName {
	//case "Location":
	//	domModel, ok := relationData.(*domain.Location)
	//	if !ok {
	//		return fmt.Errorf("%w wrong relation data", omniErr.ErrInternal)
	//	}
	//
	//	srcModelNew, err = new(src.Location).ScanFromWeb(domModel)
	//	if err != nil {
	//		return err
	//	}
	//
	//case "Roles":
	//	domModel, ok := relationData.([]*domain.Role)
	//	if !ok {
	//		return fmt.Errorf("%w wrong relation data", omniErr.ErrInternal)
	//	}
	//	srcModelNew, err = src.Roles.ScanFromWeb(nil, domModel)
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//db := r.db.WithContext(ctx)
	//err = db.Model(&src.User{ID: idSrc}).Association(relationName).Replace(srcModelNew)
	//if err != nil {
	//	return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	//}

	return nil
}

func (r *UserRepo) AppendRelation(ctx context.Context, id string, relationName string, relationData interface{}) error {
	//db := r.db.WithContext(ctx)
	//
	//err := db.Model(&src.User{ID: id}).Association(relationName).Append(relationData)
	//if err != nil {
	//	return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	//}

	return nil
}

func (r *UserRepo) DeleteRelation(ctx context.Context, id string, relationName string, relationData interface{}) error {
	//db := r.db.WithContext(ctx)
	//
	//err := db.Model(&src.User{ID: id}).Association(relationName).Delete(relationData)
	//if err != nil {
	//	return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	//}

	return nil
}
