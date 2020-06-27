package log

import (
	"fmt"
	base "log"
	"log/syslog"
	"os"
	"strings"
)

var mapPriority = map[string]syslog.Priority{
	"panic":    syslog.LOG_EMERG,
	"alert":    syslog.LOG_ALERT,
	"critical": syslog.LOG_CRIT,
	"error":    syslog.LOG_ERR,
	"warning":  syslog.LOG_WARNING,
	"notice":   syslog.LOG_NOTICE,
	"info":     syslog.LOG_INFO,
	"debug":    syslog.LOG_DEBUG,
}

type logDestination int

const (
	Stdout logDestination = iota
	Stderr
	Syslog
)

var mapDestination = map[string]logDestination{
	"stdout": Stdout,
	"stderr": Stderr,
	"syslog": Syslog,
}

type log struct {
	priority syslog.Priority
	syslog   *syslog.Writer
	stdout   *base.Logger
}

var logger log

func Init(destinationStr string, logTag string, priorityStr string) (err error) {
	priority, ok := mapPriority[strings.ToLower(priorityStr)]
	if !ok {
		priority = syslog.LOG_DEBUG
	}

	destination, ok := mapDestination[strings.ToLower(destinationStr)]
	if !ok {
		destination = Stderr
	}

	logger.priority = priority

	switch destination {
	case Syslog:
		if logger.syslog, err = syslog.New(priority|syslog.LOG_LOCAL0, strings.ToLower(logTag)); err != nil {
			return
		}
	case Stderr:
		logger.stdout = base.New(os.Stderr, "", base.LstdFlags)
	default:
		logger.stdout = base.New(os.Stdout, "", base.LstdFlags)
	}

	return
}

func Close() {
	if logger.syslog != nil {
		_ = logger.syslog.Close()
	}
}

func Debug(args ...interface{}) {
	if logger.priority >= syslog.LOG_DEBUG {
		if logger.syslog != nil {
			_ = logger.syslog.Debug(fmt.Sprint(args...))
		} else if logger.stdout != nil {
			logger.stdout.Print(args...)
		}
	}
}

func Info(args ...interface{}) {
	if logger.priority >= syslog.LOG_INFO {
		if logger.syslog != nil {
			_ = logger.syslog.Info(fmt.Sprint(args...))
		} else if logger.stdout != nil {
			logger.stdout.Print(args...)
		}
	}
}

func Notice(args ...interface{}) {
	if logger.priority >= syslog.LOG_NOTICE {
		if logger.syslog != nil {
			_ = logger.syslog.Notice(fmt.Sprint(args...))
		} else if logger.stdout != nil {
			logger.stdout.Print(args...)
		}
	}
}

func Warning(args ...interface{}) {
	if logger.priority >= syslog.LOG_WARNING {
		if logger.syslog != nil {
			_ = logger.syslog.Warning(fmt.Sprint(args...))
		} else if logger.stdout != nil {
			logger.stdout.Print(args...)
		}
	}
}

func Error(args ...interface{}) {
	if logger.priority >= syslog.LOG_ERR {
		if logger.syslog != nil {
			_ = logger.syslog.Err(fmt.Sprint(args...))
		} else if logger.stdout != nil {
			logger.stdout.Print(args...)
		}
	}
}

func Critical(args ...interface{}) {
	if logger.priority >= syslog.LOG_CRIT {
		if logger.syslog != nil {
			_ = logger.syslog.Crit(fmt.Sprint(args...))
		} else if logger.stdout != nil {
			logger.stdout.Print(args...)
		}
	}
}

func Alert(args ...interface{}) {
	if logger.priority >= syslog.LOG_ALERT {
		if logger.syslog != nil {
			_ = logger.syslog.Alert(fmt.Sprint(args...))
		} else if logger.stdout != nil {
			logger.stdout.Print(args...)
		}
	}
}

func Emerg(args ...interface{}) {
	if logger.priority >= syslog.LOG_EMERG {
		if logger.syslog != nil {
			_ = logger.syslog.Emerg(fmt.Sprint(args...))
		} else if logger.stdout != nil {
			logger.stdout.Print(args...)
		}
	}
}
