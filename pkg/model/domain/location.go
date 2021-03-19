package domain

//go:generate go run $PWD/pkg/utils/codegen/ --mode model-controller --model Location --file_in $GOFILE --company_resource true
//go:generate go run $PWD/pkg/utils/codegen/ --mode model-store --model Location --file_in $GOFILE --company_resource true
//go:generate go run $PWD/pkg/utils/codegen/ --mode model-store-interface --model Location --file_in $GOFILE --company_resource true
//go:generate go run $PWD/pkg/utils/codegen/ --mode model-service --model Location --file_in $GOFILE --company_resource true
//go:generate go run $PWD/pkg/utils/codegen/ --mode model-service-interface --model Location --file_in $GOFILE --company_resource true

type Location struct {
	ID       string                 `jsonapi:"primary,locations" default:"123"`
	Name     string                 `jsonapi:"attr,name" default:"LocationName"`
	Timezone string                 `jsonapi:"attr,timezone" default:"+00:00"`
	Info     map[string]interface{} `jsonapi:"attr,info" default:"{\"Foo\": 123}"`
	Company  *Company               `jsonapi:"relation,company" default:"{\"ID\": \"123\"}"`
	Children []*Location            `jsonapi:"relation,children"`
	Users    []*User                `jsonapi:"relation,users"` //default:"[{\"ID\": \"123\"}]"`
}
