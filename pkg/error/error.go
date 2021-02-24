package error

import (
	"errors"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"net/http"
	httpUtils "omnimanage/pkg/utils/http"
	"strconv"
)

const (
	ErrCodeUnknown  = "UNKNOWN_ERROR"
	ErrCodeResource = "RESOURCE_ERROR"
	ErrCodeInternal = "INTERNAL_ERROR"

	ErrTitleUnknown          = "UNKNOWN_ERROR"
	ErrTitleNoAuth           = "AUTH_REQUIRED"
	ErrTitleResourceNotFound = "RESOURCE_OBJECT_NOT_FOUND"
	ErrTitleResourceExists   = "RESOURCE_OBJECT_EXISTS"
	ErrTitleInternal         = "INTERNAL"
)

var (
	ErrResourceNotFound = errors.New("resource not found")
	ErrResourceExists   = errors.New("resource already exists")
	ErrNoAuth           = errors.New("authorization required")
	ErrBadRequest       = errors.New("bad request")
	ErrInternal         = errors.New("internal error")
)

func MakeSliceJSONAPI(errObj *jsonapi.ErrorObject) []*jsonapi.ErrorObject {
	errs := make([]*jsonapi.ErrorObject, 0, 1)
	errs = append(errs, errObj)
	return errs
}

func NewHTTPError(ctx echo.Context, status int, code string, title string, err error) *echo.HTTPError {
	rid := ctx.Request().Header.Get(echo.HeaderXRequestID)
	errObj := &HTTPErrorObj{
		ID:    rid,
		Code:  code,
		Title: title,
		Err:   err,
	}

	return echo.NewHTTPError(status, errObj)
}

type HTTPErrorObj struct {
	ID    string
	Code  string
	Title string
	Err   error
}

// ErrHandler implements a custom echo error handler
func ErrHandler(err error, ctx echo.Context) {

	errObj := &jsonapi.ErrorObject{}

	echoErr, ok := err.(*echo.HTTPError)
	if !ok {
		errObj.Status = strconv.Itoa(http.StatusInternalServerError)
		errObj.Code = ErrCodeUnknown
		errObj.Title = ErrTitleUnknown
		errObj.Detail = err.Error()

		errToResponse(ctx, http.StatusInternalServerError, errObj)
		return
	}

	switch errInt := echoErr.Message.(type) {
	case *HTTPErrorObj:
		errObj.ID = errInt.ID
		errObj.Status = strconv.Itoa(echoErr.Code)
		errObj.Code = errInt.Code
		errObj.Title = errInt.Title
		errObj.Detail = errInt.Err.Error()

	default:

	}

	errToResponse(ctx, echoErr.Code, errObj)
}

func errToResponse(ctx echo.Context, code int, errObj *jsonapi.ErrorObject) {
	if !ctx.Response().Committed {
		if ctx.Request().Method == echo.HEAD {
			ctx.NoContent(code)
		} else {
			httpUtils.SetResponseError(ctx, code, MakeSliceJSONAPI(errObj))
		}
	}
}
