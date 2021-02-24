package store

import (
	"context"
	"gorm.io/gorm"
	"omnimanage/internal/model"
	"omnimanage/pkg/filters"
)

type Users interface {
	GetOne(ctx context.Context, id int) (*model.User, error)
	GetList(ctx context.Context, f []*filters.Filter) (model.Users, error)
	Create(ctx context.Context, modelIn *model.User) (*model.User, error)
	Update(ctx context.Context, modelIn *model.User) (*model.User, error)
}

type Locations interface {
	GetOne(ctx context.Context, id int) (*model.Location, error)
	GetList(ctx context.Context, f []*filters.Filter) (model.Locations, error)
	Create(ctx context.Context, modelIn *model.Location) (*model.Location, error)
	Update(ctx context.Context, modelIn *model.Location) (*model.Location, error)
}

type Roles interface {
	GetOne(ctx context.Context, id int) (*model.Role, error)
	GetList(ctx context.Context, f []*filters.Filter) (model.Roles, error)
	Create(ctx context.Context, modelIn *model.Role) (*model.Role, error)
	Update(ctx context.Context, modelIn *model.Role) (*model.Role, error)
}

//type Companies interface {
//	GetById(ctx context.Context, id int) (*model.Company, error)
//}
//
//type Roles interface {
//	GetById(ctx context.Context, id int) (*model.Role, error)
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
