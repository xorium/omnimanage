package main

import (
	"fmt"
	"os"
	"text/template"
)

var codeServiceInterface = `
type {{.Name}}ServiceI interface {
	GetOne(ctx context.Context, id string) (*domain.{{.Name}}, error)
	GetList(ctx context.Context, f []*filters.Filter) ([]*domain.{{.Name}}, error)
	Create(ctx context.Context, modelIn *domain.{{.Name}}) (*domain.{{.Name}}, error)
	Update(ctx context.Context, modelIn *domain.{{.Name}}) (*domain.{{.Name}}, error)
	Delete(ctx context.Context, id string) error
{{range .Relations}}
	Get{{.Name}}(ctx context.Context, id string) ({{if .Multiple}} []*domain.{{.TypeName}} {{else}}*domain.{{.TypeName}}{{end}}, error) 
	Append{{.Name}}(ctx context.Context, id string, relationData {{if .Multiple}} []*domain.{{.TypeName}} {{else}}*domain.{{.TypeName}}{{end}}) error
	Replace{{.Name}}(ctx context.Context, id string, relationData {{if .Multiple}} []*domain.{{.TypeName}} {{else}}*domain.{{.TypeName}}{{end}}) error
	Delete{{.Name}}(ctx context.Context, id string, relationData {{if .Multiple}} []*domain.{{.TypeName}} {{else}}*domain.{{.TypeName}}{{end}}) error
{{end}}
}

`

func runModelServiceInterfaceGenerator(cwd string, fileIn string, modelName string, companyResource bool) error {
	fmt.Printf("Generating SERVICE INTERFACE model=%v, file=%v, wd=%v \n", modelName, fileIn, cwd)

	// get description
	modelDesc, err := getModelDescription(modelName, cwd+`\`+fileIn, companyResource)
	if err != nil {
		return err
	}

	// get template
	t, err := template.New("dummy").Parse(codeServiceInterface)
	if err != nil {
		return err
	}

	// generate code
	err = t.Execute(os.Stdout, modelDesc)
	if err != nil {
		return fmt.Errorf("Executing template error %v. Model = %v", err, modelName)
	}
	fmt.Printf("success! model = %v \n", modelName)

	return nil
}
