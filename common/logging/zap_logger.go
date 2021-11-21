package logging

import (
	"go.uber.org/zap"
	settings "uroborus/common/setting"
)

// ZapLogger -
type ZapLogger struct {
	*zap.Logger
}

// NewZapLogger -
func NewZapLogger(config *settings.Config) *ZapLogger {
	var logger *zap.Logger
	var err error
	if config.InReleaseMode() {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
	return &ZapLogger{logger}
}
