package logging

import (
	"time"
)

type Timer struct {
	msg    string
	start  time.Time
	logger *Logger
}

func (log *Logger) NewTimer(msg string) *Timer {
	return &Timer{msg, time.Now(), log}
}

func (t *Timer) End() {
	ms := time.Since(t.start).Seconds() * 1000
	t.logger.Debugf("%s [%3.3fms]", t.msg, ms)
}
