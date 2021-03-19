package main

import "testing"

func TestGetModelDescription(t *testing.T) {
	_, _ = getModelDescription("Location", `C:\_projects\netcube\omninanage\pkg\model\domain\location.go`, true)
}
