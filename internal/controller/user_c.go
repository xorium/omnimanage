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

type UserController struct {
	manager *service.Manager
}

func NewUserController(manager *service.Manager) *UserController {
	return &UserController{manager: manager}
}

// Init initializes routes and swag doc
func (ctr *UserController) Init(g echoswagger.ApiGroup) error {
	g.SetDescription("Operations about User")

	outModel, err := converter.GetExampleModelSwagOutput(new(domain.User))
	if err != nil {
		return err
	}

	outModelList, err := converter.GetExampleModelListSwagOutput(new(domain.User))
	if err != nil {
		return err
	}

	g.GET("/:"+strings.ToLower("User")+"_id", ctr.GetOne).
		AddParamPath("", "company_id", "Company ID").
		AddParamPath("", strings.ToLower("User")+"_id", "User ID").
		AddResponse(http.StatusOK, "User in JSON:API format", &outModel, nil).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Gets one User by id")

	g.GET("", ctr.GetList).
		AddParamPath("", "company_id", "Company ID").
		AddParamQuery("", "filter", "filter schema: filter[relation.relation_field][operator]=value", false).
		AddResponse(http.StatusOK, "Users list in JSON:API format", &outModelList, nil).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Gets Users list")

	g.POST("/", ctr.Create).
		AddParamBody(outModel, "body", "New user object", true).
		SetRequestContentType(jsonapi.MediaType).
		AddResponse(http.StatusOK, "User in JSON:API format", &outModel, nil).
		SetSummary("Creates user")

	g.PATCH("/:"+strings.ToLower("User")+"_id", ctr.Update).
		AddParamBody(outModel, "body", "Updates User object", true).
		SetRequestContentType(jsonapi.MediaType).
		AddResponse(http.StatusOK, "User in JSON:API format", &outModel, nil).
		SetSummary("Updates User")

	g.DELETE("/:"+strings.ToLower("User")+"_id", ctr.Delete).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Deletes user by id")

	// relations
	g.GET("/:"+strings.ToLower("User")+"_id/relationships/:relation_name", ctr.GetRelation).
		AddParamQuery("", "relation_name", "User relations. Possible values: company location roles subscriptions ", true).
		SetSummary("Gets User relation")

	g.POST("/:"+strings.ToLower("User")+"_id/relationships/:relation_name", ctr.ModifyRelation).
		AddParamQuery("", "relation_name", "User relations. Possible values: company location roles subscriptions ", true).
		SetSummary("Adds user relation")

	g.PATCH("/:"+strings.ToLower("User")+"_id/relationships/:relation_name", ctr.ModifyRelation).
		AddParamQuery("", "relation_name", "User relations. Possible values: company location roles subscriptions ", true).
		SetSummary("Replace user relation")

	g.DELETE("/:"+strings.ToLower("User")+"_id/relationships/:relation_name", ctr.ModifyRelation).
		AddParamQuery("", "relation_name", "User relations. Possible values: company location roles subscriptions ", true).
		SetSummary("Deletes user relation")

	return nil
}

// GetOne returns User
func (ctr *UserController) GetOne(ctx echo.Context) error {

	model, err := ctr.manager.User.GetOne(ctx.Request().Context(), ctx.Param(ctx.Param(strings.ToLower("User")+"_id")))
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

// GetList returns Users list
func (ctr *UserController) GetList(ctx echo.Context) error {

	filterStrings, err := filt.ParseQueryString(ctx.Request().URL.RawQuery, &domain.User{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	models, err := ctr.manager.User.GetList(ctx.Request().Context(), filterStrings)
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

// Create creates User
func (ctr *UserController) Create(ctx echo.Context) error {
	domModel := new(domain.User)
	err := httpUtils.UnmarshalFromRequest(domModel, ctx.Request().Body)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	err = validator.ValidateStruct(domModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusUnprocessableEntity, omniErr.ErrTitleValidation, err)
	}

	newModel, err := ctr.manager.User.Create(ctx.Request().Context(), domModel)
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

// Update updates User attributes
func (ctr *UserController) Update(ctx echo.Context) error {

	domModel := new(domain.User)
	err := httpUtils.UnmarshalFromRequest(domModel, ctx.Request().Body)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	err = validator.ValidateStruct(domModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusUnprocessableEntity, omniErr.ErrTitleValidation, err)
	}

	newModel, err := ctr.manager.User.Update(ctx.Request().Context(), domModel)
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

// Delete deletes User
func (ctr *UserController) Delete(ctx echo.Context) error {
	err := ctr.manager.User.Delete(ctx.Request().Context(), ctx.Param(strings.ToLower("User")+"_id"))
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
func (ctr *UserController) GetRelation(ctx echo.Context) error {
	webRelName := ctx.Param("relation_name")
	modelID := ctx.Param(strings.ToLower("User") + "_id")

	var (
		webData interface{}
		err     error
	)
	switch webRelName {

	case "company":
		webData, err = ctr.manager.User.GetCompany(ctx.Request().Context(), modelID)

	case "location":
		webData, err = ctr.manager.User.GetLocation(ctx.Request().Context(), modelID)

	case "roles":
		webData, err = ctr.manager.User.GetRoles(ctx.Request().Context(), modelID)

	case "subscriptions":
		webData, err = ctr.manager.User.GetSubscriptions(ctx.Request().Context(), modelID)

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
func (ctr *UserController) ModifyRelation(ctx echo.Context) error {
	webRelName := ctx.Param("relation_name")
	modelID := ctx.Param(strings.ToLower("User") + "_id")

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
			err = ctr.manager.User.ReplaceCompany(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodPost:
			err = ctr.manager.User.AppendCompany(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodDelete:
			err = ctr.manager.User.DeleteCompany(ctx.Request().Context(), modelID, domRelModel)
		}

	case "location":

		domRelModel := new(domain.Location)
		err := httpUtils.UnmarshalFromRequest(domRelModel, ctx.Request().Body)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleBadRequest, err)
		}

		switch ctx.Request().Method {
		case http.MethodPatch:
			err = ctr.manager.User.ReplaceLocation(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodPost:
			err = ctr.manager.User.AppendLocation(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodDelete:
			err = ctr.manager.User.DeleteLocation(ctx.Request().Context(), modelID, domRelModel)
		}

	case "roles":

		domRecordsIntf, err := httpUtils.UnmarshalManyFromRequest(new(domain.Role), ctx.Request().Body)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleBadRequest, err)
		}

		var domRelModel []*domain.Role
		err = converter.SliceI2SliceModel(domRecordsIntf, &domRelModel)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}

		switch ctx.Request().Method {
		case http.MethodPatch:
			err = ctr.manager.User.ReplaceRoles(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodPost:
			err = ctr.manager.User.AppendRoles(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodDelete:
			err = ctr.manager.User.DeleteRoles(ctx.Request().Context(), modelID, domRelModel)
		}

	case "subscriptions":

		domRecordsIntf, err := httpUtils.UnmarshalManyFromRequest(new(domain.Subscription), ctx.Request().Body)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleBadRequest, err)
		}

		var domRelModel []*domain.Subscription
		err = converter.SliceI2SliceModel(domRecordsIntf, &domRelModel)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}

		switch ctx.Request().Method {
		case http.MethodPatch:
			err = ctr.manager.User.ReplaceSubscriptions(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodPost:
			err = ctr.manager.User.AppendSubscriptions(ctx.Request().Context(), modelID, domRelModel)
		case http.MethodDelete:
			err = ctr.manager.User.DeleteSubscriptions(ctx.Request().Context(), modelID, domRelModel)
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
