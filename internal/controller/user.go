package controller

import (
	"fmt"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/pkg/errors"
	"net/http"
	"omnimanage/internal/store"
	"omnimanage/internal/validator"
	omniErr "omnimanage/pkg/error"
	filt "omnimanage/pkg/filters"
	"omnimanage/pkg/mapper"
	"omnimanage/pkg/model/src"
	webmodels "omnimanage/pkg/model/web"
	"omnimanage/pkg/utils/converter"
	httpUtils "omnimanage/pkg/utils/http"
)

type UserController struct {
	store *store.Store
	//mapper *mapper.ModelMapper
	//logger
}

func NewUserController(store *store.Store) *UserController {
	return &UserController{store: store}
}

func (ctr *UserController) Init(g echoswagger.ApiGroup) error {
	g.SetDescription("Operations about user")

	outModelRec, err := converter.ModelSwagOutput(new(webmodels.User))
	if err != nil {
		return err
	}

	g.GET("/:id", ctr.GetOne).
		AddResponse(http.StatusOK, "successful operation", &outModelRec, nil).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Gets one user by id")

	g.GET("", ctr.GetList).
		SetSummary("Gets users list")

	g.POST("/", ctr.Create).
		AddParamBody(outModelRec, "body", "Created user object", true).
		SetRequestContentType(jsonapi.MediaType).
		SetSummary("Creates user")

	g.PATCH("/:id", ctr.Update).
		AddParamBody(outModelRec, "body", "Updates user object", true).
		SetRequestContentType(jsonapi.MediaType).
		SetSummary("Updates user")

	g.DELETE("/:id", ctr.Delete).
		SetResponseContentType(jsonapi.MediaType).
		SetSummary("Deletes user by id")

	// relations
	g.GET("/:id/relationships/:rel", ctr.GetRelation)

	g.POST("/:id/relationships/:rel", ctr.ModifyRelation)
	g.PATCH("/:id/relationships/:rel", ctr.ModifyRelation)
	g.DELETE("/:id/relationships/:rel", ctr.ModifyRelation)

	return nil
}

// GetOne returns User
func (ctr *UserController) GetOne(ctx echo.Context) error {
	idSrc, err := mapper.Get().GetSrcID(ctx.Param("id"), &src.User{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	user, err := ctr.store.Users.GetOne(ctx.Request().Context(), idSrc)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	webUser, err := user.ToWeb()
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, webUser)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}

	return nil
}

// GetList returns users list
func (ctr *UserController) GetList(ctx echo.Context) error {

	srcFilters, err := filt.ParseFiltersFromQueryToSrcModel(ctx.Request().URL.RawQuery, &webmodels.User{}, &src.User{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	users, err := ctr.store.Users.GetList(ctx.Request().Context(), srcFilters)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	webUsers, err := users.ToWeb()
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, webUsers)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}
	return nil
}

// GetRelation returns relation data
func (ctr *UserController) GetRelation(ctx echo.Context) error {

	idSrc, err := mapper.Get().GetSrcID(ctx.Param("id"), &src.User{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	user, err := ctr.store.Users.GetOne(ctx.Request().Context(), idSrc)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	relName := ctx.Param("rel")
	switch relName {
	case "location":
		loc, err := ctr.store.Locations.GetOne(ctx.Request().Context(), user.LocationID)
		if err != nil {
			switch {
			case errors.Cause(err) == omniErr.ErrResourceNotFound:
				return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
			default:
				return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
			}
		}

		web, err := loc.ToWeb()
		if err != nil {
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
		err = httpUtils.SetResponse(ctx, http.StatusOK, web)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}

	case "roles":
		srcFilters, err := filt.GetSrcFiltersFromRelationID(user.Roles)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}

		srcList, err := ctr.store.Roles.GetList(ctx.Request().Context(), srcFilters)
		if err != nil {
			switch {
			case errors.Cause(err) == omniErr.ErrResourceNotFound:
				return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
			default:
				return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
			}
		}

		webList, err := srcList.ToWeb()
		if err != nil {
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}

		err = httpUtils.SetResponse(ctx, http.StatusOK, webList)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	default:
		return omniErr.NewHTTPError(http.StatusForbidden, omniErr.ErrTitleResourceNotFound,
			fmt.Errorf("%w wrong relation name '%v'", omniErr.ErrResourceNotFound, relName))
	}

	return nil
}

// ModifyRelation - create, delete, replace relations
func (ctr *UserController) ModifyRelation(ctx echo.Context) error {

	idSrc, err := mapper.Get().GetSrcID(ctx.Param("id"), &src.User{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	webRelName := ctx.Param("rel")
	switch webRelName {
	case "location":
		webModel := new(webmodels.Location)
		err := httpUtils.UnmarshalFromRequest(webModel, ctx.Request().Body)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleBadRequest, err)
		}

		srcModelsNew, err := new(src.Location).ScanFromWeb(webModel)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		}

		srcRelName := "Location"
		switch ctx.Request().Method {
		case http.MethodPatch:
			err = ctr.store.Users.ReplaceRelation(ctx.Request().Context(), idSrc, srcRelName, srcModelsNew)
		case http.MethodPost:
			err = ctr.store.Users.AppendRelation(ctx.Request().Context(), idSrc, srcRelName, srcModelsNew)
		case http.MethodDelete:
			err = ctr.store.Users.DeleteRelation(ctx.Request().Context(), idSrc, srcRelName, srcModelsNew)
		}
		if err != nil {
			switch {
			case errors.Cause(err) == omniErr.ErrResourceNotFound:
				return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
			default:
				return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
			}
		}

	case "roles":

		webRecordsIntf, err := httpUtils.UnmarshalManyFromRequest(new(webmodels.Role), ctx.Request().Body)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleBadRequest, err)
		}

		var webModels []*webmodels.Role
		err = converter.SliceI2SliceModel(webRecordsIntf, &webModels)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}

		srcModelsNew, err := src.Roles.ScanFromWeb(nil, webModels)
		if err != nil {
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		}

		srcRelName := "Roles"
		switch ctx.Request().Method {
		case http.MethodPatch:
			err = ctr.store.Users.ReplaceRelation(ctx.Request().Context(), idSrc, srcRelName, srcModelsNew)
		case http.MethodPost:
			err = ctr.store.Users.AppendRelation(ctx.Request().Context(), idSrc, srcRelName, srcModelsNew)
		case http.MethodDelete:
			err = ctr.store.Users.DeleteRelation(ctx.Request().Context(), idSrc, srcRelName, srcModelsNew)
		}
		if err != nil {
			switch {
			case errors.Cause(err) == omniErr.ErrResourceNotFound:
				return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
			default:
				return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
			}
		}

	default:
		return omniErr.NewHTTPError(http.StatusForbidden, omniErr.ErrTitleResourceNotFound,
			fmt.Errorf("%w wrong relation name '%v'", omniErr.ErrResourceNotFound, webRelName))
	}

	return nil

}

// Create creates user
func (ctr *UserController) Create(ctx echo.Context) error {
	webModel := new(webmodels.User)
	err := httpUtils.UnmarshalFromRequest(webModel, ctx.Request().Body)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	err = validator.ValidateStruct(webModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusUnprocessableEntity, omniErr.ErrTitleValidation, err)
	}

	srcUser, err := new(src.User).ScanFromWeb(webModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	user, err := ctr.store.Users.Create(ctx.Request().Context(), srcUser)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceExists:
			return omniErr.NewHTTPError(http.StatusConflict, omniErr.ErrTitleResourceExists, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	webUser, err := user.ToWeb()
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, webUser)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}

	return nil
}

// Update updates user attributes
func (ctr *UserController) Update(ctx echo.Context) error {

	webModel := new(webmodels.User)
	err := httpUtils.UnmarshalFromRequest(webModel, ctx.Request().Body)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	srcUser, err := new(src.User).ScanFromWeb(webModel)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
	}

	user, err := ctr.store.Users.Update(ctx.Request().Context(), srcUser)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	webUser, err := user.ToWeb()
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, webUser)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}

	return nil

}

// Delete deletes user
func (ctr *UserController) Delete(ctx echo.Context) error {
	idSrc, err := mapper.Get().GetSrcID(ctx.Param("id"), &src.User{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	err = ctr.store.Users.Delete(ctx.Request().Context(), idSrc)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	ctx.NoContent(http.StatusNoContent)
	return nil
}
