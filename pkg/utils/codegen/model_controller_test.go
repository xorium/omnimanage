package main

import "testing"

func TestGetModelDescription(t *testing.T) {
	_, _ = getModelDescription("User", `C:\_projects\netcube\omninanage\pkg\model\domain\user.go`)
}
