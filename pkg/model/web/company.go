package web

type Company struct {
	ID   string `jsonapi:"primary,companies"`
	Name string `jsonapi:"attr,name"`
}
