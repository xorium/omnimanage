package error

import (
	"errors"
	"fmt"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"net/http"
	httpUtils "omnimanage/pkg/utils/http"
	"strconv"
)

type InternalErrObj struct {
	ID    string
	Code  string
	Title string
	Err   error
}

const (
	ErrCodeUnknown  = "UNKNOWN_ERROR"
	ErrCodeResource = "RESOURCE_ERROR"
	ErrCodeInternal = "INTERNAL_ERROR"

	ErrTitleUnknown          = "UNKNOWN_ERROR"
	ErrTitleNoAuth           = "AUTH_REQUIRED"
	ErrTitleResourceNotFound = "RESOURCE_OBJECT_NOT_FOUND"
	ErrTitleInternal         = "INTERNAL"
)

var (
	ErrResourceNotFound = errors.New("resource not found")
	ErrNoAuth           = errors.New("authorization required")
	ErrBadRequest       = errors.New("bad request")
	ErrInternal         = errors.New("internal error")
)

func PutErrorJSONAPI(errObj *jsonapi.ErrorObject) []*jsonapi.ErrorObject {
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

		httpUtils.MarshalToResponse(PutErrorJSONAPI(errObj), ctx.Response())
		return
	}

	fmt.Printf("echo Err: %v", echoErr)
	//errEcho.

	//errObj := HTTPError{
	//	Code:    http.StatusInternalServerError,
	//	Message: err.ErrHandler(),
	//}
	//
	//switch err {
	//case types.ErrBadRequest:
	//	errObj.Code = http.StatusBadRequest
	//case types.ErrNotFound:
	//	errObj.Code = http.StatusNotFound
	//case types.ErrDuplicateEntry, types.ErrConflict:
	//	errObj.Code = http.StatusConflict
	//case types.ErrForbidden:
	//	errObj.Code = http.StatusForbidden
	//case types.ErrUnprocessableEntity:
	//	errObj.Code = http.StatusUnprocessableEntity
	//case types.ErrPartialOk:
	//	errObj.Code = http.StatusPartialContent
	//case types.ErrGone:
	//	errObj.Code = http.StatusGone
	//case types.ErrUnauthorized:
	//	errObj.Code = http.StatusUnauthorized
	//}
	//he, ok := err.(*echo.HTTPError)
	//if ok {
	//	errObj.Code = he.Code
	//	errObj.Message = fmt.Sprintf("%v", he.Message)
	//}
	//errObj.UserName = http.StatusText(errObj.Code)
	//if !ctx.Response().Committed {
	//	if ctx.Request().Method == echo.HEAD {
	//		ctx.NoContent(errObj.Code)
	//	} else {
	//		ctx.JSON(errObj.Code, errObj)
	//	}
	//}
}
