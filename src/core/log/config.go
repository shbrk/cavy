package log

import (
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
)

type Config struct {
	DevMode    bool   // 是不是开发模式，开发模式会控制台打印,会开启debug level日志
	Filepath   string // 日志文件夹路径
	CallerSkip int    // 调用函数层级
	MaxSize    int    // 单个日志文件最大Size
	MaxBackups int    // 最多同时保留的日志数量,0是全部
	MaxAge     int    // 日志保留的持续天数，0是永久
	Compress   bool   // 日志文件是否用zip压缩
}

func createFileCore(cfg *Config, fileDebugFlag *atomic.Bool) []zapcore.Core {
	var levels = []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.DPanicLevel,
		zapcore.PanicLevel,
		zapcore.FatalLevel,
	}

	var cores = make([]zapcore.Core, 0, 7)
	// 为每一个level创建一种日志文件
	for _, level := range levels {
		hook := lumberjack.Logger{
			Filename:   path.Join(cfg.Filepath, level.String()+".log"),
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		writer := zapcore.AddSync(&hook)

		var encoderConfig = zap.NewProductionEncoderConfig()
		if cfg.DevMode {
			encoderConfig = zap.NewDevelopmentEncoderConfig()
		}
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		var curLevel = level
		var levelFunc = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl == curLevel })
		if cfg.DevMode == false && level == zapcore.DebugLevel {
			levelFunc = zap.LevelEnablerFunc(
				func(lvl zapcore.Level) bool { return lvl == curLevel && fileDebugFlag.Load() == true })
		}

		core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), writer, levelFunc)
		cores = append(cores, core)
	}
	return cores
}

func createConsoleCore(cfg *Config) []zapcore.Core {
	var consoleDebugging = zapcore.Lock(os.Stdout)
	var encoderConfig = zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	var encoder = zapcore.NewConsoleEncoder(encoderConfig)
	return []zapcore.Core{
		zapcore.NewCore(encoder, consoleDebugging, zapcore.DebugLevel),
	}
}
