package service

import (
	"context"
	"omnimanage/internal/store"
	omniErr "omnimanage/pkg/error"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/model/domain"
)

type LocationDomService struct {
	store *store.Store
}

func NewLocationService(store *store.Store) *LocationDomService {
	return &LocationDomService{
		store: store,
	}
}

// GetOne gets one Location by ID
func (svc *LocationDomService) GetOne(ctx context.Context, id string) (*domain.Location, error) {
	return svc.store.Location.GetOne(ctx, id)
}

// GetList gets Locations list with optional filters
func (svc *LocationDomService) GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Location, error) {
	return svc.store.Location.GetList(ctx, f)
}

// Create creates new Location
func (svc *LocationDomService) Create(ctx context.Context, modelIn *domain.Location) (*domain.Location, error) {
	return svc.store.Location.Create(ctx, modelIn)
}

// Update updates Location
func (svc *LocationDomService) Update(ctx context.Context, modelIn *domain.Location) (*domain.Location, error) {
	return svc.store.Location.Update(ctx, modelIn)
}

// Delete deletes Location
func (svc *LocationDomService) Delete(ctx context.Context, id string) error {
	return svc.store.Location.Delete(ctx, id)
}

// GetCompany gets Location's Company
func (svc *LocationDomService) GetCompany(ctx context.Context, id string) (*domain.Company, error) {
	mainModel, err := svc.store.Location.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	if mainModel.Company == nil {
		return nil, omniErr.ErrResourceNotFound
	}

	return svc.store.Company.GetOne(ctx, mainModel.Company.ID)

}

// AppendCompany appends Company new relation to Location by id
func (svc *LocationDomService) AppendCompany(ctx context.Context, id string, relationData *domain.Company) error {
	return svc.store.Location.AppendCompany(ctx, id, relationData)
}

// ReplaceCompany replaces Company old relation in Location by id with new Company
func (svc *LocationDomService) ReplaceCompany(ctx context.Context, id string, relationData *domain.Company) error {
	return svc.store.Location.ReplaceCompany(ctx, id, relationData)
}

// DeleteCompany deletes Company relation in Location by id
func (svc *LocationDomService) DeleteCompany(ctx context.Context, id string, relationData *domain.Company) error {
	return svc.store.Location.DeleteCompany(ctx, id, relationData)
}

// GetChildren gets Location's Children
func (svc *LocationDomService) GetChildren(ctx context.Context, id string) ([]*domain.Location, error) {
	mainModel, err := svc.store.Location.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	if mainModel.Children == nil {
		return nil, omniErr.ErrResourceNotFound
	}

	filters, err := filters.TransformModelsIDToFilters(mainModel.Children)
	if err != nil {
		return nil, err
	}
	return svc.store.Location.GetList(ctx, filters)

}

// AppendChildren appends Children new relation to Location by id
func (svc *LocationDomService) AppendChildren(ctx context.Context, id string, relationData []*domain.Location) error {
	return svc.store.Location.AppendChildren(ctx, id, relationData)
}

// ReplaceChildren replaces Children old relation in Location by id with new Children
func (svc *LocationDomService) ReplaceChildren(ctx context.Context, id string, relationData []*domain.Location) error {
	return svc.store.Location.ReplaceChildren(ctx, id, relationData)
}

// DeleteChildren deletes Children relation in Location by id
func (svc *LocationDomService) DeleteChildren(ctx context.Context, id string, relationData []*domain.Location) error {
	return svc.store.Location.DeleteChildren(ctx, id, relationData)
}

// GetUsers gets Location's Users
func (svc *LocationDomService) GetUsers(ctx context.Context, id string) ([]*domain.User, error) {
	mainModel, err := svc.store.Location.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	if mainModel.Users == nil {
		return nil, omniErr.ErrResourceNotFound
	}

	filters, err := filters.TransformModelsIDToFilters(mainModel.Users)
	if err != nil {
		return nil, err
	}
	return svc.store.User.GetList(ctx, filters)

}

// AppendUsers appends Users new relation to Location by id
func (svc *LocationDomService) AppendUsers(ctx context.Context, id string, relationData []*domain.User) error {
	return svc.store.Location.AppendUsers(ctx, id, relationData)
}

// ReplaceUsers replaces Users old relation in Location by id with new Users
func (svc *LocationDomService) ReplaceUsers(ctx context.Context, id string, relationData []*domain.User) error {
	return svc.store.Location.ReplaceUsers(ctx, id, relationData)
}

// DeleteUsers deletes Users relation in Location by id
func (svc *LocationDomService) DeleteUsers(ctx context.Context, id string, relationData []*domain.User) error {
	return svc.store.Location.DeleteUsers(ctx, id, relationData)
}
