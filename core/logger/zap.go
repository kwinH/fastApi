package logger

import (
	"context"
	"fastApi/core/global"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"io"
	"os"
)

const loggerKey = "Log"
const loggerSugarKey = "LogSugar"
const TraceId = "traceId"
const SpanId = "spanId"

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
	global.SLog = global.Log.Sugar()
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

func CalcTraceId(ctx context.Context) (traceId, spanId string) {
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID()
	if traceID.IsValid() {
		return traceID.String(), span.SpanContext().SpanID().String()
	}

	return uuid.New().String(), ""
}

func injectToContext(log *zap.Logger, slog *zap.SugaredLogger, db *gorm.DB, store func(key string, value interface{})) {
	store(loggerKey, log)
	store(loggerSugarKey, slog)
	store(global.DBKey, db)
}

func createLoggerAndDB(fields ...zap.Field) (*zap.Logger, *zap.SugaredLogger, *gorm.DB) {
	log := logger.With(fields...)
	slog := log.Sugar()
	db := global.GDB
	db.Logger = NewGormLog(log)
	return log, slog, db
}

func With(c *gin.Context, fields ...zap.Field) {
	log, slog, db := createLoggerAndDB(fields...)
	injectToContext(log, slog, db, c.Set)
}

func WithC(c context.Context, fields ...zap.Field) context.Context {
	log, slog, db := createLoggerAndDB(fields...)
	injectToContext(log, slog, db, func(key string, value interface{}) {
		c = context.WithValue(c, key, value)
	})
	return c
}

func Log(c context.Context) *zap.Logger {
	return c.Value(loggerKey).(*zap.Logger)
}

func SLog(c context.Context) *zap.SugaredLogger {
	return c.Value(loggerSugarKey).(*zap.SugaredLogger)
}
