package main

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/fx"
	"net/http"
	"time"
	"uroborus/common"
	"uroborus/common/logging"
	settings "uroborus/common/setting"

	"uroborus/router"
	server "uroborus/server/fx"
	service "uroborus/service/fx"
	store "uroborus/store/fx"
)

// ServiceLifetimeHooks -
func ServiceLifetimeHooks(lc fx.Lifecycle, srv *http.Server, logger *logging.ZapLogger) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Sugar().Info("starting web server listen and serve at ", srv.Addr, " ...")
				go func() {
					if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						logger.Sugar().Fatalf("listen: %s\n", err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Sugar().Info("closing web server ...")
				ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				return srv.Shutdown(ctx)
			},
		},
	)
}

func main() {
	port := flag.String("p", "8082", "port to listen on")
	flag.Parse()
	app := fx.New(
		fx.Provide(fx.Annotated{
			Name: "addr",
			Target: func() string {
				return fmt.Sprintf(":%s", *port)
			},
		}),
		fx.Provide(func() context.Context {
			return context.Background()
		}),
		settings.Module,
		common.Module,
		fx.Provide(logging.NewZapLogger),
		store.Module,
		service.Module,
		server.Module,
		router.Module,
		fx.Invoke(ServiceLifetimeHooks),
	)
	app.Run()
}
