package web

type Location struct {
	ID       string                 `jsonapi:"primary,locations"`
	Name     string                 `jsonapi:"attr,name"`
	Timezone string                 `jsonapi:"attr,timezone"`
	Info     map[string]interface{} `jsonapi:"attr,info"`
	Company  *Company               `jsonapi:"relation,company"`
	Children []*Location            `jsonapi:"relation,children"`
	Users    []*User                `jsonapi:"relation,users"`
}
