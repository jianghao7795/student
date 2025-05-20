package errors

import (
	"fmt"

	errors "github.com/go-kratos/kratos/v2/errors"
)

func Error400(err error) error {
	return errors.New(400, "CONTENT_MISSING", fmt.Sprintf("bad request: %s", err.Error())).WithMetadata(map[string]string{
		"reason": err.Error(),
	})
}

func Error401(err error) error {
	return errors.New(401, "UNAUTHORIZED", fmt.Sprintf("unauthorized: %s", err.Error())).WithMetadata(map[string]string{
		"reason": err.Error(),
	})
}

func Error403() error {
	return errors.New(403, "FORBIDDEN", "forbidden")
}

func Error404() error {
	return errors.New(404, "NOT_FOUND", "record not found").WithMetadata(map[string]string{
		"reason": "record not found",
	})
}
