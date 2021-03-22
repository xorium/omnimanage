package controller

import (
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

type CompanyController struct {
	manager *service.Manager
}

func NewCompanyController(manager *service.Manager) *CompanyController {
	return &CompanyController{manager: manager}
}

// Init initializes routes and swag doc
func (ctr *CompanyController) Init(g echoswagger.ApiGroup) error {
	g.SetDescription("Operations about Company")

	outModel, err := converter.GetExampleModelSwagOutput(new(domain.Company))
	if err != nil {
		return err
	}

	outModelList, err := converter.GetExampleModelListSwagOutput(new(domain.Company))
	if err != nil {
		return err
	}

	g.GET("/:"+strings.ToLower("Company")+"_id", ctr.GetOne).
		AddParamPath("", "company_id", "Company ID").
		AddParamPath("", strings.ToLower("Company")+"_id", "Company ID").
		AddResponse(http.StatusOK, "Company in JSON:API format", &outModel, nil).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Gets one Company by id")

	g.GET("", ctr.GetList).
		AddParamPath("", "company_id", "Company ID").
		AddParamQuery("", "filter", "filter schema: filter[relation.relation_field][operator]=value", false).
		AddResponse(http.StatusOK, "Companys list in JSON:API format", &outModelList, nil).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Gets Companys list")

	g.POST("/", ctr.Create).
		AddParamBody(outModel, "body", "New user object", true).
		SetRequestContentType(jsonapi.MediaType).
		AddResponse(http.StatusOK, "User in JSON:API format", &outModel, nil).
		SetSummary("Creates user")

	g.PATCH("/:"+strings.ToLower("Company")+"_id", ctr.Update).
		AddParamBody(outModel, "body", "Updates Company object", true).
		SetRequestContentType(jsonapi.MediaType).
		AddResponse(http.StatusOK, "User in JSON:API format", &outModel, nil).
		SetSummary("Updates Company")

	g.DELETE("/:"+strings.ToLower("Company")+"_id", ctr.Delete).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Deletes user by id")

	return nil
}

// GetOne returns Company
func (ctr *CompanyController) GetOne(ctx echo.Context) error {

	model, err := ctr.manager.Company.GetOne(ctx.Request().Context(), ctx.Param(ctx.Param(strings.ToLower("Company")+"_id")))
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

// GetList returns Companys list
func (ctr *CompanyController) GetList(ctx echo.Context) error {

	filterStrings, err := filt.ParseQueryString(ctx.Request().URL.RawQuery, &domain.Company{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	models, err := ctr.manager.Company.GetList(ctx.Request().Context(), filterStrings)
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

// Create creates Company
func (ctr *CompanyController) Create(ctx echo.Context) error {
	domModel := new(domain.Company)
	err := httpUtils.UnmarshalFromRequest(domModel, ctx.Request().Body)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	err = validator.ValidateStruct(domModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusUnprocessableEntity, omniErr.ErrTitleValidation, err)
	}

	newModel, err := ctr.manager.Company.Create(ctx.Request().Context(), domModel)
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

// Update updates Company attributes
func (ctr *CompanyController) Update(ctx echo.Context) error {

	domModel := new(domain.Company)
	err := httpUtils.UnmarshalFromRequest(domModel, ctx.Request().Body)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	err = validator.ValidateStruct(domModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusUnprocessableEntity, omniErr.ErrTitleValidation, err)
	}

	newModel, err := ctr.manager.Company.Update(ctx.Request().Context(), domModel)
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

// Delete deletes Company
func (ctr *CompanyController) Delete(ctx echo.Context) error {
	err := ctr.manager.Company.Delete(ctx.Request().Context(), ctx.Param(strings.ToLower("Company")+"_id"))
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
