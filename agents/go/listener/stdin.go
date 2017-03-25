package listener

import (
	"io"
	"os"
	"os/signal"

	"github.com/maxlaverse/reverse-shell/agents/go/handler"
	"github.com/maxlaverse/reverse-shell/common"
)

type stdin struct {
	processOutput     chan *handler.ProcessOutput
	processTerminated chan *handler.ProcessTerminated
	input             chan []byte
	interruptChannel  chan os.Signal
	handler           *handler.Handler
}

func NewStdin() *stdin {
	processOutput := make(chan *handler.ProcessOutput)
	processTerminated := make(chan *handler.ProcessTerminated)

	return &stdin{
		processOutput:     processOutput,
		processTerminated: processTerminated,
		input:             make(chan []byte),
		interruptChannel:  make(chan os.Signal, 1),
		handler:           handler.New(processOutput, processTerminated),
	}
}

func (l *stdin) Start() error {
	return nil
}

func (l *stdin) Listen() {
	go l.pipeFromProcessOutput()
	go l.pipeToProcessInput()

	processId := l.handler.CreateProcess("bash --norc")

	signal.Notify(l.interruptChannel, os.Interrupt)

	for {
		select {
		case <-l.interruptChannel:
			l.handler.ExecuteCommand(processId, []byte{'\u0003'})
		case msg := <-l.input:
			l.handler.ExecuteCommand(processId, msg)
		}
	}
}

func (l *stdin) pipeToProcessInput() {
	for {
		select {
		default:
			var msg = make([]byte, 1024)
			size, err := os.Stdin.Read(msg)
			common.Logger.Debugf("%s", msg)
			if err == io.EOF {
				return
			} else if err != nil {
				panic(err)
			} else {
				l.input <- msg[0:size]
			}
		}
	}
}

func (l *stdin) pipeFromProcessOutput() {
	for {
		select {
		case a := <-l.processOutput:
			common.Logger.Debugf("%s", a.Result)
		case <-l.processTerminated:
			common.Logger.Debugf("Received a processTerminated signal!")
			os.Exit(0)
		}
	}
}
