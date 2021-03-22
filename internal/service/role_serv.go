package service

import (
	"context"
	"omnimanage/internal/store"
	omniErr "omnimanage/pkg/error"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/model/domain"
)

type RoleDomService struct {
	store *store.Store
}

func NewRoleService(store *store.Store) *RoleDomService {
	return &RoleDomService{
		store: store,
	}
}

// GetOne gets one Role by ID
func (svc *RoleDomService) GetOne(ctx context.Context, id string) (*domain.Role, error) {
	return svc.store.Role.GetOne(ctx, id)
}

// GetList gets Roles list with optional filters
func (svc *RoleDomService) GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Role, error) {
	return svc.store.Role.GetList(ctx, f)
}

// Create creates new Role
func (svc *RoleDomService) Create(ctx context.Context, modelIn *domain.Role) (*domain.Role, error) {
	return svc.store.Role.Create(ctx, modelIn)
}

// Update updates Role
func (svc *RoleDomService) Update(ctx context.Context, modelIn *domain.Role) (*domain.Role, error) {
	return svc.store.Role.Update(ctx, modelIn)
}

// Delete deletes Role
func (svc *RoleDomService) Delete(ctx context.Context, id string) error {
	return svc.store.Role.Delete(ctx, id)
}

// GetCompany gets Role's Company
func (svc *RoleDomService) GetCompany(ctx context.Context, id string) (*domain.Company, error) {
	mainModel, err := svc.store.Role.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	if mainModel.Company == nil {
		return nil, omniErr.ErrResourceNotFound
	}

	return svc.store.Company.GetOne(ctx, mainModel.Company.ID)

}

// AppendCompany appends Company new relation to Role by id
func (svc *RoleDomService) AppendCompany(ctx context.Context, id string, relationData *domain.Company) error {
	return svc.store.Role.AppendCompany(ctx, id, relationData)
}

// ReplaceCompany replaces Company old relation in Role by id with new Company
func (svc *RoleDomService) ReplaceCompany(ctx context.Context, id string, relationData *domain.Company) error {
	return svc.store.Role.ReplaceCompany(ctx, id, relationData)
}

// DeleteCompany deletes Company relation in Role by id
func (svc *RoleDomService) DeleteCompany(ctx context.Context, id string, relationData *domain.Company) error {
	return svc.store.Role.DeleteCompany(ctx, id, relationData)
}
