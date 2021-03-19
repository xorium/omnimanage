package controller

//
//import (
//	"fmt"
//	"github.com/google/jsonapi"
//	"github.com/labstack/echo/v4"
//	"github.com/pangpanglabs/echoswagger/v2"
//	"github.com/pkg/errors"
//	"net/http"
//	"omnimanage/internal/service"
//	"omnimanage/internal/validator"
//	omniErr "omnimanage/pkg/error"
//	filt "omnimanage/pkg/filters"
//	"omnimanage/pkg/model/domain"
//	"omnimanage/pkg/utils/converter"
//	httpUtils "omnimanage/pkg/utils/http"
//)
//
//type UserController struct {
//	//store *store.Store
//	manager *service.Manager
//	//logger
//}
//
//func NewUserController(manager *service.Manager) *UserController {
//	return &UserController{manager: manager}
//}
//
//// Init initializes routes and swag doc
//func (ctr *UserController) Init(g echoswagger.ApiGroup) error {
//	g.SetDescription("Operations about user")
//
//	outModel, err := converter.GetExampleModelSwagOutput(new(domain.User))
//	if err != nil {
//		return err
//	}
//
//	outModelList, err := converter.GetExampleModelListSwagOutput(new(domain.User))
//	if err != nil {
//		return err
//	}
//
//	g.GET("/:user_id", ctr.GetOne).
//		AddParamPath("", "company_id", "Company ID").
//		AddParamPath("", "user_id", "User ID").
//		AddResponse(http.StatusOK, "User in JSON:API format", &outModel, nil).
//		SetResponseContentType(jsonapi.MediaType).
//		SetSummary("Gets one user by id")
//
//	g.GET("", ctr.GetList).
//		AddParamPath("", "company_id", "Company ID").
//		AddParamQuery("", "filter", "filter schema: filter[relation.relation_field][operator]=value", false).
//		AddResponse(http.StatusOK, "Users list in JSON:API format", &outModelList, nil).
//		SetResponseContentType(jsonapi.MediaType).
//		SetSummary("Gets users list")
//
//	g.POST("/", ctr.Create).
//		AddParamBody(outModel, "body", "New user object", true).
//		SetRequestContentType(jsonapi.MediaType).
//		AddResponse(http.StatusOK, "User in JSON:API format", &outModel, nil).
//		SetSummary("Creates user")
//
//	g.PATCH("/:user_id", ctr.Update).
//		AddParamBody(outModel, "body", "Updates user object", true).
//		SetRequestContentType(jsonapi.MediaType).
//		AddResponse(http.StatusOK, "User in JSON:API format", &outModel, nil).
//		SetSummary("Updates user")
//
//	g.DELETE("/:user_id", ctr.Delete).
//		SetResponseContentType(jsonapi.MediaType).
//		SetSummary("Deletes user by id")
//
//	// relations
//	g.GET("/:user_id/relationships/:relation_name", ctr.GetRelation).
//		AddParamQuery("", "relation_name", "User relations. Possible values: company, location, roles, subscriptions", true).
//		SetSummary("Gets user relation (company, location, roles, subscriptions)")
//
//	g.POST("/:user_id/relationships/:relation_name", ctr.ModifyRelation).
//		AddParamQuery("", "relation_name", "User relations. Possible values: company, location, roles, subscriptions", true).
//		SetSummary("Adds user relation (company, location, roles, subscriptions)")
//
//	g.PATCH("/:user_id/relationships/:relation_name", ctr.ModifyRelation).
//		AddParamQuery("", "relation_name", "User relations. Possible values: company, location, roles, subscriptions", true).
//		SetSummary("Updates user relation (company, location, roles, subscriptions)")
//
//	g.DELETE("/:user_id/relationships/:relation_name", ctr.ModifyRelation).
//		AddParamQuery("", "relation_name", "User relations. Possible values: company, location, roles, subscriptions", true).
//		SetSummary("Deletes user relation (company, location, roles, subscriptions)")
//
//	return nil
//}
//
//// GetOne returns User
//func (ctr *UserController) GetOne(ctx echo.Context) error {
//
//	user, err := ctr.manager.User.GetOne(ctx.Request().Context(), ctx.Param("user_id"))
//	if err != nil {
//		switch {
//		case errors.Cause(err) == omniErr.ErrResourceNotFound:
//			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
//		default:
//			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//		}
//	}
//
//	err = httpUtils.SetResponse(ctx, http.StatusOK, user)
//	if err != nil {
//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	}
//
//	return nil
//}
//
//// GetList returns users list
//func (ctr *UserController) GetList(ctx echo.Context) error {
//
//	//srcFilters, err := filt.ParseFiltersFromQueryToSrcModel(ctx.Request().URL.RawQuery, &webmodels.User{}, &src.User{})
//	filterStrings, err := filt.ParseQueryString(ctx.Request().URL.RawQuery, &domain.User{})
//	if err != nil {
//		switch {
//		case errors.Cause(err) == omniErr.ErrBadRequest:
//			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//		default:
//			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//		}
//	}
//
//	users, err := ctr.manager.User.GetList(ctx.Request().Context(), filterStrings)
//	if err != nil {
//		switch {
//		case errors.Cause(err) == omniErr.ErrResourceNotFound:
//			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
//		default:
//			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//		}
//	}
//
//	//webUsers, err := users.ToWeb()
//	//if err != nil {
//	//	return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//}
//
//	err = httpUtils.SetResponse(ctx, http.StatusOK, users)
//	if err != nil {
//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	}
//	return nil
//}
//
//// GetRelation returns relation data
//func (ctr *UserController) GetRelation(ctx echo.Context) error {
//
//	relModel, err := ctr.manager.User.GetRelation(ctx.Request().Context(), ctx.Param("user_id"), ctx.Param("relation_name"))
//	if err != nil {
//		switch {
//		case errors.Cause(err) == omniErr.ErrResourceNotFound:
//			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
//		case errors.Cause(err) == omniErr.ErrBadRequest:
//			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//
//		default:
//			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//		}
//	}
//
//	err = httpUtils.SetResponse(ctx, http.StatusOK, relModel)
//	if err != nil {
//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	}
//
//	//idSrc, err := mapper.Get().GetSrcID(ctx.Param("user_id"), &src.User{})
//	//if err != nil {
//	//	switch {
//	//	case errors.Cause(err) == omniErr.ErrBadRequest:
//	//		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//	//	default:
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//}
//	//
//	//user, err := ctr.store.Users.GetOne(ctx.Request().Context(), idSrc)
//	//if err != nil {
//	//	switch {
//	//	case errors.Cause(err) == omniErr.ErrResourceNotFound:
//	//		return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
//	//	default:
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//}
//	//
//	//
//	//relName := ctx.Param("relation_name")
//	//switch relName {
//	//case "location":
//	//	loc, err := ctr.store.Locations.GetOne(ctx.Request().Context(), user.LocationID)
//	//	if err != nil {
//	//		switch {
//	//		case errors.Cause(err) == omniErr.ErrResourceNotFound:
//	//			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
//	//		default:
//	//			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//		}
//	//	}
//	//
//	//	web, err := loc.ToWeb()
//	//	if err != nil {
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//	err = httpUtils.SetResponse(ctx, http.StatusOK, web)
//	//	if err != nil {
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//
//	//case "roles":
//	//	srcFilters, err := filt.GetSrcFiltersFromRelationID(user.Roles)
//	//	if err != nil {
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//
//	//	srcList, err := ctr.store.Roles.GetList(ctx.Request().Context(), srcFilters)
//	//	if err != nil {
//	//		switch {
//	//		case errors.Cause(err) == omniErr.ErrResourceNotFound:
//	//			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
//	//		default:
//	//			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//		}
//	//	}
//	//
//	//	webList, err := srcList.ToWeb()
//	//	if err != nil {
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//
//	//	err = httpUtils.SetResponse(ctx, http.StatusOK, webList)
//	//	if err != nil {
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//default:
//	//	return omniErr.NewHTTPError(http.StatusForbidden, omniErr.ErrTitleResourceNotFound,
//	//		fmt.Errorf("%w wrong relation name '%v'", omniErr.ErrResourceNotFound, relName))
//	//}
//
//	return nil
//}
//
//// ModifyRelation - create, delete, replace relations
//func (ctr *UserController) ModifyRelation(ctx echo.Context) error {
//
//	//idSrc, err := mapper.Get().GetSrcID(ctx.Param("user_id"), &src.User{})
//	//if err != nil {
//	//	switch {
//	//	case errors.Cause(err) == omniErr.ErrBadRequest:
//	//		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//	//	default:
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//}
//	var domRelationModel interface{}
//	webRelName := ctx.Param("relation_name")
//	var domRelName string
//	switch webRelName {
//	case "location":
//		domRelName = "Location"
//
//		domRelationModel = new(domain.Location)
//		err := httpUtils.UnmarshalFromRequest(domRelationModel, ctx.Request().Body)
//		if err != nil {
//			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleBadRequest, err)
//		}
//		//
//		//srcModelsNew, err := new(src.Location).ScanFromWeb(webModel)
//		//if err != nil {
//		//	return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//		//}
//		//
//		//srcRelName := "Location"
//		//switch ctx.Request().Method {
//		//case http.MethodPatch:
//		//	err = ctr.store.Users.ReplaceRelation(ctx.Request().Context(), idSrc, srcRelName, srcModelsNew)
//		//case http.MethodPost:
//		//	err = ctr.store.Users.AppendRelation(ctx.Request().Context(), idSrc, srcRelName, srcModelsNew)
//		//case http.MethodDelete:
//		//	err = ctr.store.Users.DeleteRelation(ctx.Request().Context(), idSrc, srcRelName, srcModelsNew)
//		//}
//		//if err != nil {
//		//	switch {
//		//	case errors.Cause(err) == omniErr.ErrResourceNotFound:
//		//		return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
//		//	default:
//		//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//		//	}
//		//}
//
//	case "roles":
//		domRelName = "Roles"
//
//		domRecordsIntf, err := httpUtils.UnmarshalManyFromRequest(new(domain.Role), ctx.Request().Body)
//		if err != nil {
//			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleBadRequest, err)
//		}
//
//		var domModels []*domain.Role
//		err = converter.SliceI2SliceModel(domRecordsIntf, &domModels)
//		if err != nil {
//			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//		}
//
//		domRelationModel = domModels
//		//srcModelsNew, err := src.Roles.ScanFromWeb(nil, webModels)
//		//if err != nil {
//		//	return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//		//}
//		//
//		//srcRelName := "Roles"
//		//switch ctx.Request().Method {
//		//case http.MethodPatch:
//		//	err = ctr.store.Users.ReplaceRelation(ctx.Request().Context(), idSrc, srcRelName, srcModelsNew)
//		//case http.MethodPost:
//		//	err = ctr.store.Users.AppendRelation(ctx.Request().Context(), idSrc, srcRelName, srcModelsNew)
//		//case http.MethodDelete:
//		//	err = ctr.store.Users.DeleteRelation(ctx.Request().Context(), idSrc, srcRelName, srcModelsNew)
//		//}
//		//if err != nil {
//		//	switch {
//		//	case errors.Cause(err) == omniErr.ErrResourceNotFound:
//		//		return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
//		//	default:
//		//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//		//	}
//		//}
//
//	default:
//		return omniErr.NewHTTPError(http.StatusForbidden, omniErr.ErrTitleResourceNotFound,
//			fmt.Errorf("%w wrong relation name '%v'", omniErr.ErrResourceNotFound, webRelName))
//	}
//
//	err := ctr.manager.User.ModifyRelation(
//		ctx.Request().Context(),
//		ctx.Param("user_id"),
//		domRelName,
//		service.GetRelationOperFromHTTPMethod(ctx.Request().Method),
//		domRelationModel,
//	)
//	if err != nil {
//		switch {
//		case errors.Cause(err) == omniErr.ErrResourceNotFound:
//			return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
//		case errors.Cause(err) == omniErr.ErrBadRequest:
//			return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//		default:
//			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//		}
//	}
//
//	return nil
//}
//
//// Create creates user
//func (ctr *UserController) Create(ctx echo.Context) error {
//	domModel := new(domain.User)
//	err := httpUtils.UnmarshalFromRequest(domModel, ctx.Request().Body)
//	if err != nil {
//		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//	}
//
//	err = validator.ValidateStruct(domModel)
//	if err != nil {
//		return omniErr.NewHTTPError(http.StatusUnprocessableEntity, omniErr.ErrTitleValidation, err)
//	}
//
//	//srcUser, err := new(src.User).ScanFromWeb(domModel)
//	//if err != nil {
//	//	return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//	//}
//	//
//	//user, err := ctr.store.Users.Create(ctx.Request().Context(), srcUser)
//	//if err != nil {
//	//	switch {
//	//	case errors.Cause(err) == omniErr.ErrResourceExists:
//	//		return omniErr.NewHTTPError(http.StatusConflict, omniErr.ErrTitleResourceExists, err)
//	//	default:
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//}
//	//
//	//webUser, err := user.ToWeb()
//	//if err != nil {
//	//	return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//}
//	newUser, err := ctr.manager.User.Create(ctx.Request().Context(), domModel)
//	if err != nil {
//		switch {
//		case errors.Cause(err) == omniErr.ErrResourceExists:
//			return omniErr.NewHTTPError(http.StatusConflict, omniErr.ErrTitleResourceExists, err)
//		default:
//			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//		}
//	}
//
//	err = httpUtils.SetResponse(ctx, http.StatusOK, newUser)
//	if err != nil {
//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	}
//
//	return nil
//}
//
//// Update updates user attributes
//func (ctr *UserController) Update(ctx echo.Context) error {
//
//	domModel := new(domain.User)
//	err := httpUtils.UnmarshalFromRequest(domModel, ctx.Request().Body)
//	if err != nil {
//		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//	}
//
//	err = validator.ValidateStruct(domModel)
//	if err != nil {
//		return omniErr.NewHTTPError(http.StatusUnprocessableEntity, omniErr.ErrTitleValidation, err)
//	}
//
//	//srcUser, err := new(src.User).ScanFromWeb(domModel)
//	//if err != nil {
//	//	return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//	//}
//	//
//	//user, err := ctr.store.Users.Update(ctx.Request().Context(), srcUser)
//	//if err != nil {
//	//	switch {
//	//	case errors.Cause(err) == omniErr.ErrResourceNotFound:
//	//		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//	//	default:
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//}
//	//
//	//webUser, err := user.ToWeb()
//	//if err != nil {
//	//	return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//}
//
//	newUser, err := ctr.manager.User.Update(ctx.Request().Context(), domModel)
//	if err != nil {
//		switch {
//		case errors.Cause(err) == omniErr.ErrResourceExists:
//			return omniErr.NewHTTPError(http.StatusConflict, omniErr.ErrTitleResourceExists, err)
//		default:
//			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//		}
//	}
//
//	err = httpUtils.SetResponse(ctx, http.StatusOK, newUser)
//	if err != nil {
//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	}
//
//	return nil
//
//}
//
//// Delete deletes user
//func (ctr *UserController) Delete(ctx echo.Context) error {
//	//idSrc, err := mapper.Get().GetSrcID(ctx.Param("user_id"), &src.User{})
//	//if err != nil {
//	//	switch {
//	//	case errors.Cause(err) == omniErr.ErrBadRequest:
//	//		return omniErr.NewHTTPError(http.StatusBadRequest, omniErr.ErrTitleResourceNotFound, err)
//	//	default:
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//}
//	//
//	//err = ctr.store.Users.Delete(ctx.Request().Context(), idSrc)
//	//if err != nil {
//	//	switch {
//	//	case errors.Cause(err) == omniErr.ErrResourceNotFound:
//	//		return omniErr.NewHTTPError(http.StatusNotFound, omniErr.ErrTitleResourceNotFound, err)
//	//	default:
//	//		return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//	//	}
//	//}
//
//	err := ctr.manager.User.Delete(ctx.Request().Context(), ctx.Param("user_id"))
//	if err != nil {
//		switch {
//		case errors.Cause(err) == omniErr.ErrResourceExists:
//			return omniErr.NewHTTPError(http.StatusConflict, omniErr.ErrTitleResourceExists, err)
//		default:
//			return omniErr.NewHTTPError(http.StatusInternalServerError, omniErr.ErrTitleInternal, err)
//		}
//	}
//
//	ctx.NoContent(http.StatusNoContent)
//	return nil
//}
