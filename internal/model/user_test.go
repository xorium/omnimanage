package model

import "testing"

func TestUser_ToWebM(t *testing.T) {
	userDb := User{
		ID:      23,
		Name:    "name1",
		Company: &Company{ID: 2, Name: "Sespel"},
	}

	userOmni, err := userDb.ToWebM()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", userOmni)
}
