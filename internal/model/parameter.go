package model

type Parameter struct {
	ID                int                    `jsonapi:"primary,parameters"`
	Name              string                 `jsonapi:"attr,name"`
	Description       string                 `jsonapi:"attr,desc"`
	Type              string                 `jsonapi:"attr,type"`
	IsValuesSetFinite bool                   `jsonapi:"attr,isValuesSetFinite"`
	Info              map[string]interface{} `jsonapi:"attr,info"`
}
