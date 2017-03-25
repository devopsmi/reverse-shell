package main

import (
	"fmt"
	"net/http"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/maxlaverse/reverse-shell/common"
)

var agentTable = NewAgentTable()
var sessionTable = NewSessionTable()
var responseTable map[string]chan string = make(map[string]chan string)

func Start(port int32) {
	go http.Handle("/agent/listen", onAgentConnection{})
	go http.Handle("/agent/list", onAgentList{})
	go http.Handle("/session/list", onSessionList{})
	go http.Handle("/session/attach/", onSessionAttach{})
	go http.Handle("/session/create", onSessionCreate{})

	common.Logger.Infof("Ready for incoming connections")
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

var options struct {
	LogLevel func(string) `short:"l" long:"log-level" description:"log level" default:"info"`
	Port     int32        `short:"P" long:"port" env:"PORT" description:"Port" required:"true"`
}

func main() {
	options.LogLevel = common.InitLogger
	var parser = flags.NewParser(&options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	Start(options.Port)
}
