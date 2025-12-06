package api

import (
	"onboarding/internal/handler"
	"onboarding/pkg/config"
	"onboarding/pkg/token"
	"onboarding/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	router                *gin.Engine
	jwtImpl               token.JWT
	authHandler           *handler.AuthHandler
	userHandler           *handler.UserHandler
	forgotPasswordHandler *handler.ForgotPasswordHandler
}

func NewServer(
	cfg config.App,
	jwtImpl token.JWT,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	forgotPasswordHandler *handler.ForgotPasswordHandler,
) *Server {
	server := &Server{
		jwtImpl:               jwtImpl,
		authHandler:           authHandler,
		userHandler:           userHandler,
		forgotPasswordHandler: forgotPasswordHandler,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validEmail", validation.ValidEmail)
		v.RegisterValidation("validPassword", validation.ValidPassword)
		v.RegisterValidation("validUUID", validation.ValidUUID)
	}

	server.setupRouter(cfg)
	return server
}

func (server *Server) Start(port string) error {
	return server.router.Run(port)
}

func (server *Server) setupRouter(cfg config.App) {
	gin.SetMode(cfg.GinMode)
	router := gin.Default()

	formRoutes := router.Group("/").Use(
		ContentTypeValidation(),
		Timeout(cfg.Timeout),
	)
	{
		formRoutes.POST("/auth/register", server.authHandler.Register)
		formRoutes.POST("/auth/login", server.authHandler.Login)
		formRoutes.POST("/forgot-password", server.forgotPasswordHandler.RequestResetPassword)
		formRoutes.POST("/reset-password", server.forgotPasswordHandler.ResetPassword)
	}

	authRoutes := router.Group("/").Use(
		Authentication(server.jwtImpl),
		Timeout(cfg.Timeout),
	)
	{
		authRoutes.DELETE("/auth/logout", server.authHandler.Logout)
	}

	authFormRoutes := router.Group("/").Use(
		ContentTypeValidation(),
		Authentication(server.jwtImpl),
		Timeout(cfg.Timeout),
	)
	{
		authFormRoutes.GET("/user", server.userHandler.GetUser)
		authFormRoutes.GET("/user/:uuid", server.userHandler.GetUser)
	}

	server.router = router
}

func (server *Server) start(port string) error {
	return server.router.Run(port)
}
