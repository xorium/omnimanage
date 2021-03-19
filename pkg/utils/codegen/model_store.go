package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

var codeStore = `
package store

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	omniErr "omnimanage/pkg/error"
	"omnimanage/pkg/filters"
	"omnimanage/pkg/mapper"
	"omnimanage/pkg/model/domain"
	"omnimanage/pkg/model/src"
)

type {{.Name}}Repo struct {
	db *gorm.DB
}

func New{{.Name}}Repo(db *gorm.DB) *{{.Name}}Repo {
	return &{{.Name}}Repo{db: db}
}

// GetOne gets one {{.Name}} by ID
func (r *{{.Name}}Repo) GetOne(ctx context.Context, id string) (*domain.{{.Name}}, error) {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.{{.Name}}{})
	if err != nil {
		return nil, err
	}

	// Gets from db
	db := r.db.WithContext(ctx)

	srcModel := new(src.{{.Name}})
	dbResult := db.Where("id = ?", idSrc).Preload(clause.Associations).First(srcModel)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, omniErr.ErrResourceNotFound
	} else if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	// src Model -> domain
	dom, err := srcModel.ToWeb()
	if err != nil {
		return nil, err
	}

	return dom, nil
}

// GetList gets {{.Name}}s list with optional filters
func (r *{{.Name}}Repo) GetList(ctx context.Context, f []*filters.Filter) ([]*domain.{{.Name}}, error) {
	// string filters -> src filters
	srcFilters, err := filters.TransformWebToSrc(f, &domain.{{.Name}}{}, &src.{{.Name}}{})
	if err != nil {
		return nil, err
	}

	// Gets from db
	srcModels := make(src.{{.Name}}s, 0, 1)
	db := r.db.WithContext(ctx)
	db, err = filters.SetGormFilters(db, &srcModels, srcFilters)
	if err != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	dbResult := db.Preload(clause.Associations).Find(&srcModels)
	if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	if dbResult.RowsAffected == 0 {
		return nil, omniErr.ErrResourceNotFound
	}

	// src Model -> domain
	domModels, err := srcModels.ToWeb()
	if err != nil {
		return nil, err
	}

	return domModels, nil
}

// Create creates new {{.Name}}
func (r *{{.Name}}Repo) Create(ctx context.Context, modelIn *domain.{{.Name}}) (*domain.{{.Name}}, error) {
	// domain Model -> src model
	srcModel, err := new(src.{{.Name}}).ScanFromWeb(modelIn)
	if err != nil {
		return nil, err
	}

	// db operations
	db := r.db.WithContext(ctx)

	// check existence
	tmpRec := new(src.{{.Name}})
	dbResult := db.Where("id = ?", srcModel.ID).First(tmpRec)
	if dbResult.Error != nil && !errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}
	if dbResult.RowsAffected > 0 {
		return nil, fmt.Errorf("%w", omniErr.ErrResourceExists)
	}

	// create data
	dbResult = db.Preload(clause.Associations).Create(&srcModel)
	if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	// src Model -> domain model
	domModel, err := srcModel.ToWeb()
	if err != nil {
		return nil, err
	}

	return domModel, nil
}

// Update updates {{.Name}}
func (r *{{.Name}}Repo) Update(ctx context.Context, modelIn *domain.{{.Name}}) (*domain.{{.Name}}, error) {
	// domain Model -> src model
	srcModel, err := new(src.{{.Name}}).ScanFromWeb(modelIn)
	if err != nil {
		return nil, err
	}

	// db operations
	db := r.db.WithContext(ctx)

	// check existence
	tmpRec := new(src.{{.Name}})
	dbResult := db.Where("id = ?", srcModel.ID).First(tmpRec)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, omniErr.ErrResourceNotFound
	} else if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	// update data
	dbResult = db.Preload(clause.Associations).Save(&srcModel)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return nil, omniErr.ErrResourceNotFound
	} else if dbResult.Error != nil {
		return nil, fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}

	// src Model -> domain model
	domModel, err := srcModel.ToWeb()
	if err != nil {
		return nil, err
	}

	return domModel, nil
}

// Delete deletes {{.Name}}
func (r *{{.Name}}Repo) Delete(ctx context.Context, id string) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.{{.Name}}{})
	if err != nil {
		return err
	}

	// db operations
	db := r.db.WithContext(ctx)
	dbResult := db.Delete(&src.{{.Name}}{}, idSrc)
	if dbResult.Error != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, dbResult.Error)
	}
	if dbResult.RowsAffected == 0 {
		return omniErr.ErrResourceNotFound
	}

	return nil
}

{{range .Relations}}
// Append{{.Name}} appends {{.Name}} new relation to {{$.Name}} by id
func (r *{{$.Name}}Repo) Append{{.Name}}(ctx context.Context, id string, relationData {{if .Multiple}} []*domain.{{.TypeName}} {{else}}*domain.{{.TypeName}}{{end}}) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.{{$.Name}}{})
	if err != nil {
		return err
	}
	{{if .Multiple}}
	srcModelNew, err := src.{{.TypeNameMulti}}.ScanFromWeb(nil, relationData)	
	{{else}}
	srcModelNew, err := new(src.{{.TypeName}}).ScanFromWeb(relationData)
	{{end}}
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.{{$.Name}}{ID: idSrc}).Association("{{.Name}}").Append(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// Replace{{.Name}} replaces {{.Name}} old relation in {{$.Name}} by id with new {{.Name}}
func (r *{{$.Name}}Repo) Replace{{.Name}}(ctx context.Context, id string, relationData {{if .Multiple}} []*domain.{{.TypeName}} {{else}}*domain.{{.TypeName}}{{end}}) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.{{$.Name}}{})
	if err != nil {
		return err
	}
	{{if .Multiple}}
	srcModelNew, err := src.{{.TypeNameMulti}}.ScanFromWeb(nil, relationData)	
	{{else}}
	srcModelNew, err := new(src.{{.TypeName}}).ScanFromWeb(relationData)
	{{end}}
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.{{$.Name}}{ID: idSrc}).Association("{{.Name}}").Replace(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}

// Delete{{.Name}} deletes {{.Name}} relation in {{$.Name}} by id
func (r *{{$.Name}}Repo) Delete{{.Name}}(ctx context.Context, id string, relationData {{if .Multiple}} []*domain.{{.TypeName}} {{else}}*domain.{{.TypeName}}{{end}}) error {
	// domain ID -> src ID
	idSrc, err := mapper.Get().GetSrcID(id, &src.{{$.Name}}{})
	if err != nil {
		return err
	}
	{{if .Multiple}}
	srcModelNew, err := src.{{.TypeNameMulti}}.ScanFromWeb(nil, relationData)	
	{{else}}
	srcModelNew, err := new(src.{{.TypeName}}).ScanFromWeb(relationData)
	{{end}}
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	db := r.db.WithContext(ctx)
	err = db.Model(&src.{{$.Name}}{ID: idSrc}).Association("{{.Name}}").Delete(srcModelNew)
	if err != nil {
		return fmt.Errorf("%w %v", omniErr.ErrInternal, err)
	}

	return nil
}
{{end}}
`

func runModelStoreGenerator(cwd string, fileIn string, modelName string, toStdOut bool, companyResource bool) error {
	fmt.Printf("Generating STORE model=%v, file=%v, wd=%v \n", modelName, fileIn, cwd)

	// get description
	modelDesc, err := getModelDescription(modelName, cwd+`\`+fileIn, companyResource)
	if err != nil {
		return err
	}

	// get template
	t, err := template.New("dummy").Parse(codeStore)
	if err != nil {
		return err
	}

	var newFileName string
	if !toStdOut {
		// create destination file
		newFileName = "../../../internal/store/" + strings.ToLower(modelName) + "_store.go"
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
