package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"gitlab.omnicube.ru/libs/omnilib/utils/converter"
	"net/http"
	"omnimanage/internal/model"
	"omnimanage/internal/store"
	omniErr "omnimanage/pkg/error"
	filt "omnimanage/pkg/filters"
	"omnimanage/pkg/mapper"
	httpUtils "omnimanage/pkg/utils/http"
)

type UserController struct {
	store *store.Store
	//logger
}

func NewUserController(store *store.Store) *UserController {
	return &UserController{store: store}
}

// GetOne returns User
func (ctr *UserController) GetOne(ctx echo.Context) error {
	idSrc, err := mapper.GetSrcID(ctx.Param("id"), &model.User{}, &omnimodels.User{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(ctx, http.StatusBadRequest,
				omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}
	}

	user, err := ctr.store.Users.GetOne(ctx.Request().Context(), idSrc)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(ctx, http.StatusNotFound,
				omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}
	}

	webUser, err := user.ToWeb()
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
			omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, webUser)
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
			omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
	}

	return nil
}

func (ctr *UserController) GetList(ctx echo.Context) error {

	srcFilters, err := filt.ParseFiltersFromQueryToSrcModel(ctx.Request().URL.RawQuery, &omnimodels.User{}, &model.User{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(ctx, http.StatusBadRequest,
				omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}
	}

	users, err := ctr.store.Users.GetList(ctx.Request().Context(), srcFilters)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(ctx, http.StatusNotFound,
				omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}
	}

	webUsers, err := users.ToWeb()
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
			omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, webUsers)
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
			omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
	}
	return nil
}

func (ctr *UserController) GetRelation(ctx echo.Context) error {

	idSrc, err := mapper.GetSrcID(ctx.Param("id"), &model.User{}, &omnimodels.User{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(ctx, http.StatusBadRequest,
				omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}
	}

	user, err := ctr.store.Users.GetOne(ctx.Request().Context(), idSrc)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(ctx, http.StatusNotFound,
				omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}
	}

	relName := ctx.Param("rel")
	switch relName {
	case "location":
		//loc, err := ctr.store.Locations.GetOne(ctx.Request().Context(), user.LocationID)
		//if err != nil {
		//	return echo.NewHTTPError(http.StatusBadRequest, err)
		//}
		//web, err := loc.ToWeb()
		//if err != nil {
		//	return echo.NewHTTPError(http.StatusBadRequest, err)
		//}
		//err = httpUtils.SetResponse(ctx, http.StatusOK, web)
		//if err != nil {
		//	return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
		//		omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		//}
	case "roles":
		srcFilters, err := filt.GetSrcFiltersFromRelationID(user.Roles)
		if err != nil {
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}

		srcList, err := ctr.store.Roles.GetList(ctx.Request().Context(), srcFilters)
		if err != nil {
			switch {
			case errors.Cause(err) == omniErr.ErrResourceNotFound:
				return omniErr.NewHTTPError(ctx, http.StatusNotFound,
					omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
			default:
				return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
					omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
			}
		}

		webList, err := srcList.ToWeb()
		if err != nil {
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}

		err = httpUtils.SetResponse(ctx, http.StatusOK, webList)
		if err != nil {
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}
	default:
		return omniErr.NewHTTPError(ctx, http.StatusForbidden,
			omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound,
			fmt.Errorf("%w wrong relation name '%v'", omniErr.ErrResourceNotFound, relName))
	}

	return nil
}

func (ctr *UserController) ModifyRelation(ctx echo.Context) error {

	idSrc, err := mapper.GetSrcID(ctx.Param("id"), &model.User{}, &omnimodels.User{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(ctx, http.StatusBadRequest,
				omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}
	}

	webRelName := ctx.Param("rel")
	switch webRelName {
	case "location":

	case "roles":

		rolesIntf, err := httpUtils.UnmarshalManyFromRequest(new(omnimodels.Role), ctx.Request().Body)
		if err != nil {
			return omniErr.NewHTTPError(ctx, http.StatusBadRequest,
				omniErr.ErrCodeResource, omniErr.ErrTitleBadRequest, err)
		}

		var webRoles []*omnimodels.Role
		err = converter.SliceI2SliceModel(rolesIntf, &webRoles)
		if err != nil {
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}

		srcRolesNew, err := model.Roles{}.ScanFromWeb(webRoles)
		if err != nil {
			return omniErr.NewHTTPError(ctx, http.StatusBadRequest,
				omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
		}

		srcRelName := "Roles"
		switch ctx.Request().Method {
		case http.MethodPatch:
			err = ctr.store.Users.ReplaceRelation(ctx.Request().Context(), idSrc, srcRelName, srcRolesNew)
		case http.MethodPost:
			err = ctr.store.Users.AppendRelation(ctx.Request().Context(), idSrc, srcRelName, srcRolesNew)
		case http.MethodDelete:
			err = ctr.store.Users.DeleteRelation(ctx.Request().Context(), idSrc, srcRelName, srcRolesNew)
		}
		if err != nil {
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}

	default:
		return omniErr.NewHTTPError(ctx, http.StatusForbidden,
			omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound,
			fmt.Errorf("%w wrong relation name '%v'", omniErr.ErrResourceNotFound, webRelName))
	}

	return nil

}

func (ctr *UserController) Create(ctx echo.Context) error {
	webModel := new(omnimodels.User)
	err := httpUtils.UnmarshalFromRequest(webModel, ctx.Request().Body)
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusBadRequest,
			omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
	}

	srcUser, err := new(model.User).ScanFromWeb(webModel)
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusBadRequest,
			omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
	}

	user, err := ctr.store.Users.Create(ctx.Request().Context(), srcUser)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceExists:
			return omniErr.NewHTTPError(ctx, http.StatusConflict,
				omniErr.ErrCodeResource, omniErr.ErrTitleResourceExists, err)
		default:
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}
	}

	webUser, err := user.ToWeb()
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
			omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, webUser)
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
			omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
	}

	return nil
}

func (ctr *UserController) Update(ctx echo.Context) error {

	webModel := new(omnimodels.User)
	err := httpUtils.UnmarshalFromRequest(webModel, ctx.Request().Body)
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusBadRequest,
			omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
	}

	srcUser, err := new(model.User).ScanFromWeb(webModel)
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusBadRequest,
			omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
	}

	user, err := ctr.store.Users.Update(ctx.Request().Context(), srcUser)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(ctx, http.StatusBadRequest,
				omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}
	}

	webUser, err := user.ToWeb()
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
			omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, webUser)
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
			omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
	}

	return nil

}

func (ctr *UserController) Delete(ctx echo.Context) error {
	idSrc, err := mapper.GetSrcID(ctx.Param("id"), &model.User{}, &omnimodels.User{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(ctx, http.StatusBadRequest,
				omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}
	}

	err = ctr.store.Users.Delete(ctx.Request().Context(), idSrc)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(ctx, http.StatusNotFound,
				omniErr.ErrCodeResource, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
				omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
		}
	}

	ctx.NoContent(http.StatusNoContent)
	return nil
}
