package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	pkgerror "github.com/mine/just-projecting/pkg/error"
)

//QueryValidate validate query params
func QueryValidate(qs ...string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		valid := false
		for _, q := range qs {
			if c.Query(q) != "" {
				valid = true
				break
			}
		}

		if valid {
			return c.Next()
		}

		err := fmt.Errorf("one of [%s] is required", strings.Join(qs, ", "))
		return pkgerror.BadRequest(err)
	}
}
