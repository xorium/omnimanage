package service

import (
	"context"
	"omnimanage/internal/store"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/model/domain"
)

type CompanyDomService struct {
	store *store.Store
}

func NewCompanyService(store *store.Store) *CompanyDomService {
	return &CompanyDomService{
		store: store,
	}
}

// GetOne gets one Company by ID
func (svc *CompanyDomService) GetOne(ctx context.Context, id string) (*domain.Company, error) {
	return svc.store.Company.GetOne(ctx, id)
}

// GetList gets Companys list with optional filters
func (svc *CompanyDomService) GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Company, error) {
	return svc.store.Company.GetList(ctx, f)
}

// Create creates new Company
func (svc *CompanyDomService) Create(ctx context.Context, modelIn *domain.Company) (*domain.Company, error) {
	return svc.store.Company.Create(ctx, modelIn)
}

// Update updates Company
func (svc *CompanyDomService) Update(ctx context.Context, modelIn *domain.Company) (*domain.Company, error) {
	return svc.store.Company.Update(ctx, modelIn)
}

// Delete deletes Company
func (svc *CompanyDomService) Delete(ctx context.Context, id string) error {
	return svc.store.Company.Delete(ctx, id)
}
