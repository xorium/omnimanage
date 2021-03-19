package service

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"omnimanage/internal/store"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/model/domain"
)

const (
	OperationRelationReplace int = 1
	OperationRelationAppend  int = 2
	OperationRelationDelete  int = 3
)

func GetRelationOperFromHTTPMethod(method string) int {
	switch method {
	case http.MethodPost:
		return OperationRelationReplace
	case http.MethodPatch:
		return OperationRelationAppend
	case http.MethodDelete:
		return OperationRelationDelete
	}
	return 0
}

type CompanyServiceI interface {
	GetOne(ctx context.Context, id string) (*domain.Company, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Company, error)
	Create(ctx context.Context, modelIn *domain.Company) (*domain.Company, error)
	Update(ctx context.Context, modelIn *domain.Company) (*domain.Company, error)
	Delete(ctx context.Context, id string) error
}

type UserServiceI interface {
	GetOne(ctx context.Context, id string) (*domain.User, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.User, error)
	Create(ctx context.Context, modelIn *domain.User) (*domain.User, error)
	Update(ctx context.Context, modelIn *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id string) error

	GetCompany(ctx context.Context, id string) (*domain.Company, error)
	AppendCompany(ctx context.Context, id string, relationData *domain.Company) error
	ReplaceCompany(ctx context.Context, id string, relationData *domain.Company) error
	DeleteCompany(ctx context.Context, id string, relationData *domain.Company) error

	GetLocation(ctx context.Context, id string) (*domain.Location, error)
	AppendLocation(ctx context.Context, id string, relationData *domain.Location) error
	ReplaceLocation(ctx context.Context, id string, relationData *domain.Location) error
	DeleteLocation(ctx context.Context, id string, relationData *domain.Location) error

	GetRoles(ctx context.Context, id string) ([]*domain.Role, error)
	AppendRoles(ctx context.Context, id string, relationData []*domain.Role) error
	ReplaceRoles(ctx context.Context, id string, relationData []*domain.Role) error
	DeleteRoles(ctx context.Context, id string, relationData []*domain.Role) error

	GetSubscriptions(ctx context.Context, id string) ([]*domain.Subscription, error)
	AppendSubscriptions(ctx context.Context, id string, relationData []*domain.Subscription) error
	ReplaceSubscriptions(ctx context.Context, id string, relationData []*domain.Subscription) error
	DeleteSubscriptions(ctx context.Context, id string, relationData []*domain.Subscription) error
}

type RoleServiceI interface {
	GetOne(ctx context.Context, id string) (*domain.Role, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Role, error)
	Create(ctx context.Context, modelIn *domain.Role) (*domain.Role, error)
	Update(ctx context.Context, modelIn *domain.Role) (*domain.Role, error)
	Delete(ctx context.Context, id string) error

	GetCompany(ctx context.Context, id string) (*domain.Company, error)
	AppendCompany(ctx context.Context, id string, relationData *domain.Company) error
	ReplaceCompany(ctx context.Context, id string, relationData *domain.Company) error
	DeleteCompany(ctx context.Context, id string, relationData *domain.Company) error
}

type LocationServiceI interface {
	GetOne(ctx context.Context, id string) (*domain.Location, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Location, error)
	Create(ctx context.Context, modelIn *domain.Location) (*domain.Location, error)
	Update(ctx context.Context, modelIn *domain.Location) (*domain.Location, error)
	Delete(ctx context.Context, id string) error

	GetCompany(ctx context.Context, id string) (*domain.Company, error)
	AppendCompany(ctx context.Context, id string, relationData *domain.Company) error
	ReplaceCompany(ctx context.Context, id string, relationData *domain.Company) error
	DeleteCompany(ctx context.Context, id string, relationData *domain.Company) error

	GetChildren(ctx context.Context, id string) ([]*domain.Location, error)
	AppendChildren(ctx context.Context, id string, relationData []*domain.Location) error
	ReplaceChildren(ctx context.Context, id string, relationData []*domain.Location) error
	DeleteChildren(ctx context.Context, id string, relationData []*domain.Location) error

	GetUsers(ctx context.Context, id string) ([]*domain.User, error)
	AppendUsers(ctx context.Context, id string, relationData []*domain.User) error
	ReplaceUsers(ctx context.Context, id string, relationData []*domain.User) error
	DeleteUsers(ctx context.Context, id string, relationData []*domain.User) error
}

type Manager struct {
	Company  CompanyServiceI
	User     UserServiceI
	Location LocationServiceI
	Role     RoleServiceI
}

func NewManager(store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, errors.New("No store provided")
	}

	return &Manager{
		Company:  NewCompanyService(store),
		User:     NewUserService(store),
		Location: NewLocationService(store),
		Role:     NewRoleService(store),
	}, nil
}
