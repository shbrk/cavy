package log

//包默认日志对象
var logger *Logger

func InitDefaultLogger(cfg *Config) {
	if logger == nil {
		logger = NewLogger(cfg)
	}
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
