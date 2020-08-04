package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
    "github.com/natefinch/lumberjack"
)


//init logger
func InitLogger() *zap.Logger{
    writeSyncer := getLogWriter()
    encoder := getEncoder()
    core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
    // logger := zap.New(core, zap.AddCaller())
    return zap.New(core, zap.AddCaller())
}

//get log writer
func getLogWriter() zapcore.WriteSyncer {
    logMaxSize := 10
    logMaxBackups := 200
    logMaxAge := 30

    lumberJackLogger := &lumberjack.Logger{
        Filename:   "logs/opentracing/go/http_server_trace.log",
        MaxSize:    logMaxSize,
        MaxBackups: logMaxBackups,
        MaxAge:     logMaxAge,
        Compress:   false,
    }
    return zapcore.AddSync(lumberJackLogger)
}

//get log encoder
func getEncoder() zapcore.Encoder {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
    return zapcore.NewConsoleEncoder(encoderConfig)
}