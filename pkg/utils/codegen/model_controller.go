package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

// go generate ./pkg/model/domain/.

var codeController = `
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
	"strings"
)

type {{.Name}}Controller struct {
	manager *service.Manager
}

func New{{.Name}}Controller(manager *service.Manager) *{{.Name}}Controller {
	return &{{.Name}}Controller{manager: manager}
}

// Init initializes routes and swag doc
func (ctr *{{.Name}}Controller) Init(g echoswagger.ApiGroup) error {
	g.SetDescription("Operations about {{.Name}}")

	outModel, err := converter.GetExampleModelSwagOutput(new(domain.{{.Name}}))
	if err != nil {
		return err
	}

	outModelList, err := converter.GetExampleModelListSwagOutput(new(domain.{{.Name}}))
	if err != nil {
		return err
	}
	
	g.GET("/:" + strings.ToLower("{{.Name}}") + "_id", ctr.GetOne).
		{{if .CompanyResource}}AddParamPath("", "company_id", "Company ID").{{end}}
		AddParamPath("", strings.ToLower("{{.Name}}") + "_id", "{{.Name}} ID").
		AddResponse(http.StatusOK, "{{.Name}} in JSON:API format", &outModel, nil).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Gets one {{.Name}} by id")

	g.GET("", ctr.GetList).
		{{if .CompanyResource}}AddParamPath("", "company_id", "Company ID").{{end}}
		AddParamQuery("", "filter", "filter schema: filter[relation.relation_field][operator]=value", false).
		AddResponse(http.StatusOK, "{{.Name}}s list in JSON:API format", &outModelList, nil).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Gets {{.Name}}s list")

	g.POST("/", ctr.Create).
		AddParamBody(outModel, "body", "New user object", true).
		SetRequestContentType(jsonapi.MediaType).
		AddResponse(http.StatusOK, "User in JSON:API format", &outModel, nil).
		SetSummary("Creates user")

	g.PATCH("/:" + strings.ToLower("{{.Name}}") + "_id", ctr.Update).
		AddParamBody(outModel, "body", "Updates {{.Name}} object", true).
		SetRequestContentType(jsonapi.MediaType).
		AddResponse(http.StatusOK, "User in JSON:API format", &outModel, nil).
		SetSummary("Updates {{.Name}}")

	g.DELETE("/:" + strings.ToLower("{{.Name}}") + "_id", ctr.Delete).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Deletes user by id")
		
	{{if .Relations }}
	// relations
	g.GET("/:" + strings.ToLower("{{.Name}}") + "_id/relationships/:relation_name", ctr.GetRelation).
		AddParamQuery("", "relation_name", "{{.Name}} relations. Possible values: {{range .Relations}}{{.WebName}} {{end}}", true).
		SetSummary("Gets {{.Name}} relation")	
	
	g.POST("/:" + strings.ToLower("{{.Name}}") + "_id/relationships/:relation_name", ctr.ModifyRelation).
		AddParamQuery("", "relation_name", "{{.Name}} relations. Possible values: {{range .Relations}}{{.WebName}} {{end}}", true).
		SetSummary("Adds user relation")

	g.PATCH("/:" + strings.ToLower("{{.Name}}") + "_id/relationships/:relation_name", ctr.ModifyRelation).
		AddParamQuery("", "relation_name", "{{.Name}} relations. Possible values: {{range .Relations}}{{.WebName}} {{end}}", true).
		SetSummary("Replace user relation")

	g.DELETE("/:" + strings.ToLower("{{.Name}}") + "_id/relationships/:relation_name", ctr.ModifyRelation).
		AddParamQuery("", "relation_name", "{{.Name}} relations. Possible values: {{range .Relations}}{{.WebName}} {{end}}", true).
		SetSummary("Deletes user relation")
	{{end}}

	return nil
}

// GetOne returns {{.Name}}
func (ctr *{{.Name}}Controller) GetOne(ctx echo.Context) error {

	model, err := ctr.manager.{{.Name}}.GetOne(ctx.Request().Context(), ctx.Param(ctx.Param(strings.ToLower("{{.Name}}") + "_id")))
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
func (ctr *{{.Name}}Controller) GetList(ctx echo.Context) error {

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

// Create creates {{.Name}}
func (ctr *{{.Name}}Controller) Create(ctx echo.Context) error {
	domModel := new(domain.{{.Name}})
	err := httpUtils.UnmarshalFromRequest(domModel, ctx.Request().Body)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	err = validator.ValidateStruct(domModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusUnprocessableEntity, omniErr.ErrTitleValidation, err)
	}

	newModel, err := ctr.manager.{{.Name}}.Create(ctx.Request().Context(), domModel)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceExists:
			return omniErr.NewHTTPError(http.StatusConflict, omniErr.ErrTitleResourceExists, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, newModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}

	return nil
}

// Update updates {{.Name}} attributes
func (ctr *{{.Name}}Controller) Update(ctx echo.Context) error {

	domModel := new(domain.{{.Name}})
	err := httpUtils.UnmarshalFromRequest(domModel, ctx.Request().Body)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	err = validator.ValidateStruct(domModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusUnprocessableEntity, omniErr.ErrTitleValidation, err)
	}

	newModel, err := ctr.manager.{{.Name}}.Update(ctx.Request().Context(), domModel)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceExists:
			return omniErr.NewHTTPError(http.StatusConflict, omniErr.ErrTitleResourceExists, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, newModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}

	return nil

}

// Delete deletes {{.Name}}
func (ctr *{{.Name}}Controller) Delete(ctx echo.Context) error {
	err := ctr.manager.{{.Name}}.Delete(ctx.Request().Context(), ctx.Param(strings.ToLower("{{.Name}}") + "_id"))
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceExists:
			return omniErr.NewHTTPError(http.StatusConflict, omniErr.ErrTitleResourceExists, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	ctx.NoContent(http.StatusNoContent)
	return nil
}

{{if .Relations}}
// GetRelation returns relation data
func (ctr *{{.Name}}Controller) GetRelation(ctx echo.Context) error {
	webRelName := ctx.Param("relation_name")
	modelID := ctx.Param(strings.ToLower("{{.Name}}") + "_id")

	var ( 
		webData interface{}
		err error
	)	
	switch webRelName {
	{{range .Relations}}
		case "{{.WebName}}":
			webData, err = ctr.manager.{{$.Name}}.Get{{.Name}}(ctx.Request().Context(), modelID )
	{{end}}
	default:
		return omniErr.NewHTTPError(http.StatusForbidden, omniErr.ErrTitleResourceNotFound,
			fmt.Errorf("%w wrong relation name '%v'", omniErr.ErrResourceNotFound, webRelName))
	}
	
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, webData)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}

	return nil
}

// ModifyRelation modifies relation data
func (ctr *{{.Name}}Controller) ModifyRelation(ctx echo.Context) error {
	webRelName := ctx.Param("relation_name")
	modelID := ctx.Param(strings.ToLower("{{.Name}}") + "_id")

	var err error
	
	switch webRelName {
	{{range .Relations}}
		case "{{.WebName}}":
			{{if .Multiple}}
			domRecordsIntf, err := httpUtils.UnmarshalManyFromRequest(new(domain.{{.NameSingle}}), ctx.Request().Body)
			if err != nil {
				return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleBadRequest, err)
			}
	
			var domRelModel []*domain.{{.NameSingle}}
			err = converter.SliceI2SliceModel(domRecordsIntf, &domRelModel)
			if err != nil {
				return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
			}
			{{else}}
			domRelModel := new(domain.{{.Name}})
			err := httpUtils.UnmarshalFromRequest(domRelModel, ctx.Request().Body)
			if err != nil {
				return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleBadRequest, err)
			}
			{{end}}
			switch ctx.Request().Method {
			case http.MethodPatch:
				err = ctr.manager.{{$.Name}}.Replace{{.Name}}(ctx.Request().Context(), modelID, domRelModel)
			case http.MethodPost:
				err = ctr.manager.{{$.Name}}.Append{{.Name}}(ctx.Request().Context(), modelID, domRelModel)
			case http.MethodDelete:
				err = ctr.manager.{{$.Name}}.Delete{{.Name}}(ctx.Request().Context(), modelID, domRelModel)
			}
	{{end}}
	default:
		return omniErr.NewHTTPError(http.StatusForbidden, omniErr.ErrTitleResourceNotFound,
			fmt.Errorf("%w wrong relation name '%v'", omniErr.ErrResourceNotFound, webRelName))
	}
	
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}
	
	ctx.NoContent(http.StatusNoContent)
	return nil	
}
{{end}}
`

func runModelControllerGenerator(cwd string, fileIn string, modelName string, toStdOut bool, companyResource bool) error {
	fmt.Printf("Generating CONTROLLER model=%v, file=%v, wd=%v \n", modelName, fileIn, cwd)

	// get description
	modelDesc, err := getModelDescription(modelName, cwd+`\`+fileIn, companyResource)
	if err != nil {
		return err
	}

	// get template
	t, err := template.New("dummy").Parse(codeController)
	if err != nil {
		return err
	}

	var newFileName string
	if !toStdOut {
		// create destination file
		newFileName = "../../../internal/controller/" + strings.ToLower(modelName) + "_c.go"
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
