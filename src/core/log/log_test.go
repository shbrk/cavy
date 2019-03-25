package log

import (
	"fmt"
	"testing"
)

func TestDevInitLogger(t *testing.T) {
	var cfg = Config{
		DevMode:    true,
		Filepath:   "d:/log/dev/", // 日志文件夹路径
		MaxSize:    1,             // 单个日志文件最大Size
		MaxBackups: 2,             // 最多同时保留的日志数量,0是全部
		MaxAge:     1,             // 日志保留的持续天数，0是永久
		Compress:   false,         // 日志文件是否用zip压缩
	}
	logger := NewLogger(&cfg)
	fmt.Println("Log Start")
	logger.Debug("HelloWorld!", Int("A", 1))
	logger.Info("HelloWorld!", Int("A", 2))
	logger.Warn("HelloWorld!", Int("A", 3))
	logger.Error("HelloWorld!")
	fmt.Println("Log Done")
	logger.Debug("should appear in files!", Int("A", 1))
	logger.Debug("should appear in files!", Int("A", 1))
	logger.Debug("should appear in files!", Int("A", 1))
	logger.Sync()
	fmt.Println("Sync Done")
}
