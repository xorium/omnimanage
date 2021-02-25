package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"net/http"
	"omnimanage/internal/model"
	"omnimanage/internal/store"
	omniErr "omnimanage/pkg/error"
	filt "omnimanage/pkg/filters"
	"omnimanage/pkg/mapper"
	httpUtils "omnimanage/pkg/utils/http"
)

type RoleController struct {
	store *store.Store
	//logger
}

func NewRoleController(store *store.Store) *RoleController {
	return &RoleController{store: store}
}

// GetOne returns Role
func (ctr *RoleController) GetOne(ctx echo.Context) error {
	idSrc, err := mapper.GetSrcID(ctx.Param("id"), &model.Role{}, &omnimodels.Role{})
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

	srcRes, err := ctr.store.Roles.GetOne(ctx.Request().Context(), idSrc)
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

	webRes, err := srcRes.ToWeb()
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
			omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, webRes)
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
			omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
	}

	return nil
}

func (ctr *RoleController) GetList(ctx echo.Context) error {

	srcFilters, err := filt.ParseFiltersFromQueryToSrcModel(ctx.Request().URL.RawQuery, &omnimodels.Role{}, &model.Role{})
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

	srcRes, err := ctr.store.Roles.GetList(ctx.Request().Context(), srcFilters)
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

	webRes, err := srcRes.ToWeb()
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
			omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, webRes)
	if err != nil {
		return omniErr.NewHTTPError(ctx, http.StatusInternalServerError,
			omniErr.ErrCodeInternal, omniErr.ErrTitleInternal, err)
	}
	return nil
}
