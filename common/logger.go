package common

import (
	"os"

	logging "github.com/op/go-logging"
)

var Logger = logging.MustGetLogger("example")
var Format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func InitLogger(l string) {
	var s logging.Level
	switch l {
	case "info":
		s = logging.INFO
	case "error":
		s = logging.ERROR
	case "warning":
		s = logging.WARNING
	case "debug":
		s = logging.DEBUG
	default:
		panic("Unknown log level")
	}

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, Format)
	backend1Leveled := logging.AddModuleLevel(backendFormatter)
	backend1Leveled.SetLevel(s, "")

	logging.SetBackend(backend1Leveled)
}
