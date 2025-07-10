package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger *zap.Logger
	mu     sync.Mutex
	currentDate string
)

// 日志轮转配置
type LogConfig struct {
	// 日志文件路径
	Filename string
	// 单个日志文件最大大小，单位MB
	MaxSize int
	// 保留的旧日志文件最大数量
	MaxBackups int
	// 保留的旧日志文件最大天数
	MaxAge int
	// 是否压缩旧日志文件
	Compress bool
	// 日志级别
	Level string
}

// 默认日志配置
var defaultLogConfig = LogConfig{
	Filename:    "logs/main.log",
	MaxSize:     100, // 100MB
	MaxBackups:  30,  // 保留30个备份文件
	MaxAge:      30,  // 保留30天
	Compress:    true, // 压缩旧文件
	Level:       "info",
}

func Init() {
	InitWithConfig(defaultLogConfig)
}

func InitWithConfig(config LogConfig) {
	mu.Lock()
	defer mu.Unlock()

	// 确保 logs 目录存在
	logDir := filepath.Dir(config.Filename)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}

	// 获取当前日期
	currentDate = time.Now().Format("2006-01-02")

	// 配置 lumberjack 日志轮转
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	// 配置 zap 编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.TimeKey = "time"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "msg"
	encoderConfig.LevelKey = "level"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 创建编码器
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 解析日志级别
	level := parseLogLevel(config.Level)

	// 创建核心
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(lumberJackLogger),
		level,
	)

	// 创建 logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	
	// 直接使用 Logger 记录初始化信息，避免递归调用
	Logger.Info("日志系统初始化完成", 
		zap.String("filename", config.Filename),
		zap.Int("maxSize", config.MaxSize),
		zap.Int("maxBackups", config.MaxBackups),
		zap.Int("maxAge", config.MaxAge),
		zap.Bool("compress", config.Compress),
		zap.String("level", config.Level),
		zap.String("currentDate", currentDate))
}

// 检查并执行日志轮转
func checkAndRotate() {
	mu.Lock()
	defer mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if currentDate != today {
		// 日期变化，需要轮转日志
		rotateLogFile()
		currentDate = today
	}
}

// 轮转日志文件
func rotateLogFile() {
	if Logger == nil {
		return
	}

	// 获取当前日志文件路径
	logPath := defaultLogConfig.Filename
	logDir := filepath.Dir(logPath)
	logName := filepath.Base(logPath)
	
	// 生成归档文件名
	archiveName := fmt.Sprintf("%s.%s", logName, currentDate)
	archivePath := filepath.Join(logDir, archiveName)

	// 如果 main.log 存在且不为空，则重命名为归档文件
	if fileInfo, err := os.Stat(logPath); err == nil && fileInfo.Size() > 0 {
		// 如果归档文件已存在，先删除
		if _, err := os.Stat(archivePath); err == nil {
			os.Remove(archivePath)
		}
		
		// 重命名当前日志文件为归档文件
		if err := os.Rename(logPath, archivePath); err == nil {
			// 直接使用 Logger 记录，避免递归调用
			Logger.Info("日志文件已轮转", 
				zap.String("from", logPath),
				zap.String("to", archivePath),
				zap.String("date", currentDate))
		}
	}

	// 重新初始化 logger
	InitWithConfig(defaultLogConfig)
}

// 解析日志级别
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// 获取当前日期作为文件名
func getCurrentDateFileName() string {
	return time.Now().Format("2006-01-02") + ".log"
}

// 手动轮转日志文件
func RotateLog() {
	checkAndRotate()
}

// 获取日志配置
func GetLogConfig() LogConfig {
	return defaultLogConfig
}

// 设置日志配置
func SetLogConfig(config LogConfig) {
	defaultLogConfig = config
}

func Info(msg string, fields ...zap.Field) {
	checkAndRotate()
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	checkAndRotate()
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

func Debug(msg string, fields ...zap.Field) {
	checkAndRotate()
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	checkAndRotate()
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	checkAndRotate()
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	}
}
