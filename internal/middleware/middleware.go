package middleware

import (
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
)

func ResponseType(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if ctx.Path() == "/docs/" {
			return next(ctx)
		}
		ctx.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
		return next(ctx)
	}
}
