package store

import (
	"context"
	"gorm.io/gorm"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/model/domain"
)

//type Users interface {
//	GetOne(ctx context.Context, id string) (*domain.User, error)
//	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.User, error)
//	Create(ctx context.Context, modelIn *domain.User) (*domain.User, error)
//	Update(ctx context.Context, modelIn *domain.User) (*domain.User, error)
//	Delete(ctx context.Context, id string) error
//	ModifyRelation(ctx context.Context, id string, relationName string, operation int, relationData interface{}) error
//	//ReplaceRelation(ctx context.Context, id string, relationName string, relationData interface{}) error
//	//AppendRelation(ctx context.Context, id string, relationName string, relationData interface{}) error
//	//DeleteRelation(ctx context.Context, id string, relationName string, relationData interface{}) error
//}

type CompanyStoreI interface {
	GetOne(ctx context.Context, id string) (*domain.Company, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Company, error)
	Create(ctx context.Context, modelIn *domain.Company) (*domain.Company, error)
	Update(ctx context.Context, modelIn *domain.Company) (*domain.Company, error)
	Delete(ctx context.Context, id string) error
}

type UserStoreI interface {
	GetOne(ctx context.Context, id string) (*domain.User, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.User, error)
	Create(ctx context.Context, modelIn *domain.User) (*domain.User, error)
	Update(ctx context.Context, modelIn *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id string) error

	AppendCompany(ctx context.Context, id string, relationData *domain.Company) error
	ReplaceCompany(ctx context.Context, id string, relationData *domain.Company) error
	DeleteCompany(ctx context.Context, id string, relationData *domain.Company) error

	AppendLocation(ctx context.Context, id string, relationData *domain.Location) error
	ReplaceLocation(ctx context.Context, id string, relationData *domain.Location) error
	DeleteLocation(ctx context.Context, id string, relationData *domain.Location) error

	AppendRoles(ctx context.Context, id string, relationData []*domain.Role) error
	ReplaceRoles(ctx context.Context, id string, relationData []*domain.Role) error
	DeleteRoles(ctx context.Context, id string, relationData []*domain.Role) error

	AppendSubscriptions(ctx context.Context, id string, relationData []*domain.Subscription) error
	ReplaceSubscriptions(ctx context.Context, id string, relationData []*domain.Subscription) error
	DeleteSubscriptions(ctx context.Context, id string, relationData []*domain.Subscription) error
}

type LocationStoreI interface {
	GetOne(ctx context.Context, id string) (*domain.Location, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Location, error)
	Create(ctx context.Context, modelIn *domain.Location) (*domain.Location, error)
	Update(ctx context.Context, modelIn *domain.Location) (*domain.Location, error)
	Delete(ctx context.Context, id string) error

	AppendCompany(ctx context.Context, id string, relationData *domain.Company) error
	ReplaceCompany(ctx context.Context, id string, relationData *domain.Company) error
	DeleteCompany(ctx context.Context, id string, relationData *domain.Company) error

	AppendChildren(ctx context.Context, id string, relationData []*domain.Location) error
	ReplaceChildren(ctx context.Context, id string, relationData []*domain.Location) error
	DeleteChildren(ctx context.Context, id string, relationData []*domain.Location) error

	AppendUsers(ctx context.Context, id string, relationData []*domain.User) error
	ReplaceUsers(ctx context.Context, id string, relationData []*domain.User) error
	DeleteUsers(ctx context.Context, id string, relationData []*domain.User) error
}

type RoleStoreI interface {
	GetOne(ctx context.Context, id string) (*domain.Role, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Role, error)
	Create(ctx context.Context, modelIn *domain.Role) (*domain.Role, error)
	Update(ctx context.Context, modelIn *domain.Role) (*domain.Role, error)
	Delete(ctx context.Context, id string) error

	AppendCompany(ctx context.Context, id string, relationData *domain.Company) error
	ReplaceCompany(ctx context.Context, id string, relationData *domain.Company) error
	DeleteCompany(ctx context.Context, id string, relationData *domain.Company) error
}

//type Locations interface {
//	GetOne(ctx context.Context, id string) (*domain.Location, error)
//	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Location, error)
//	Create(ctx context.Context, modelIn *domain.Location) (*domain.Location, error)
//	Update(ctx context.Context, modelIn *domain.Location) (*domain.Location, error)
//	Delete(ctx context.Context, id string) error
//}
//
//type Roles interface {
//	GetOne(ctx context.Context, id string) (*domain.Role, error)
//	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.Role, error)
//	Create(ctx context.Context, modelIn *domain.Role) (*domain.Role, error)
//	Update(ctx context.Context, modelIn *domain.Role) (*domain.Role, error)
//	Delete(ctx context.Context, id string) error
//}

// Store contains all repositories
type Store struct {
	Company  CompanyStoreI
	User     UserStoreI
	Location LocationStoreI
	Role     RoleStoreI
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		User:     NewUserRepo(db),
		Location: NewLocationRepo(db),
		Role:     NewRoleRepo(db),
	}
}
