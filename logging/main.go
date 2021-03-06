package logging

import (
	"os"
	"time"

	"github.com/op/go-logging"
)

const (
	size = 1024
)

var (
	Log    *Logger
	exited bool
	level  = logging.INFO // default level
	format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
)

type Logger struct {
	*logging.Logger
	backend *logging.MemoryBackend
}

func Init() *Logger {
	if Log == nil {
		logging.SetFormatter(format) // setup default formatter

		Log = &Logger{
			logging.MustGetLogger("tcolors"),
			logging.NewMemoryBackend(size),
		}

		if debugMode() {
			level = logging.DEBUG
			StartServer()
		}

		backendLvl := logging.AddModuleLevel(Log.backend)
		backendLvl.SetLevel(level, "")

		logging.SetBackend(backendLvl)
		Log.Notice("logger initialized")
	}
	return Log
}

func (log *Logger) tail() chan string {
	stream := make(chan string)

	node := log.backend.Head()
	go func() {
		for {
			stream <- node.Record.Formatted(0)
			for {
				nnode := node.Next()
				if nnode != nil {
					node = nnode
					break
				}
				if exited {
					close(stream)
					return
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return stream
}

func (log *Logger) Exit() {
	exited = true
	StopServer()
}

func debugMode() bool    { return os.Getenv("TCOLORS_DEBUG") == "1" }
func debugModeTCP() bool { return os.Getenv("TCOLORS_DEBUG_TCP") == "1" }
