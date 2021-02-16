package model

import "testing"

func TestUser_ToWeb(t *testing.T) {
	userDb := User{
		ID:       23,
		UserName: "name1",
		Settings: []byte("{}"),
		Roles: []*Role{
			&Role{ID: 1, Name: "Name1", Info: []byte("{}")},
			&Role{ID: 2, Name: "Name2", Info: []byte("{}")},
		},
	}

	userOmni, err := userDb.ToWeb()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", userOmni)
}
