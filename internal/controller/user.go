package controller

import (
	"bytes"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"io"
	"net/http"
	"omnimanage/internal/model"
	"omnimanage/internal/store"
	filt "omnimanage/pkg/filters"
	"strconv"
)

type UserController struct {
	store *store.Store
	//logger
}

func NewUserController(store *store.Store) *UserController {
	return &UserController{store: store}
}

// Get returns User
func (ctr *UserController) GetOne(ctx echo.Context) error {
	ctx.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse ID"))
	}

	user, err := ctr.store.Users.GetOne(ctx.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	webUser, err := user.ToWeb()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = MarshalToResponse(webUser, ctx.Response())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	ctx.Response().WriteHeader(http.StatusOK)

	//userID, err := uuid.Parse(ctx.Param("id"))
	//if err != nil {
	//	return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse user UUID"))
	//}
	//user, err := ctr.services.User.GetUser(ctx.Request().Context(), userID)
	//if err != nil {
	//	switch {
	//	case errors.Cause(err) == types.ErrNotFound:
	//		return echo.NewHTTPError(http.StatusNotFound, err)
	//	case errors.Cause(err) == types.ErrBadRequest:
	//		return echo.NewHTTPError(http.StatusBadRequest, err)
	//	default:
	//		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not get user"))
	//	}
	//}
	//return ctx.JSON(http.StatusOK, user)
	return nil
}

func MarshalToResponse(model interface{}, w io.Writer) error {
	var b bytes.Buffer
	err := jsonapi.MarshalPayload(&b, model)
	if err != nil {
		return err
	}

	_, err = b.WriteTo(w)
	if err != nil {
		return err
	}
	return nil
}

func (ctr *UserController) GetList(ctx echo.Context) error {
	ctx.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	srcFilters, err := filt.ParseFiltersFromQueryToSrcModel(ctx.Request().URL.RawQuery, &omnimodels.User{}, &model.User{})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	users, err := ctr.store.Users.GetList(ctx.Request().Context(), srcFilters)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	webUsers, err := model.UsersToWeb(users)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = MarshalToResponse(webUsers, ctx.Response())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}

func (ctr *UserController) GetRelation(ctx echo.Context) error {
	ctx.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse ID"))
	}

	relName := ctx.Param("rel")

	user, err := ctr.store.Users.GetOne(ctx.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	switch relName {
	case "location":
		loc, err := ctr.store.Locations.GetOne(ctx.Request().Context(), user.LocationID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		err = MarshalToResponse(loc.ToWeb(), ctx.Response())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

	default:
		return echo.NewHTTPError(http.StatusBadRequest, "wrong relation name: %v", relName)
	}

	return nil
}
