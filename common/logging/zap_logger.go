package logging

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	settings "uroborus/common/setting"
)

// ZapLogger -
type ZapLogger struct {
	*zap.Logger
}

// NewZapLogger -
func NewZapLogger(config *settings.Config) *ZapLogger {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		//EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, //这里可以指定颜色
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	// 设置日志级别
	var atom zap.AtomicLevel
	if config.InReleaseMode() {
		atom = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		atom = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	zapConfig := zap.Config{
		Level:       atom, // 日志级别
		Development: true, // 开发模式，堆栈跟踪
		//Encoding:         "json",                                              // 输出格式 console 或 json
		Encoding:         "console",          // 输出格式 console 或 json
		EncoderConfig:    encoderConfig,      // 编码器配置
		OutputPaths:      []string{"stdout"}, // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
	}
	zapConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder

	logger, err := zapConfig.Build()
	if err != nil {
		panic(fmt.Sprintf("log 初始化失败: %v", err))
	}

	return &ZapLogger{logger}
}
