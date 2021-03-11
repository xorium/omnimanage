package domain

type Location struct {
	ID       string                 `jsonapi:"primary,locations" default:"123"`
	Name     string                 `jsonapi:"attr,name" default:"LocationName"`
	Timezone string                 `jsonapi:"attr,timezone" default:"+00:00"`
	Info     map[string]interface{} `jsonapi:"attr,info" default:"{\"Foo\": 123}"`
	Company  *Company               `jsonapi:"relation,company" default:"{\"ID\": \"123\"}"`
	Children []*Location            `jsonapi:"relation,children"`
	Users    []*User                `jsonapi:"relation,users"` //default:"[{\"ID\": \"123\"}]"`
}
