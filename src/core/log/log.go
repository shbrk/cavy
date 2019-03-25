package log

import (
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var defaultConfig = Config{
	DevMode:    true, // 开发者模式
	CallerSkip: 3,
	Filepath:   "../log/", // 日志文件夹路径
	MaxSize:    10,        // 单个日志文件最大Size
	MaxBackups: 5,         // 最多同时保留的日志数量,0是全部
	MaxAge:     2,         // 日志保留的持续天数，0是永久
	Compress:   false,     // 日志文件是否用zip压缩
}

type Logger struct {
	logger         *zap.Logger
	msgChan        chan *MsgEntity
	wg             sync.WaitGroup
	fileDebugLevel *atomic.Bool
	callerSkip     int
}

func NewLogger(cfg *Config) *Logger {
	var logger = &Logger{
		msgChan: make(chan *MsgEntity, 1024*4),
	}
	if cfg == nil {
		cfg = &defaultConfig
	}
	logger.callerSkip = cfg.CallerSkip
	logger.fileDebugLevel = atomic.NewBool(cfg.DevMode)
	var fileCores = createFileCore(cfg, logger.fileDebugLevel)
	var kafkaCores = createKafkaCore(cfg)
	var cores = append(fileCores, kafkaCores...)
	if cfg.DevMode {
		var consoleCores = createConsoleCore(cfg)
		cores = append(cores, consoleCores...)
	}
	var tee = zapcore.NewTee(cores...)
	logger.logger = zap.New(tee)
	logger.wg.Add(1)
	// 开启异步消费队列
	go logger.consume()
	return logger
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.produce(zapcore.DebugLevel, msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.produce(zapcore.InfoLevel, msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.produce(zapcore.WarnLevel, msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.produce(zapcore.ErrorLevel, msg, fields...)
}

//调用后直接panic， 直接在当前线程执行
func (l *Logger) Panic(msg string, fields ...Field) {
	l.logger.Panic(msg, fields...)
}

//调用后直接退出， 直接在当前线程执行
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
}

// 动态开启debug文件日志输出
func (l *Logger) SetFileDebugOutput(enabled bool) {
	l.fileDebugLevel.Store(enabled)
}
