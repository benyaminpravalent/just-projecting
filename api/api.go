package api

import (
	"context"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"

	"github.com/mine/just-projecting/api/handler"
	"github.com/mine/just-projecting/api/helper"
	"github.com/mine/just-projecting/domain/repository"
	"github.com/mine/just-projecting/pkg/config"
	log "github.com/mine/just-projecting/pkg/logger"
	"github.com/mine/just-projecting/service"
)

//Server API server
type Server interface {
	Configure(ctx context.Context) error
	Serve() error
}

//FiberServer gofiber api server
type FiberServer struct {
	app *fiber.App
}

var l *logrus.Entry

//NewFiberServer create instance of fiber server
func NewFiberServer() Server {
	return &FiberServer{
		app: fiber.New(fiber.Config{
			ErrorHandler: helper.Error,
		}),
	}
}

//Configure configure fiber server
func (f *FiberServer) Configure(ctx context.Context) error {
	l = log.GetLoggerContext(ctx, "api", "Configure")

	if err := f.registerHandlers(); err != nil {
		l.WithError(err).Error("Error registering handler")
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_ = f.app.Shutdown()
	}()

	return nil
}

//Serve start fiber server
func (f *FiberServer) Serve() error {
	return f.app.Listen(config.GetString("address") + ":" + config.GetString("port"))
}

func (f *FiberServer) registerHandlers() error {
	f.app.Use(recover.New())
	f.app.Use(logger.New())
	f.app.Use(cors.New())

	f.app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	userRepo := repository.NewUserRepository()

	userService := service.NewUserService().
		SetUserRepository(userRepo).
		Validate()

	userHandler := handler.NewUserHandler().
		SetUserService(userService).
		Validate()

	user := f.app.Group("/users")
	user.Get("/", userHandler.GetUsers)
	user.Get("/count/:id", userHandler.CountUsers)
	user.Get("/:id", userHandler.GetUser)
	user.Post("/publish", userHandler.PublishMessage)

	return nil
}
