package domain

type Role struct {
	ID         string                 `jsonapi:"primary,roles"`
	Name       string                 `jsonapi:"attr,name"`
	Slug       string                 `jsonapi:"attr,slug"`
	Assigned   bool                   `jsonapi:"attr,assigned"`
	Persistent bool                   `jsonapi:"attr,persistent"`
	Info       map[string]interface{} `jsonapi:"attr,info"`
	Company    *Company               `jsonapi:"relation,company"`
}
