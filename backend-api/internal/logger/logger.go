package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Init() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.TimeKey = "time"
	config.CallerKey = "caller"
	config.MessageKey = "msg"
	config.LevelKey = "level"
	config.EncodeLevel = zapcore.CapitalLevelEncoder

	// 确保 logs 目录存在
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}
	// 按天分割日志文件，文件名为 logs/yyyy-MM-dd.log
	logFileName := time.Now().Format("2006-01-02") + ".log"
	logFilePath := "logs/" + logFileName
	logFile, _ := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	encoder := zapcore.NewConsoleEncoder(config)
	core := zapcore.NewCore(encoder, zapcore.AddSync(logFile), zap.InfoLevel)
	Logger = zap.New(core, zap.AddCaller())
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}
