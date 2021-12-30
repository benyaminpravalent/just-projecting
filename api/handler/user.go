package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mine/just-projecting/service"
)

type UserHandler struct {
	user service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (handler *UserHandler) SetUserService(service service.UserService) *UserHandler {
	handler.user = service
	return handler
}

func (handler *UserHandler) Validate() *UserHandler {
	if handler.user == nil {
		panic("handler need user service")
	}
	return handler
}

func (handler *UserHandler) GetUsers(c *fiber.Ctx) error {
	return c.SendString("users")
}

func (handler *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.SendString(fmt.Sprintf("user: %s", id))
}

func (handler *UserHandler) CountUsers(c *fiber.Ctx) error {
	id := c.Params("id")

	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic("error while parsing")
	}

	res, err := handler.user.CountUsersUsecase(n)
	if err != nil {
		panic("failed when call usecase")
	}

	// return c.SendString(fmt.Sprintf("count user: %s", res))
	return c.JSON(res)
}

func (handler *UserHandler) PublishMessage(c *fiber.Ctx) error {
	return handler.user.PublishMessageUser(c.Context())
}
