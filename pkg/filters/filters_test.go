package filters

import (
	"gitlab.omnicube.ru/libs/omnilib/models"
	"testing"
)

func TestTransform(t *testing.T) {
	f := []*Filter{
		&Filter{
			Relation: "company",
			Field:    "name",
			Operator: "=",
			Value:    "aaa",
		},
		&Filter{
			Relation: "",
			Field:    "time",
			Operator: "=",
			Value:    "123123",
		},
	}

	ev := models.Event{}

	fNew, err := Transform(f, ev, nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(fNew)

}
