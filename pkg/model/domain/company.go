package domain

//go:generate go run $PWD/pkg/tools/codegen/ --mode model-controller --model Company --file_in $GOFILE

type Company struct {
	ID   string `jsonapi:"primary,companies" default:"1"`
	Name string `jsonapi:"attr,name" default:"CompanyName"`
}
