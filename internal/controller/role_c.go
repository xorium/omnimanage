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

type RoleController struct {
	manager *service.Manager
}

func NewRoleController(manager *service.Manager) *RoleController {
	return &RoleController{manager: manager}
}

// Init initializes routes and swag doc
func (ctr *RoleController) Init(g echoswagger.ApiGroup) error {
	g.SetDescription("Operations about Role")

	outModel, err := converter.GetExampleModelSwagOutput(new(domain.Role))
	if err != nil {
		return err
	}

	outModelList, err := converter.GetExampleModelListSwagOutput(new(domain.Role))
	if err != nil {
		return err
	}

	g.GET("/:"+strings.ToLower("Role")+"_id", ctr.GetOne).
		AddParamPath("", "company_id", "Company ID").
		AddParamPath("", strings.ToLower("Role")+"_id", "Role ID").
		AddResponse(http.StatusOK, "Role in JSON:API format", &outModel, nil).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Gets one Role by id")

	g.GET("", ctr.GetList).
		AddParamPath("", "company_id", "Company ID").
		AddParamQuery("", "filter", "filter schema: filter[relation.relation_field][operator]=value", false).
		AddResponse(http.StatusOK, "Roles list in JSON:API format", &outModelList, nil).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Gets Roles list")

	g.POST("/", ctr.Create).
		AddParamBody(outModel, "body", "New user object", true).
		SetRequestContentType(jsonapi.MediaType).
		AddResponse(http.StatusOK, "User in JSON:API format", &outModel, nil).
		SetSummary("Creates user")

	g.PATCH("/:"+strings.ToLower("Role")+"_id", ctr.Update).
		AddParamBody(outModel, "body", "Updates Role object", true).
		SetRequestContentType(jsonapi.MediaType).
		AddResponse(http.StatusOK, "User in JSON:API format", &outModel, nil).
		SetSummary("Updates Role")

	g.DELETE("/:"+strings.ToLower("Role")+"_id", ctr.Delete).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Deletes user by id")

	// relations
	g.GET("/:"+strings.ToLower("Role")+"_id/relationships/:relation_name", ctr.GetRelation).
		AddParamQuery("", "relation_name", "Role relations. Possible values: company ", true).
		SetSummary("Gets Role relation")

	g.POST("/:"+strings.ToLower("Role")+"_id/relationships/:relation_name", ctr.ModifyRelation).
		AddParamQuery("", "relation_name", "Role relations. Possible values: company ", true).
		SetSummary("Adds user relation")

	g.PATCH("/:"+strings.ToLower("Role")+"_id/relationships/:relation_name", ctr.ModifyRelation).
		AddParamQuery("", "relation_name", "Role relations. Possible values: company ", true).
		SetSummary("Replace user relation")

	g.DELETE("/:"+strings.ToLower("Role")+"_id/relationships/:relation_name", ctr.ModifyRelation).
		AddParamQuery("", "relation_name", "Role relations. Possible values: company ", true).
		SetSummary("Deletes user relation")

	return nil
}

// GetOne returns Role
func (ctr *RoleController) GetOne(ctx echo.Context) error {

	model, err := ctr.manager.Role.GetOne(ctx.Request().Context(), ctx.Param(ctx.Param(strings.ToLower("Role")+"_id")))
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

// GetList returns Roles list
func (ctr *RoleController) GetList(ctx echo.Context) error {

	filterStrings, err := filt.ParseQueryString(ctx.Request().URL.RawQuery, &domain.Role{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	models, err := ctr.manager.Role.GetList(ctx.Request().Context(), filterStrings)
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

// Create creates Role
func (ctr *RoleController) Create(ctx echo.Context) error {
	domModel := new(domain.Role)
	err := httpUtils.UnmarshalFromRequest(domModel, ctx.Request().Body)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	err = validator.ValidateStruct(domModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusUnprocessableEntity, omniErr.ErrTitleValidation, err)
	}

	newModel, err := ctr.manager.Role.Create(ctx.Request().Context(), domModel)
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

// Update updates Role attributes
func (ctr *RoleController) Update(ctx echo.Context) error {

	domModel := new(domain.Role)
	err := httpUtils.UnmarshalFromRequest(domModel, ctx.Request().Body)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	err = validator.ValidateStruct(domModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusUnprocessableEntity, omniErr.ErrTitleValidation, err)
	}

	newModel, err := ctr.manager.Role.Update(ctx.Request().Context(), domModel)
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

// Delete deletes Role
func (ctr *RoleController) Delete(ctx echo.Context) error {
	err := ctr.manager.Role.Delete(ctx.Request().Context(), ctx.Param(strings.ToLower("Role")+"_id"))
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

// GetRelation returns relation data
func (ctr *RoleController) GetRelation(ctx echo.Context) error {
	webRelName := ctx.Param("relation_name")
	modelID := ctx.Param(strings.ToLower("Role") + "_id")

	var (
		webData interface{}
		err     error
	)
	switch webRelName {

	case "company":
		webData, err = ctr.manager.Role.GetCompany(ctx.Request().Context(), modelID)

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
func (ctr *RoleController) ModifyRelation(ctx echo.Context) error {
	webRelName := ctx.Param("relation_name")
	modelID := ctx.Param(strings.ToLower("Role") + "_id")

	var err error

	switch webRelName {

	case "company":

		domRelModel := new(domain.Company)
		err := httpUtils.UnmarshalFromRequest(domRelModel, ctx.Request().Body)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleBadRequest, err)
		}

		switch ctx.Request().Method {
		case http.MethodPatch:
			err = ctr.manager.Role.ReplaceCompany(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodPost:
			err = ctr.manager.Role.AppendCompany(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodDelete:
			err = ctr.manager.Role.DeleteCompany(ctx.Request().Context(), modelID, domRelModel)
		}

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
