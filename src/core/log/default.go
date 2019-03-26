package log

import "go.uber.org/zap"

//包默认日志对象
var logger *Logger

var defaultConfig = Config{
	DevMode:    true, // 开发者模式
	CallerSkip: 3,
	Filepath:   "../log/", // 日志文件夹路径
	MaxSize:    10,        // 单个日志文件最大Size
	MaxBackups: 5,         // 最多同时保留的日志数量,0是全部
	MaxAge:     2,         // 日志保留的持续天数，0是永久
	Compress:   false,     // 日志文件是否用zip压缩
}

func InitDefaultLogger(cfg *Config) {
	if logger == nil {
		logger = NewLogger(cfg)
	}
}

func GetRaw() *zap.Logger {
	return logger.logger
}

func Debug(msg string, fields ...Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	logger.Error(msg, fields...)
}

//调用后直接panic， 直接在当前线程执行
func Panic(msg string, fields ...Field) {
	logger.Panic(msg, fields...)
}

//调用后直接panic， 直接在当前线程执行
func Fatal(msg string, fields ...Field) {
	logger.Fatal(msg, fields...)
}

// 动态开启debug文件日志输出
func SetFileDebugOutput(enabled bool) {
	logger.SetFileDebugOutput(enabled)
}

func Sync() error {
	return logger.Sync()
}
