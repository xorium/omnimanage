package model

import "testing"

func TestUser_ToWeb(t *testing.T) {
	userDb := User{
		ID:       23,
		UserName: "name1",
		Company:  &Company{ID: 2, Name: "Sespel"},
	}

	userOmni, err := userDb.ToWeb()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", userOmni)
}
