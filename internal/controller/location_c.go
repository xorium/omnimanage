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

type LocationController struct {
	manager *service.Manager
}

func NewLocationController(manager *service.Manager) *LocationController {
	return &LocationController{manager: manager}
}

// Init initializes routes and swag doc
func (ctr *LocationController) Init(g echoswagger.ApiGroup) error {
	g.SetDescription("Operations about Location")

	outModel, err := converter.GetExampleModelSwagOutput(new(domain.Location))
	if err != nil {
		return err
	}

	outModelList, err := converter.GetExampleModelListSwagOutput(new(domain.Location))
	if err != nil {
		return err
	}

	g.GET("/:"+strings.ToLower("Location")+"_id", ctr.GetOne).
		AddParamPath("", "company_id", "Company ID").
		AddParamPath("", strings.ToLower("Location")+"_id", "Location ID").
		AddResponse(http.StatusOK, "Location in JSON:API format", &outModel, nil).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Gets one Location by id")

	g.GET("", ctr.GetList).
		AddParamPath("", "company_id", "Company ID").
		AddParamQuery("", "filter", "filter schema: filter[relation.relation_field][operator]=value", false).
		AddResponse(http.StatusOK, "Locations list in JSON:API format", &outModelList, nil).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Gets Locations list")

	g.POST("/", ctr.Create).
		AddParamBody(outModel, "body", "New user object", true).
		SetRequestContentType(jsonapi.MediaType).
		AddResponse(http.StatusOK, "User in JSON:API format", &outModel, nil).
		SetSummary("Creates user")

	g.PATCH("/:"+strings.ToLower("Location")+"_id", ctr.Update).
		AddParamBody(outModel, "body", "Updates Location object", true).
		SetRequestContentType(jsonapi.MediaType).
		AddResponse(http.StatusOK, "User in JSON:API format", &outModel, nil).
		SetSummary("Updates Location")

	g.DELETE("/:"+strings.ToLower("Location")+"_id", ctr.Delete).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Deletes user by id")

	// relations
	g.GET("/:"+strings.ToLower("Location")+"_id/relationships/:relation_name", ctr.GetRelation).
		AddParamQuery("", "relation_name", "Location relations. Possible values: company children users ", true).
		SetSummary("Gets Location relation")

	g.POST("/:"+strings.ToLower("Location")+"_id/relationships/:relation_name", ctr.ModifyRelation).
		AddParamQuery("", "relation_name", "Location relations. Possible values: company children users ", true).
		SetSummary("Adds user relation")

	g.PATCH("/:"+strings.ToLower("Location")+"_id/relationships/:relation_name", ctr.ModifyRelation).
		AddParamQuery("", "relation_name", "Location relations. Possible values: company children users ", true).
		SetSummary("Replace user relation")

	g.DELETE("/:"+strings.ToLower("Location")+"_id/relationships/:relation_name", ctr.ModifyRelation).
		AddParamQuery("", "relation_name", "Location relations. Possible values: company children users ", true).
		SetSummary("Deletes user relation")

	return nil
}

// GetOne returns Location
func (ctr *LocationController) GetOne(ctx echo.Context) error {

	model, err := ctr.manager.Location.GetOne(ctx.Request().Context(), ctx.Param(ctx.Param(strings.ToLower("Location")+"_id")))
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

// GetList returns Locations list
func (ctr *LocationController) GetList(ctx echo.Context) error {

	filterStrings, err := filt.ParseQueryString(ctx.Request().URL.RawQuery, &domain.Location{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	models, err := ctr.manager.Location.GetList(ctx.Request().Context(), filterStrings)
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

// Create creates Location
func (ctr *LocationController) Create(ctx echo.Context) error {
	domModel := new(domain.Location)
	err := httpUtils.UnmarshalFromRequest(domModel, ctx.Request().Body)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	err = validator.ValidateStruct(domModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusUnprocessableEntity, omniErr.ErrTitleValidation, err)
	}

	newModel, err := ctr.manager.Location.Create(ctx.Request().Context(), domModel)
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

// Update updates Location attributes
func (ctr *LocationController) Update(ctx echo.Context) error {

	domModel := new(domain.Location)
	err := httpUtils.UnmarshalFromRequest(domModel, ctx.Request().Body)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	err = validator.ValidateStruct(domModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusUnprocessableEntity, omniErr.ErrTitleValidation, err)
	}

	newModel, err := ctr.manager.Location.Update(ctx.Request().Context(), domModel)
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

// Delete deletes Location
func (ctr *LocationController) Delete(ctx echo.Context) error {
	err := ctr.manager.Location.Delete(ctx.Request().Context(), ctx.Param(strings.ToLower("Location")+"_id"))
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
func (ctr *LocationController) GetRelation(ctx echo.Context) error {
	webRelName := ctx.Param("relation_name")
	modelID := ctx.Param(strings.ToLower("Location") + "_id")

	var (
		webData interface{}
		err     error
	)
	switch webRelName {

	case "company":
		webData, err = ctr.manager.Location.GetCompany(ctx.Request().Context(), modelID)

	case "children":
		webData, err = ctr.manager.Location.GetChildren(ctx.Request().Context(), modelID)

	case "users":
		webData, err = ctr.manager.Location.GetUsers(ctx.Request().Context(), modelID)

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
func (ctr *LocationController) ModifyRelation(ctx echo.Context) error {
	webRelName := ctx.Param("relation_name")
	modelID := ctx.Param(strings.ToLower("Location") + "_id")

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
			err = ctr.manager.Location.ReplaceCompany(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodPost:
			err = ctr.manager.Location.AppendCompany(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodDelete:
			err = ctr.manager.Location.DeleteCompany(ctx.Request().Context(), modelID, domRelModel)
		}

	case "children":

		domRecordsIntf, err := httpUtils.UnmarshalManyFromRequest(new(domain.Location), ctx.Request().Body)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleBadRequest, err)
		}

		var domRelModel []*domain.Location
		err = converter.SliceI2SliceModel(domRecordsIntf, &domRelModel)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}

		switch ctx.Request().Method {
		case http.MethodPatch:
			err = ctr.manager.Location.ReplaceChildren(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodPost:
			err = ctr.manager.Location.AppendChildren(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodDelete:
			err = ctr.manager.Location.DeleteChildren(ctx.Request().Context(), modelID, domRelModel)
		}

	case "users":

		domRecordsIntf, err := httpUtils.UnmarshalManyFromRequest(new(domain.User), ctx.Request().Body)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleBadRequest, err)
		}

		var domRelModel []*domain.User
		err = converter.SliceI2SliceModel(domRecordsIntf, &domRelModel)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}

		switch ctx.Request().Method {
		case http.MethodPatch:
			err = ctr.manager.Location.ReplaceUsers(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodPost:
			err = ctr.manager.Location.AppendUsers(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodDelete:
			err = ctr.manager.Location.DeleteUsers(ctx.Request().Context(), modelID, domRelModel)
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
