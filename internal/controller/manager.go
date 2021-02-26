package controller

import (
	"omnimanage/internal/store"
	"omnimanage/pkg/mapper"
)

// Manager is just a collection of all controllers we have in the project
type Manager struct {
	User *UserController
	Role *RoleController
}

func NewManager(store *store.Store, mapper *mapper.ModelMapper) *Manager {
	return &Manager{
		User: NewUserController(store, mapper),
		Role: NewRoleController(store),
	}
}
