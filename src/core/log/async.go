package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
)

type MsgEntity struct {
	level  zapcore.Level
	msg    string
	caller string
	fields []Field
}

//消费函数
func (l *Logger) consume() {
	for {
		msg, ok := <-l.msgChan
		if !ok {
			l.wg.Done()
			return
		}
		if ce := l.logger.Check(msg.level, msg.msg); ce != nil {
			// 将caller加入到field 不破坏zap结构
			fields := append(msg.fields, zap.String("caller", msg.caller))
			ce.Write(fields...)
		}
	}
}

//生产函数
func (l *Logger) produce(level zapcore.Level, msg string, fields ...Field) {
	caller := zapcore.NewEntryCaller(runtime.Caller(l.callerSkip))
	l.msgChan <- &MsgEntity{level, msg, caller.TrimmedPath(), fields}
}

// 阻塞同步，只允许程序退出之前调用
func (l *Logger) Sync() error {
	close(l.msgChan)
	l.wg.Wait()
	return l.logger.Sync()
}
