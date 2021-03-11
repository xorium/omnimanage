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

type UserService interface {
	GetOne(ctx context.Context, id string) (*domain.User, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.User, error)
	Create(ctx context.Context, modelIn *domain.User) (*domain.User, error)
	Update(ctx context.Context, modelIn *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	GetRelation(ctx context.Context, id string, relationName string) (interface{}, error)
	ModifyRelation(ctx context.Context, id string, relationName string, operation int, relationData interface{}) error
}

type LocationService interface {
	GetOne(ctx context.Context, id int) (*domain.Location, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Location, error)
	Create(ctx context.Context, modelIn *domain.Location) (*domain.Location, error)
	Update(ctx context.Context, modelIn *domain.Location) (*domain.Location, error)
	Delete(ctx context.Context, id int) error
}

type RoleService interface {
	GetOne(ctx context.Context, id int) (*domain.Role, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Role, error)
	Create(ctx context.Context, modelIn *domain.Role) (*domain.Role, error)
	Update(ctx context.Context, modelIn *domain.Role) (*domain.Role, error)
	Delete(ctx context.Context, id int) error
}

type Manager struct {
	User     UserService
	Location LocationService
	Role     RoleService
}

func NewManager(store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, errors.New("No store provided")
	}

	return &Manager{
		User: NewUserService(store),
	}, nil
}
