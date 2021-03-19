package controller

//import (
//	"github.com/labstack/echo/v4"
//	"github.com/pangpanglabs/echoswagger/v2"
//	"omnimanage/internal/service"
//)
//
//type RoleController struct {
//	//store *store.Store
//	//logger
//	manager *service.Manager
//}
//
//func NewRoleController(manager *service.Manager) *RoleController {
//	return &RoleController{manager: manager}
//}
//
//// Init initializes routes and swag doc
//func (ctr *RoleController) Init(g echoswagger.ApiGroup) error {
//
//	g.GET("", ctr.GetList)
//	g.GET("/:id", ctr.GetOne)
//
//	return nil
//}
//
//// GetOne returns Role
//func (ctr *RoleController) GetOne(ctx echo.Context) error {
//	//idSrc, err := mapper.Get().GetSrcID(ctx.Param("role_id"), &src.Role{})
//	//if err != nil {
//	//	switch {
//	//	case errors.Cause(err) == omniErr.ErrBadRequest:
//	//		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//	//	default:
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//}
//	//
//	//srcRes, err := ctr.manager.Role.GetOne(ctx.Request().Context(), idSrc)
//	//if err != nil {
//	//	switch {
//	//	case errors.Cause(err) == omniErr.ErrResourceNotFound:
//	//		return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
//	//	default:
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//}
//	//
//	//webRes, err := srcRes.ToWeb()
//	//if err != nil {
//	//	return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//}
//	//
//	//err = httpUtils.SetResponse(ctx, http.StatusOK, webRes)
//	//if err != nil {
//	//	return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//}
//
//	return nil
//}
//
//func (ctr *RoleController) GetList(ctx echo.Context) error {
//
//	//srcFilters, err := filt.ParseFiltersFromQueryToSrcModel(ctx.Request().URL.RawQuery, &webmodels.Role{}, &src.Role{})
//	//if err != nil {
//	//	switch {
//	//	case errors.Cause(err) == omniErr.ErrBadRequest:
//	//		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//	//	default:
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//}
//	//
//	//srcRes, err := ctr.store.Roles.GetList(ctx.Request().Context(), srcFilters)
//	//if err != nil {
//	//	switch {
//	//	case errors.Cause(err) == omniErr.ErrResourceNotFound:
//	//		return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
//	//	default:
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//}
//	//
//	//webRes, err := srcRes.ToWeb()
//	//if err != nil {
//	//	return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//}
//	//
//	//err = httpUtils.SetResponse(ctx, http.StatusOK, webRes)
//	//if err != nil {
//	//	return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//}
//	return nil
//}
