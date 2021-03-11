package domain

type Company struct {
	ID   string `jsonapi:"primary,companies" default:"1"`
	Name string `jsonapi:"attr,name" default:"CompanyName"`
}
