package store

import (
	"context"
	"gorm.io/gorm"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/model/src"
)

type Users interface {
	GetOne(ctx context.Context, id int) (*src.User, error)
	GetList(ctx context.Context, f []*filters.Filter) (src.Users, error)
	Create(ctx context.Context, modelIn *src.User) (*src.User, error)
	Update(ctx context.Context, modelIn *src.User) (*src.User, error)
	Delete(ctx context.Context, id int) error
	ReplaceRelation(ctx context.Context, id int, relationName string, relationData interface{}) error
	AppendRelation(ctx context.Context, id int, relationName string, relationData interface{}) error
	DeleteRelation(ctx context.Context, id int, relationName string, relationData interface{}) error
}

type Locations interface {
	GetOne(ctx context.Context, id int) (*src.Location, error)
	GetList(ctx context.Context, f []*filters.Filter) (src.Locations, error)
	Create(ctx context.Context, modelIn *src.Location) (*src.Location, error)
	Update(ctx context.Context, modelIn *src.Location) (*src.Location, error)
	Delete(ctx context.Context, id int) error
}

type Roles interface {
	GetOne(ctx context.Context, id int) (*src.Role, error)
	GetList(ctx context.Context, f []*filters.Filter) (src.Roles, error)
	Create(ctx context.Context, modelIn *src.Role) (*src.Role, error)
	Update(ctx context.Context, modelIn *src.Role) (*src.Role, error)
	Delete(ctx context.Context, id int) error
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
