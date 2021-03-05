package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"omnimanage/internal/store"
	omniErr "omnimanage/pkg/error"
	filt "omnimanage/pkg/filters"
	"omnimanage/pkg/mapper"
	"omnimanage/pkg/model/src"
	webmodels "omnimanage/pkg/model/web"
	httpUtils "omnimanage/pkg/utils/http"
)

type RoleController struct {
	store *store.Store
	//mapper *mapper.ModelMapper
	//logger
}

func NewRoleController(store *store.Store) *RoleController {
	return &RoleController{store: store}
}

// GetOne returns Role
func (ctr *RoleController) GetOne(ctx echo.Context) error {
	idSrc, err := mapper.Get().GetSrcID(ctx.Param("id"), &src.Role{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	srcRes, err := ctr.store.Roles.GetOne(ctx.Request().Context(), idSrc)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	webRes, err := srcRes.ToWeb()
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, webRes)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}

	return nil
}

func (ctr *RoleController) GetList(ctx echo.Context) error {

	srcFilters, err := filt.ParseFiltersFromQueryToSrcModel(ctx.Request().URL.RawQuery, &webmodels.Role{}, &src.Role{})
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrBadRequest:
			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	srcRes, err := ctr.store.Roles.GetList(ctx.Request().Context(), srcFilters)
	if err != nil {
		switch {
		case errors.Cause(err) == omniErr.ErrResourceNotFound:
			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
		default:
			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
		}
	}

	webRes, err := srcRes.ToWeb()
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}

	err = httpUtils.SetResponse(ctx, http.StatusOK, webRes)
	if err != nil {
		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
	}
	return nil
}
