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
	deployServer    *server.DeployServer
	containerServer *server.ContainerServer
	groupServer     *server.GroupServer
	commentServer   *server.ProjectCommentServer
}

// NewRouter Generator
func NewRouter(
	config *settings.Config,
	logger *logging.ZapLogger,
	healthServer *server.HealthServer,
	userServer *server.UserServer,
	projectServer *server.ProjectServer,
	baseImageServer *server.BaseImageServer,
	deployServer *server.DeployServer,
	containerServer *server.ContainerServer,
	groupServer *server.GroupServer,
	commentServer *server.ProjectCommentServer,
) *Router {
	return &Router{
		config:          config,
		logger:          logger,
		healthServer:    healthServer,
		userServer:      userServer,
		projectServer:   projectServer,
		baseImageServer: baseImageServer,
		deployServer:    deployServer,
		containerServer: containerServer,
		groupServer:     groupServer,
		commentServer:   commentServer,
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
		}
		{
			userRoute := app.Group(baseEngine.BasePath() + "/user")
			userRoute.PUT("", r.userServer.Register)
			userRoute.POST("", r.userServer.Login)
		}
		{
			projectRoute := app.Group(baseEngine.BasePath() + "/project")
			projectRoute.Use(middleware.Auth())
			projectRoute.GET("", r.projectServer.Get)
			projectRoute.PUT("", r.projectServer.Register)
			projectRoute.DELETE("", r.projectServer.Delete)
			projectRoute.GET("/group", r.projectServer.GetGroupProjs)
			projectRoute.POST("/group", r.projectServer.RegisterGroupProj)
			projectRoute.POST("/info", r.projectServer.SaveProjectInfo)
		}
		{
			deployRoute := app.Group(baseEngine.BasePath() + "/deploy")
			deployRoute.GET("/log/ws", r.deployServer.Log)
			deployRoute.GET("/running/ws", r.deployServer.RunningLog)
			deployRoute.Use(middleware.Auth())
			deployRoute.GET("", r.deployServer.Get)
			deployRoute.POST("", r.deployServer.Deploy)

		}
		{
			baseImageRoute := app.Group(baseEngine.BasePath() + "/image")
			baseImageRoute.Use(middleware.Auth())
			baseImageRoute.PUT("", r.baseImageServer.Register)
			baseImageRoute.GET("", r.baseImageServer.Get)
		}
		{
			containerRoute := app.Group(baseEngine.BasePath() + "/container")
			containerRoute.GET("/exec", r.containerServer.Exec)
			containerRoute.Use(middleware.Auth())
			containerRoute.GET("", r.containerServer.GetAll)
		}
		{
			groupRoute := app.Group(baseEngine.BasePath() + "/group")
			groupRoute.Use(middleware.Auth())
			groupRoute.GET("", r.groupServer.Find)
			groupRoute.POST("", r.groupServer.Register)
			groupRoute.DELETE("", r.groupServer.Delete)
		}
		{
			commentRoute := app.Group(baseEngine.BasePath() + "/comment")
			commentRoute.Use(middleware.Auth())
			commentRoute.GET("", r.commentServer.Find)
			commentRoute.POST("", r.commentServer.Register)
		}

	}
	return app
}
