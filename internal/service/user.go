package service

import (
	"context"
	"fmt"
	"net/http"
	"omnimanage/internal/store"
	omniErr "omnimanage/pkg/error"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/mapper"
	"omnimanage/pkg/model/domain"
	"omnimanage/pkg/model/src"
)

type UserDomService struct {
	store *store.Store
}

func NewUserService(store *store.Store) *UserDomService {
	return &UserDomService{
		store: store,
	}
}

func (svc *UserDomService) GetOne(ctx context.Context, id string) (*domain.User, error) {
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return nil, err
	}

	src, err := svc.store.Users.GetOne(ctx, idSrc)
	if err != nil {
		return nil, err
	}

	dom, err := src.ToWeb()
	if err != nil {
		return nil, err
	}

	return dom, nil
}

func (svc *UserDomService) GetList(ctx context.Context, f []*filters.Filter) ([]*domain.User, error) {

	srcFilters, err := filters.TransformWebToSrc(f, &domain.User{}, &src.User{})
	if err != nil {
		return nil, err
	}

	src, err := svc.store.Users.GetList(ctx, srcFilters)
	if err != nil {
		return nil, err
	}

	dom, err := src.ToWeb()
	if err != nil {
		return nil, err
	}

	return dom, nil
}

func (svc *UserDomService) Create(ctx context.Context, modelIn *domain.User) (*domain.User, error) {
	srcUser, err := new(src.User).ScanFromWeb(modelIn)
	if err != nil {
		return nil, err
	}

	user, err := svc.store.Users.Create(ctx, srcUser)

	domUser, err := user.ToWeb()
	if err != nil {
		return nil, err
	}

	return domUser, nil
}

func (svc *UserDomService) Update(ctx context.Context, modelIn *domain.User) (*domain.User, error) {
	srcUser, err := new(src.User).ScanFromWeb(modelIn)
	if err != nil {
		return nil, err
	}

	user, err := svc.store.Users.Update(ctx, srcUser)

	domUser, err := user.ToWeb()
	if err != nil {
		return nil, err
	}

	return domUser, nil
}

func (svc *UserDomService) Delete(ctx context.Context, id string) error {
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	err = svc.store.Users.Delete(ctx, idSrc)
	if err != nil {
		return err
	}

	return nil
}

func (svc *UserDomService) GetRelation(ctx context.Context, id string, relationName string) (interface{}, error) {
	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return nil, err
	}

	user, err := svc.store.Users.GetOne(ctx, idSrc)
	if err != nil {
		return nil, err
	}

	switch relationName {
	case "location":
		loc, err := svc.store.Locations.GetOne(ctx, user.LocationID)
		if err != nil {
			return nil, err
		}

		dom, err := loc.ToWeb()
		if err != nil {
			return nil, err
		}
		return dom, nil

	case "roles":
		srcFilters, err := filters.GetSrcFiltersFromRelationID(user.Roles)
		if err != nil {
			return nil, err
		}

		srcList, err := svc.store.Roles.GetList(ctx, srcFilters)
		if err != nil {
			return nil, err
		}

		dom, err := srcList.ToWeb()
		if err != nil {
			return dom, nil
		}
		return dom, nil

	default:
		return nil, fmt.Errorf("%w wrong relation name '%v'", omniErr.ErrResourceNotFound, relationName)
	}
}

func (svc *UserDomService) ModifyRelation(ctx context.Context, id string, relationName string, operation int, relationData interface{}) error {
	return nil

	idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
	if err != nil {
		return err
	}

	switch relationName {
	case "location":
		domModel, ok := relationData.(*domain.Location)
		if !ok {
			return fmt.Errorf("%w wrong relation data", omniErr.ErrInternal)
		}

		srcModelsNew, err := new(src.Location).ScanFromWeb(domModel)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		}

		srcRelName := "Location"
		switch operation {
		case OperationRelationReplace:
			err = svc.store.Users.ReplaceRelation(ctx, idSrc, srcRelName, srcModelsNew)
		case OperationRelationAppend:
			err = svc.store.Users.AppendRelation(ctx, idSrc, srcRelName, srcModelsNew)
		case OperationRelationDelete:
			err = svc.store.Users.DeleteRelation(ctx, idSrc, srcRelName, srcModelsNew)
		}
		if err != nil {
			return err
		}

	case "roles":
		domModel, ok := relationData.([]*domain.Role)
		if !ok {
			return fmt.Errorf("%w wrong relation data", omniErr.ErrInternal)
		}

		srcModelsNew, err := src.Roles.ScanFromWeb(nil, domModel)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		}

		srcRelName := "Roles"
		switch operation {
		case OperationRelationReplace:
			err = svc.store.Users.ReplaceRelation(ctx, idSrc, srcRelName, srcModelsNew)
		case OperationRelationAppend:
			err = svc.store.Users.AppendRelation(ctx, idSrc, srcRelName, srcModelsNew)
		case OperationRelationDelete:
			err = svc.store.Users.DeleteRelation(ctx, idSrc, srcRelName, srcModelsNew)
		}
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("%w wrong relation name '%v'", omniErr.ErrResourceNotFound, relationName)

	}

	return nil
}
