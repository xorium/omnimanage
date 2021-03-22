package service

import (
	"context"
	"omnimanage/internal/store"
	omniErr "omnimanage/pkg/error"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/model/domain"
)

type UserDomService struct {
	store *store.Store
}

func NewUserService(store *store.Store) *UserDomService {
	return &UserDomService{
		store: store,
	}
}

// GetOne gets one User by ID
func (svc *UserDomService) GetOne(ctx context.Context, id string) (*domain.User, error) {
	return svc.store.User.GetOne(ctx, id)
}

// GetList gets Users list with optional filters
func (svc *UserDomService) GetList(ctx context.Context, f []*filters.Filter) ([]*domain.User, error) {
	return svc.store.User.GetList(ctx, f)
}

// Create creates new User
func (svc *UserDomService) Create(ctx context.Context, modelIn *domain.User) (*domain.User, error) {
	return svc.store.User.Create(ctx, modelIn)
}

// Update updates User
func (svc *UserDomService) Update(ctx context.Context, modelIn *domain.User) (*domain.User, error) {
	return svc.store.User.Update(ctx, modelIn)
}

// Delete deletes User
func (svc *UserDomService) Delete(ctx context.Context, id string) error {
	return svc.store.User.Delete(ctx, id)
}

// GetCompany gets User's Company
func (svc *UserDomService) GetCompany(ctx context.Context, id string) (*domain.Company, error) {
	mainModel, err := svc.store.User.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	if mainModel.Company == nil {
		return nil, omniErr.ErrResourceNotFound
	}

	return svc.store.Company.GetOne(ctx, mainModel.Company.ID)

}

// AppendCompany appends Company new relation to User by id
func (svc *UserDomService) AppendCompany(ctx context.Context, id string, relationData *domain.Company) error {
	return svc.store.User.AppendCompany(ctx, id, relationData)
}

// ReplaceCompany replaces Company old relation in User by id with new Company
func (svc *UserDomService) ReplaceCompany(ctx context.Context, id string, relationData *domain.Company) error {
	return svc.store.User.ReplaceCompany(ctx, id, relationData)
}

// DeleteCompany deletes Company relation in User by id
func (svc *UserDomService) DeleteCompany(ctx context.Context, id string, relationData *domain.Company) error {
	return svc.store.User.DeleteCompany(ctx, id, relationData)
}

// GetLocation gets User's Location
func (svc *UserDomService) GetLocation(ctx context.Context, id string) (*domain.Location, error) {
	mainModel, err := svc.store.User.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	if mainModel.Location == nil {
		return nil, omniErr.ErrResourceNotFound
	}

	return svc.store.Location.GetOne(ctx, mainModel.Location.ID)

}

// AppendLocation appends Location new relation to User by id
func (svc *UserDomService) AppendLocation(ctx context.Context, id string, relationData *domain.Location) error {
	return svc.store.User.AppendLocation(ctx, id, relationData)
}

// ReplaceLocation replaces Location old relation in User by id with new Location
func (svc *UserDomService) ReplaceLocation(ctx context.Context, id string, relationData *domain.Location) error {
	return svc.store.User.ReplaceLocation(ctx, id, relationData)
}

// DeleteLocation deletes Location relation in User by id
func (svc *UserDomService) DeleteLocation(ctx context.Context, id string, relationData *domain.Location) error {
	return svc.store.User.DeleteLocation(ctx, id, relationData)
}

// GetRoles gets User's Roles
func (svc *UserDomService) GetRoles(ctx context.Context, id string) ([]*domain.Role, error) {
	mainModel, err := svc.store.User.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	if mainModel.Roles == nil {
		return nil, omniErr.ErrResourceNotFound
	}

	filters, err := filters.TransformModelsIDToFilters(mainModel.Roles)
	if err != nil {
		return nil, err
	}
	return svc.store.Role.GetList(ctx, filters)

}

// AppendRoles appends Roles new relation to User by id
func (svc *UserDomService) AppendRoles(ctx context.Context, id string, relationData []*domain.Role) error {
	return svc.store.User.AppendRoles(ctx, id, relationData)
}

// ReplaceRoles replaces Roles old relation in User by id with new Roles
func (svc *UserDomService) ReplaceRoles(ctx context.Context, id string, relationData []*domain.Role) error {
	return svc.store.User.ReplaceRoles(ctx, id, relationData)
}

// DeleteRoles deletes Roles relation in User by id
func (svc *UserDomService) DeleteRoles(ctx context.Context, id string, relationData []*domain.Role) error {
	return svc.store.User.DeleteRoles(ctx, id, relationData)
}

// GetSubscriptions gets User's Subscriptions
func (svc *UserDomService) GetSubscriptions(ctx context.Context, id string) ([]*domain.Subscription, error) {
	mainModel, err := svc.store.User.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	if mainModel.Subscriptions == nil {
		return nil, omniErr.ErrResourceNotFound
	}

	filters, err := filters.TransformModelsIDToFilters(mainModel.Subscriptions)
	if err != nil {
		return nil, err
	}
	return svc.store.Subscription.GetList(ctx, filters)

}

// AppendSubscriptions appends Subscriptions new relation to User by id
func (svc *UserDomService) AppendSubscriptions(ctx context.Context, id string, relationData []*domain.Subscription) error {
	return svc.store.User.AppendSubscriptions(ctx, id, relationData)
}

// ReplaceSubscriptions replaces Subscriptions old relation in User by id with new Subscriptions
func (svc *UserDomService) ReplaceSubscriptions(ctx context.Context, id string, relationData []*domain.Subscription) error {
	return svc.store.User.ReplaceSubscriptions(ctx, id, relationData)
}

// DeleteSubscriptions deletes Subscriptions relation in User by id
func (svc *UserDomService) DeleteSubscriptions(ctx context.Context, id string, relationData []*domain.Subscription) error {
	return svc.store.User.DeleteSubscriptions(ctx, id, relationData)
}
