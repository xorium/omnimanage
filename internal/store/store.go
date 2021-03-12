package store

import (
	"context"
	"gorm.io/gorm"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/model/domain"
)

type Users interface {
	GetOne(ctx context.Context, id string) (*domain.User, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.User, error)
	Create(ctx context.Context, modelIn *domain.User) (*domain.User, error)
	Update(ctx context.Context, modelIn *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	ModifyRelation(ctx context.Context, id string, relationName string, operation int, relationData interface{}) error
	//ReplaceRelation(ctx context.Context, id string, relationName string, relationData interface{}) error
	//AppendRelation(ctx context.Context, id string, relationName string, relationData interface{}) error
	//DeleteRelation(ctx context.Context, id string, relationName string, relationData interface{}) error
}

type Locations interface {
	GetOne(ctx context.Context, id string) (*domain.Location, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Location, error)
	Create(ctx context.Context, modelIn *domain.Location) (*domain.Location, error)
	Update(ctx context.Context, modelIn *domain.Location) (*domain.Location, error)
	Delete(ctx context.Context, id string) error
}

type Roles interface {
	GetOne(ctx context.Context, id string) (*domain.Role, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Role, error)
	Create(ctx context.Context, modelIn *domain.Role) (*domain.Role, error)
	Update(ctx context.Context, modelIn *domain.Role) (*domain.Role, error)
	Delete(ctx context.Context, id string) error
}

//type Companies interface {
//	GetById(ctx context.Context, id int) (*src.Company, error)
//}
//
//type Roles interface {
//	GetById(ctx context.Context, id int) (*src.Role, error)
//}

// Store contains all repositories
type Store struct {
	Users     Users
	Locations Locations
	Roles     Roles
	//Companies Companies

}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		Users:     NewUserRepo(db),
		Locations: NewLocationRepo(db),
		Roles:     NewRoleRepo(db),
	}
}
