package middleware

import "go.uber.org/fx"

// Module export all gin.HandlerFunc type middlewares, due to
// there are sample type, when dependency injection could be
// injected as group with annotated group name, receiver also
// need annotated group.
var Module = fx.Options(
	fx.Provide(fx.Annotated{
		Target: CORS,
		Group:  "middleware",
	}),
	fx.Provide(fx.Annotated{
		Target: Error,
		Group:  "middleware",
	}),
	fx.Provide(fx.Annotated{
		Target: Log,
		Group:  "middleware",
	}),
)
