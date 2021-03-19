package domain

//go:generate go run $PWD/pkg/utils/codegen/ --mode model-controller --model Role --file_in $GOFILE --company_resource true
//go:generate go run $PWD/pkg/utils/codegen/ --mode model-store --model Role --file_in $GOFILE --company_resource true
//go:generate go run $PWD/pkg/utils/codegen/ --mode model-store-interface --model Role --file_in $GOFILE --company_resource true
//go:generate go run $PWD/pkg/utils/codegen/ --mode model-service --model Role --file_in $GOFILE --company_resource true
//go:generate go run $PWD/pkg/utils/codegen/ --mode model-service-interface --model Role --file_in $GOFILE --company_resource true

type Role struct {
	ID         string                 `jsonapi:"primary,roles"`
	Name       string                 `jsonapi:"attr,name"`
	Slug       string                 `jsonapi:"attr,slug"`
	Assigned   bool                   `jsonapi:"attr,assigned"`
	Persistent bool                   `jsonapi:"attr,persistent"`
	Info       map[string]interface{} `jsonapi:"attr,info"`
	Company    *Company               `jsonapi:"relation,company"`
}
