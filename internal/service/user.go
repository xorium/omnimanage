package service

//
//import (
//	"context"
//	"fmt"
//	"omnimanage/internal/store"
//	omniErr "omnimanage/pkg/error"
//	"omnimanage/pkg/filters"
//	"omnimanage/pkg/model/domain"
//)
//
//type UserDomService struct {
//	store *store.Store
//}
//
//func NewUserService(store *store.Store) *UserDomService {
//	return &UserDomService{
//		store: store,
//	}
//}
//
//func (svc *UserDomService) GetOne(ctx context.Context, id string) (*domain.User, error) {
//	return svc.store.User.GetOne(ctx, id)
//	//idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
//	//if err != nil {
//	//	return nil, err
//	//}
//	//src, err := svc.store.Users.GetOne(ctx, idSrc)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//dom, err := src.ToWeb()
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//return dom, nil
//}
//
//func (svc *UserDomService) GetList(ctx context.Context, f []*filters.Filter) ([]*domain.User, error) {
//	return svc.store.User.GetList(ctx, f)
//	//srcFilters, err := filters.TransformWebToSrc(f, &domain.User{}, &src.User{})
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//src, err := svc.store.Users.GetList(ctx, srcFilters)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//dom, err := src.ToWeb()
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//return dom, nil
//}
//
//func (svc *UserDomService) Create(ctx context.Context, modelIn *domain.User) (*domain.User, error) {
//	return svc.store.User.Create(ctx, modelIn)
//	//srcUser, err := new(src.User).ScanFromWeb(modelIn)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//user, err := svc.store.Users.Create(ctx, srcUser)
//	//
//	//domUser, err := user.ToWeb()
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//return domUser, nil
//}
//
//func (svc *UserDomService) Update(ctx context.Context, modelIn *domain.User) (*domain.User, error) {
//	return svc.store.User.Update(ctx, modelIn)
//	//srcUser, err := new(src.User).ScanFromWeb(modelIn)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//user, err := svc.store.Users.Update(ctx, srcUser)
//	//
//	//domUser, err := user.ToWeb()
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//return domUser, nil
//}
//
//func (svc *UserDomService) Delete(ctx context.Context, id string) error {
//	return svc.store.User.Delete(ctx, id)
//	//idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
//	//if err != nil {
//	//	return err
//	//}
//	//
//	//err = svc.store.Users.Delete(ctx, idSrc)
//	//if err != nil {
//	//	return err
//	//}
//	//
//	//return nil
//}
//
//func (svc *UserDomService) GetRelation(ctx context.Context, id string, relationName string) (interface{}, error) {
//	//idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	user, err := svc.store.User.GetOne(ctx, id)
//	if err != nil {
//		return nil, err
//	}
//
//	switch relationName {
//	case "location":
//		if user.Location == nil {
//			return nil, omniErr.ErrResourceNotFound
//		}
//		return svc.store.Location.GetOne(ctx, user.Location.ID)
//		//loc, err := svc.store.Locations.GetOne(ctx, user.LocationID)
//		//if err != nil {
//		//	return nil, err
//		//}
//		//
//		//dom, err := loc.ToWeb()
//		//if err != nil {
//		//	return nil, err
//		//}
//		//return dom, nil
//
//	case "roles":
//		filters, err := filters.TransformModelsIDToFilters(user.Roles)
//		if err != nil {
//			return nil, err
//		}
//		return svc.store.Role.GetList(ctx, filters)
//		//srcList, err := svc.store.Roles.GetList(ctx, filters)
//		//if err != nil {
//		//	return nil, err
//		//}
//		//
//		//dom, err := srcList.ToWeb()
//		//if err != nil {
//		//	return dom, nil
//		//}
//		//return dom, nil
//
//	default:
//		return nil, fmt.Errorf("%w wrong relation name '%v'", omniErr.ErrResourceNotFound, relationName)
//	}
//}
//
//func (svc *UserDomService) ModifyRelation(ctx context.Context, id string, relationName string, operation int, relationModel interface{}) error {
//	return svc.store.User.ModifyRelation(ctx, id, relationName, operation, relationModel)
//
//	//idSrc, err := mapper.Get().GetSrcID(id, &src.User{})
//	//if err != nil {
//	//	return err
//	//}
//	//switch relationName {
//	//case "location":
//	//	domModel, ok := relationData.(*domain.Location)
//	//	if !ok {
//	//		return fmt.Errorf("%w wrong relation data", omniErr.ErrInternal)
//	//	}
//	//
//	//	srcModelNew, err := new(src.Location).ScanFromWeb(domModel)
//	//	if err != nil {
//	//		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//	//	}
//	//
//	//	srcRelName := "Location"
//	//	//switch operation {
//	//	//case OperationRelationReplace:
//	//	//	err = svc.store.Users.ReplaceRelation(ctx, id, srcRelName, srcModelNew)
//	//	//case OperationRelationAppend:
//	//	//	err = svc.store.Users.AppendRelation(ctx, id, srcRelName, srcModelNew)
//	//	//case OperationRelationDelete:
//	//	//	err = svc.store.Users.DeleteRelation(ctx, id, srcRelName, srcModelNew)
//	//	//}
//	//	//if err != nil {
//	//	//	return err
//	//	//}
//	//
//	//case "roles":
//	//	domModel, ok := relationData.([]*domain.Role)
//	//	if !ok {
//	//		return fmt.Errorf("%w wrong relation data", omniErr.ErrInternal)
//	//	}
//	//
//	//	srcModelsNew, err := src.Roles.ScanFromWeb(nil, domModel)
//	//	if err != nil {
//	//		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//	//	}
//	//
//	//	srcRelName := "Roles"
//	//	switch operation {
//	//	case OperationRelationReplace:
//	//		err = svc.store.Users.ReplaceRelation(ctx, idSrc, srcRelName, srcModelsNew)
//	//	case OperationRelationAppend:
//	//		err = svc.store.Users.AppendRelation(ctx, idSrc, srcRelName, srcModelsNew)
//	//	case OperationRelationDelete:
//	//		err = svc.store.Users.DeleteRelation(ctx, idSrc, srcRelName, srcModelsNew)
//	//	}
//	//	if err != nil {
//	//		return err
//	//	}
//	//
//	//default:
//	//	return fmt.Errorf("%w wrong relation name '%v'", omniErr.ErrResourceNotFound, relationName)
//	//
//	//}
//
//	//return nil
//}
