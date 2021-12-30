package helper

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/mine/just-projecting/api/presenter"
	pkgerror "github.com/mine/just-projecting/pkg/error"
)

//Error set error
func Error(ctx *fiber.Ctx, err error) error {
	var code int = http.StatusInternalServerError
	var message string

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	} else {
		if e, ok := err.(*pkgerror.Error); ok {
			code = e.Code
			message = e.Err.Error()
		} else {
			message = err.Error()
		}
	}

	if message == "" {
		message = http.StatusText(code)
	}

	hte := presenter.HTTPError{
		Code:    code,
		Message: message,
	}

	ctx.Status(code)
	return ctx.JSON(hte)
}
