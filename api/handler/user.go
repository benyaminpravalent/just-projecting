package handler

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mine/just-projecting/domain/model"
	"github.com/mine/just-projecting/service"
	"github.com/prometheus/common/log"
)

type UserHandler struct {
	user service.UserService
}

type DecodeRes struct {
	Username string `json:"user_name"`
	Exp      int64  `json:"exp"`
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

func (handler *UserHandler) Login(c *fiber.Ctx) error {
	var form model.Login

	body := c.Body()
	if len(body) < 0 {
		panic("Validation Error - Body Required")
	}

	if err := json.Unmarshal(body, &form); err != nil {
		return c.JSON(http.StatusInternalServerError)
	}

	token, userData, err := handler.user.LoginUsecase(c, form)
	if err != nil {
		return c.JSON(http.StatusInternalServerError)
	}

	response := model.LoginResponse{
		Token: token,
		User:  userData,
	}

	return c.JSON(response)
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

func (handler *UserHandler) GetMerchantOmzet(c *fiber.Ctx) error {
	var token model.Token

	date := c.Query("date", "2021-11-01")
	temp := c.Request().Header.Peek("Authorization")

	token.Token = string(temp[:])

	jwt, err := handler.DecodeToken(token.Token)
	if err != nil {
		return err
	}

	res, err := handler.user.ListMerchantOmzet(c, date, jwt.Username)
	if err != nil {
		panic("failed when call usecase")
	}

	return c.JSON(res)
}

func (handler *UserHandler) GetOutletOmzet(c *fiber.Ctx) error {
	var token model.Token

	date := c.Query("date", "2021-11-01")
	temp := c.Request().Header.Peek("Authorization")

	if err := json.Unmarshal(temp, &token.Token); err != nil {
		return c.JSON(http.StatusInternalServerError)
	}

	jwt, err := handler.DecodeToken(token.Token)
	if err != nil {
		return err
	}

	res, err := handler.user.ListOutletOmzet(c, date, jwt.Username)
	if err != nil {
		panic("failed when call usecase")
	}

	return c.JSON(res)
}

func (handler *UserHandler) DecodeToken(tokenStr string) (decodeRes DecodeRes, err error) {
	payload := strings.Split(tokenStr, ".")

	if len(payload) < 3 {
		log.Error(fmt.Sprintf("token is not valid: %s", errors.New("Token is Invalid")))
		return decodeRes, errors.New("Token is Invalid")
	}

	byte, err := b64.RawStdEncoding.DecodeString(payload[1])
	if err != nil {
		log.Error(fmt.Sprintf("failed decode token string: %s", err.Error()))
		return decodeRes, err
	}

	err = json.Unmarshal(byte, &decodeRes)
	if err != nil {
		log.Error(fmt.Sprintf("failed unmarshal for claims: %s", err.Error()))
	}

	return decodeRes, err
}
