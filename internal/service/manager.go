package service

import (
	"fmt"
	"omnimanage/internal/store"
)

// Manager is just a collection of all services we have in the project
type Manager struct {
	User *UserService
}

func NewManager(store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, fmt.Errorf("No store provided")
	}
	return &Manager{
		User: NewUserService(store),
	}, nil
}
