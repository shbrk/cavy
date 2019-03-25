package log

import (
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

const QUEUE_SIZE = 1024 * 64

type Logger struct {
	logger         *zap.Logger
	msgChan        chan *MsgEntity
	wg             sync.WaitGroup
	fileDebugLevel *atomic.Bool
	callerSkip     int
}

func NewLogger(cfg *Config) *Logger {
	var logger = &Logger{
		msgChan: make(chan *MsgEntity, QUEUE_SIZE),
	}
	if cfg == nil {
		cfg = &defaultConfig
	}
	logger.callerSkip = cfg.CallerSkip
	logger.fileDebugLevel = atomic.NewBool(cfg.DevMode)
	var cores = createFileCore(cfg, logger.fileDebugLevel)
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
