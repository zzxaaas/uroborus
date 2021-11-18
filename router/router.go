package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
	"uroborus/common/logging"
	settings "uroborus/common/setting"
	"uroborus/server"
	"uroborus/server/doc/swagger"
)

// Router router
type Router struct {
	config       *settings.Config
	logger       *logging.ZapLogger
	healthServer *server.HealthServer
}

// NewRouter Generator
func NewRouter(
	config *settings.Config,
	logger *logging.ZapLogger,
	healthServer *server.HealthServer,
) *Router {
	return &Router{
		config:       config,
		logger:       logger,
		healthServer: healthServer,
	}
}

// ServerOption fx需要
type ServerOption struct {
	fx.In
	Addr       string            `name:"addr"`
	Middleware []gin.HandlerFunc `group:"middleware"`
}

// NewHTTPServer fx需要
func NewHTTPServer(router *Router, option ServerOption) *http.Server {
	return &http.Server{
		Addr:    option.Addr,
		Handler: router.Server(option.Middleware...),
	}
}

// Server main server
func (r *Router) Server(middleware ...gin.HandlerFunc) *gin.Engine {
	gin.DisableConsoleColor()
	app := gin.New()
	// Setup middlewares
	app.Use(middleware...)
	// Api router
	app.GET("/swagger.json", swagger.Serve)
	{
		baseEngine := app.Group(r.config.ApiPrefix + ApiV1)
		{
			baseEngine.GET("/health", r.healthServer.CheckV1)
		}
	}
	return app
}
