package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

var codeService = `
package service

import (
	"context"
	"omnimanage/internal/store"
	omniErr "omnimanage/pkg/error"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/model/domain"
)

type {{.Name}}DomService struct {
	store *store.Store
}

func New{{.Name}}Service(store *store.Store) *{{.Name}}DomService {
	return &{{.Name}}DomService{
		store: store,
	}
}

// GetOne gets one {{.Name}} by ID
func (svc *{{.Name}}DomService) GetOne(ctx context.Context, id string) (*domain.{{.Name}}, error) {
	return svc.store.{{.Name}}.GetOne(ctx, id)
}

// GetList gets {{.Name}}s list with optional filters
func (svc *{{.Name}}DomService) GetList(ctx context.Context, f []*filters.Filter) ([]*domain.{{.Name}}, error) {
	return svc.store.{{.Name}}.GetList(ctx, f)
}

// Create creates new {{.Name}}
func (svc *{{.Name}}DomService) Create(ctx context.Context, modelIn *domain.{{.Name}}) (*domain.{{.Name}}, error) {
	return svc.store.{{.Name}}.Create(ctx, modelIn)
}

// Update updates {{.Name}}
func (svc *{{.Name}}DomService) Update(ctx context.Context, modelIn *domain.{{.Name}}) (*domain.{{.Name}}, error) {
	return svc.store.{{.Name}}.Update(ctx, modelIn)
}

// Delete deletes {{.Name}}
func (svc *{{.Name}}DomService) Delete(ctx context.Context, id string) error {
	return svc.store.{{.Name}}.Delete(ctx, id)
}

{{range .Relations}}
// Get{{.Name}} gets {{$.Name}}'s {{.Name}}
func (svc *{{$.Name}}DomService) Get{{.Name}}(ctx context.Context, id string) ({{if .Multiple}} []*domain.{{.TypeName}} {{else}}*domain.{{.TypeName}}{{end}}, error) {
	mainModel, err := svc.store.{{$.Name}}.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if mainModel.{{.Name}} == nil {
		return nil, omniErr.ErrResourceNotFound
	}
	{{if .Multiple}}
	filters, err := filters.TransformModelsIDToFilters(mainModel.{{.Name}})
	if err != nil {
		return nil, err
	}
	return svc.store.{{.TypeName}}.GetList(ctx, filters)
	{{else}}
	return svc.store.{{.Name}}.GetOne(ctx, mainModel.{{.Name}}.ID)
	{{end}}
}

// Append{{.Name}} appends {{.Name}} new relation to {{$.Name}} by id
func (svc *{{$.Name}}DomService) Append{{.Name}}(ctx context.Context, id string, relationData {{if .Multiple}} []*domain.{{.TypeName}} {{else}}*domain.{{.TypeName}}{{end}}) error {
	return svc.store.{{$.Name}}.Append{{.Name}}(ctx, id, relationData)
}

// Replace{{.Name}} replaces {{.Name}} old relation in {{$.Name}} by id with new {{.Name}}
func (svc *{{$.Name}}DomService) Replace{{.Name}}(ctx context.Context, id string, relationData {{if .Multiple}} []*domain.{{.TypeName}} {{else}}*domain.{{.TypeName}}{{end}}) error {
	return svc.store.{{$.Name}}.Replace{{.Name}}(ctx, id, relationData)
}

// Delete{{.Name}} deletes {{.Name}} relation in {{$.Name}} by id
func (svc *{{$.Name}}DomService) Delete{{.Name}}(ctx context.Context, id string, relationData {{if .Multiple}} []*domain.{{.TypeName}} {{else}}*domain.{{.TypeName}}{{end}}) error {
	return svc.store.{{$.Name}}.Delete{{.Name}}(ctx, id, relationData)
}
{{end}}
`

func runModelServiceGenerator(cwd string, fileIn string, modelName string, toStdOut bool, companyResource bool) error {
	fmt.Printf("Generating SERVICE model=%v, file=%v, wd=%v \n", modelName, fileIn, cwd)

	// get description
	modelDesc, err := getModelDescription(modelName, cwd+`\`+fileIn, companyResource)
	if err != nil {
		return err
	}

	// get template
	t, err := template.New("dummy").Parse(codeService)
	if err != nil {
		return err
	}

	var newFileName string
	if !toStdOut {
		// create destination file
		newFileName = "../../../internal/service/" + strings.ToLower(modelName) + "_serv.go"
		newFileWriter, err := os.Create(newFileName)
		if err != nil {
			return fmt.Errorf("Creating file %v : %v ", newFileName, err)
		}

		// generate code
		err = t.Execute(newFileWriter, modelDesc)

	} else {
		newFileName = "stdout"
		// generate code
		err = t.Execute(os.Stdout, modelDesc)
	}
	if err != nil {
		return fmt.Errorf("Executing template error %v. Model = %v, out file = %v", err, modelName, newFileName)
	}
	fmt.Printf("success! model = %v, out = %v \n", modelName, newFileName)

	return nil
}
