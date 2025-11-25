package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/OzyKleyton/agendamento-api/config"
	"github.com/OzyKleyton/agendamento-api/config/db"
	"github.com/OzyKleyton/agendamento-api/internal/api/handler"
	"github.com/OzyKleyton/agendamento-api/internal/api/middleware"
	"github.com/OzyKleyton/agendamento-api/internal/api/router"
	"github.com/OzyKleyton/agendamento-api/internal/model/user"
	"github.com/OzyKleyton/agendamento-api/internal/repository"
	"github.com/OzyKleyton/agendamento-api/internal/service"
	jwt "github.com/OzyKleyton/agendamento-api/utils/auth"
	"github.com/gofiber/fiber/v2"
)

func Run(host, port string) error {
	address := fmt.Sprintf("%s:%s", host, port)
	log.Println("Listen app in port ", address)

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		Prefork:     config.GetConfig().Prefork,
		ProxyHeader: fiber.HeaderXForwardedFor,
	})

	db, err := db.ConnectDB(config.GetConfig().DBURL)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db = db.WithContext(ctx)

	if err := db.AutoMigrate(
		&user.User{},
	); err != nil {
		return err
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)

	// JWT
	jwtUtil := jwt.NewJWT()

	// Services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, jwtUtil)

	// Handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Setup das rotas seguindo sua estrutura
	setupRoutes(app, authMiddleware, authHandler, userHandler)

	c := make(chan os.Signal, 1)
	errc := make(chan error, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		cancel()
		errc <- app.Shutdown()
	}()

	if err := app.Listen(address); err != nil {
		return err
	}

	err = <-errc

	return err
}

func setupRoutes(app *fiber.App, authMiddleware *middleware.AuthMiddleware, authHandler *handler.AuthHandler, userHandler *handler.UserHandler) {
	allRoutes := func(route fiber.Router) {
		authHandler.Routes()(route)

		protected := route.Group("")
		protected.Use(authMiddleware.Validate())
		userHandler.Routes()(protected)
	}

	router.SetupRouter(app, allRoutes)
}
