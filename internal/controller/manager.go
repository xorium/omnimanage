package controller

import (
	"omnimanage/internal/service"
)

// Manager is just a collection of all controllers we have in the project
type Manager struct {
	User *UserController
	Role *RoleController
}

func NewManager(svc *service.Manager) *Manager {
	return &Manager{
		User: NewUserController(svc),
		Role: NewRoleController(svc),
	}
}
