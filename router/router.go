package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
	"uroborus/common/logging"
	settings "uroborus/common/setting"
	"uroborus/router/middleware"
	"uroborus/server"
	"uroborus/server/doc/swagger"
)

// Router router
type Router struct {
	config          *settings.Config
	logger          *logging.ZapLogger
	healthServer    *server.HealthServer
	userServer      *server.UserServer
	projectServer   *server.ProjectServer
	baseImageServer *server.BaseImageServer
}

// NewRouter Generator
func NewRouter(
	config *settings.Config,
	logger *logging.ZapLogger,
	healthServer *server.HealthServer,
	userServer *server.UserServer,
	projectServer *server.ProjectServer,
	baseImageServer *server.BaseImageServer,
) *Router {
	return &Router{
		config:          config,
		logger:          logger,
		healthServer:    healthServer,
		userServer:      userServer,
		projectServer:   projectServer,
		baseImageServer: baseImageServer,
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
func (r *Router) Server(middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.DisableConsoleColor()
	app := gin.New()
	// Setup middlewares
	app.Use(middlewares...)
	// Api router
	app.GET("/swagger.json", swagger.Serve)
	{
		baseEngine := app.Group(r.config.ApiPrefix + ApiV1)
		{
			baseEngine.GET("/health", r.healthServer.CheckV1)
			baseEngine.Use(middleware.Auth())
			baseEngine.GET("/health/auth", r.healthServer.CheckV1)

		}
		{
			userRoute := app.Group(baseEngine.BasePath() + "/user")
			userRoute.PUT("", r.userServer.Register)
			userRoute.POST("", r.userServer.Login)
		}
		{
			projectRoute := app.Group(baseEngine.BasePath() + "/project")
			projectRoute.Use(middleware.Auth())
			projectRoute.PUT("", r.projectServer.Register)
			projectRoute.POST("/checkout", r.projectServer.CheckOut)
		}
		{
			baseImageRoute := app.Group(baseEngine.BasePath() + "/image")
			baseImageRoute.Use(middleware.Auth())
			baseImageRoute.PUT("", r.baseImageServer.Register)
			baseImageRoute.GET("", r.baseImageServer.Get)
		}
	}
	return app
}
