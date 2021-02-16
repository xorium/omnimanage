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
}

type Locations interface {
	GetOne(ctx context.Context, id int) (*model.Location, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*model.Location, error)
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
	//Companies Companies
	//Roles     Roles
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		Users:     NewUserRepo(db),
		Locations: NewLocationRepo(db),
	}
}
