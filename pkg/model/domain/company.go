package domain

//go:generate go run $PWD/pkg/utils/codegen/ --mode model-controller --model Company --file_in $GOFILE --company_resource false
//go:generate go run $PWD/pkg/utils/codegen/ --mode model-store --model Company --file_in $GOFILE --company_resource false
//go:generate go run $PWD/pkg/utils/codegen/ --mode model-store-interface --model Company --file_in $GOFILE --company_resource false
//go:generate go run $PWD/pkg/utils/codegen/ --mode model-service --model Company --file_in $GOFILE --company_resource false
//go:generate go run $PWD/pkg/utils/codegen/ --mode model-service-interface --model Company --file_in $GOFILE --company_resource false

type Company struct {
	ID   string `jsonapi:"primary,companies" default:"1"`
	Name string `jsonapi:"attr,name" default:"CompanyName"`
}
