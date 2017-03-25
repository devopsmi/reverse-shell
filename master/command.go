package main

import (
	"fmt"

	"github.com/maxlaverse/reverse-shell/common"
)

type CreateCommand struct {
	Url string `short:"U" long:"url" env:"URL" description:"Url of the rendez-vous point." required:"true"`
}

func (x *CreateCommand) Execute(args []string) error {
	sessionId := CreateSession(x.Url, args[0])
	AttachSession(x.Url, sessionId)
	return nil
}

type AttachCommand struct {
	Url string `short:"U" long:"url" env:"URL" description:"Url of the rendez-vous point." required:"true"`
}

func (x *AttachCommand) Execute(args []string) error {
	AttachSession(x.Url, args[0])
	return nil
}

type ListSessionCommand struct {
	Url string `short:"U" long:"url" env:"URL" description:"Url of the rendez-vous point." required:"true"`
}

func (x *ListSessionCommand) Execute(args []string) error {
	fmt.Printf("List of sessions:\n")
	l, err := ListSessions(x.Url)
	for _, v := range l {
		fmt.Printf(" * %s => agent: %s, masters: %s, state: %s\n", v.Name, v.Agent, v.Masters, v.State)
	}
	return err
}

type ListAgentCommand struct {
	Url string `short:"U" long:"url" env:"URL" description:"Url of the rendez-vous point." required:"true"`
}

func (x *ListAgentCommand) Execute(args []string) error {
	fmt.Printf("List of agents:\n")
	l, err := ListAgents(x.Url)
	for _, v := range l {
		fmt.Printf(" * %s\n", v.Name)
	}
	return err
}

type ListenCommand struct {
	Port int `short:"P" long:"port" env:"PORT" description:"Port to listen to." required:"true"`
}

func (x *ListenCommand) Execute(args []string) error {
	return Listen(x.Port)
}

var options struct {
	LogLevel           func(string)       `short:"l" long:"log-level" description:"log level" default:"info"`
	createOptions      CreateCommand      `command:"create"`
	attachOptions      AttachCommand      `command:"attach"`
	listSessionOptions ListSessionCommand `command:"list-session"`
	listAgentOptions   ListAgentCommand   `command:"list-agent"`
	listenOptions      ListenCommand      `command:"listen"`
}

func init() {
	var createCommand CreateCommand
	var attachCommand AttachCommand
	var listSessionCommand ListSessionCommand
	var listAgentCommand ListAgentCommand
	var listenCommand ListenCommand

	options.LogLevel = common.InitLogger
	parser.AddCommand("attach",
		"attach to an existing session",
		"",
		&attachCommand)
	parser.AddCommand("create",
		"create a new session on a given agent",
		"",
		&createCommand)
	parser.AddCommand("list-agents",
		"list all the agents available on a rendez-vous",
		"",
		&listAgentCommand)
	parser.AddCommand("list-sessions",
		"list all the sessions available on a rendez-vous",
		"",
		&listSessionCommand)
	parser.AddCommand("listen",
		"listen for agents to connect using websocket",
		"",
		&listenCommand)
}
