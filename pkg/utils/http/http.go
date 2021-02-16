package http

import (
	"bytes"
	"github.com/google/jsonapi"
	"io"
)

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
