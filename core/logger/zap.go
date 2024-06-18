package logger

import (
	"fastApi/core/global"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

const TraceId = "traceId"

var logger *zap.Logger

func InitLogger() *zap.Logger {

	logLevel := zapcore.InfoLevel
	switch viper.GetString("logger.level") {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	}

	zapCore := zapcore.NewCore(
		getEncoder(),
		getLogWriter(),
		logLevel,
	)

	global.Log = zap.New(zapCore)
	logger = global.Log
	return logger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create(viper.GetString("logger.path"))
	ws := io.MultiWriter(file, os.Stdout) // 打印到控制台和文件
	return zapcore.AddSync(ws)
}

func CalcTraceId() (traceId string) {
	return uuid.New().String()
}

func With(fields ...zap.Field) {
	global.Log = logger.With(fields...)
	global.SLog = global.Log.Sugar()
	global.DB.Logger = NewGormLog(global.Log)
}
