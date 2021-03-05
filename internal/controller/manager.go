package controller

import (
	"omnimanage/internal/store"
)

// Manager is just a collection of all controllers we have in the project
type Manager struct {
	User *UserController
	Role *RoleController
}

func NewManager(store *store.Store) *Manager {
	return &Manager{
		User: NewUserController(store),
		Role: NewRoleController(store),
	}
}
