package domain

type Manufacturer struct {
	ID   string                 `jsonapi:"primary,manufacturers"`
	Name string                 `jsonapi:"attr,name"`
	Info map[string]interface{} `jsonapi:"attr,info"`
}
