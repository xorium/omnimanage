package http

import (
	"bytes"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"io"
)

func SetResponse(ctx echo.Context, code int, model interface{}) error {
	var b bytes.Buffer
	err := jsonapi.MarshalPayload(&b, model)
	if err != nil {
		return err
	}

	err = ctx.JSONBlob(code, b.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func SetResponseError(ctx echo.Context, code int, errObjs []*jsonapi.ErrorObject) error {
	var b bytes.Buffer
	err := jsonapi.MarshalErrors(&b, errObjs)
	if err != nil {
		return err
	}

	err = ctx.JSONBlob(code, b.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func UnmarshalFromRequest(model interface{}, r io.Reader) error {
	err := jsonapi.UnmarshalPayload(r, model)
	if err != nil {
		return err
	}
	return nil
}
