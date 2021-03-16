package main

import (
	"fmt"
	"html/template"
	"omnimanage/pkg/utils/model_parser"
	"os"
)

var code = `
package controller

import (
	"fmt"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/pkg/errors"
	"net/http"
	"omnimanage/internal/service"
	"omnimanage/internal/validator"
	omniErr "omnimanage/pkg/error"
	filt "omnimanage/pkg/filters"
	"omnimanage/pkg/model/domain"
	"omnimanage/pkg/utils/converter"
	httpUtils "omnimanage/pkg/utils/http"
)

type {{.Name}}Controller struct {
	manager *service.Manager
}

func New{{.Name}}Controller(manager *service.Manager) *{{.Name}}Controller {
	return &{{.Name}}Controller{manager: manager}
}

// Init initializes routes and swag doc
func (ctr *{{.Name}}Controller) Init(g echoswagger.ApiGroup) error {

}

// GetOne returns {{.Name}}
func (ctr *{{.Name}}Controller) GetOne(ctx echo.Context) error {

	model, err := ctr.manager.{{.Name}}.GetOne(ctx.Request().Context(), ctx.Param("{{.Name}}_id"))
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, model)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}

	return nil
}

// GetList returns {{.Name}}s list
func (ctr *UserController) GetList(ctx echo.Context) error {

	filterStrings, err := filt.ParseQueryString(ctx.Request().URL.RawQuery, &domain.{{.Name}}{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	models, err := ctr.manager.{{.Name}}.GetList(ctx.Request().Context(), filterStrings)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, models)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}
	return nil
}



`

type ModelDescription struct {
	Name      string
	Relations []*model_parser.Relation
}

func runModelControllerGenerator(cwd string, fileIn string, modelName string) error {
	fmt.Printf("Generating model=%v, file=%v, wd=%v \n", modelName, fileIn, cwd)

	modelDesc, err := getModelDescription(modelName, cwd+`\`+fileIn)
	if err != nil {
		return err
	}

	t, err := template.New("dummy").Parse(code)
	if err != nil {
		return err
	}

	t.Execute(os.Stdout, modelDesc)

	return nil
}

func getModelDescription(modelName string, file string) (*ModelDescription, error) {
	p, err := model_parser.NewParser(file)
	if err != nil {
		return nil, err
	}

	rels, err := p.GetRelations(modelName)
	if err != nil {
		return nil, err
	}

	m := &ModelDescription{
		Name:      modelName,
		Relations: rels,
	}

	return m, nil
}
