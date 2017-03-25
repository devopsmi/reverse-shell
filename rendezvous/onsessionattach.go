package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/maxlaverse/reverse-shell/common"
	"github.com/maxlaverse/reverse-shell/message"
)

type onSessionAttach struct{}

func (h onSessionAttach) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestedSession := r.URL.Path[16:]
	common.Logger.Debugf("New attachement request for '%s'", requestedSession)

	session := sessionTable.FindSession(requestedSession)
	if session == nil {
		common.Logger.Debugf("Session not found")
		w.Write([]byte("Session not found"))
		return
	}

	conn, err := common.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	session.masterConn = append(session.masterConn, conn)

	go func() {
		defer conn.Close()

		for {
			_, m, err := conn.ReadMessage()
			if err != nil {
				common.Logger.Debugf("ReadMessage error on the masterChannel: %s", err)
				return
			}
			b := message.FromBinary(m)
			switch v := b.(type) {
			case *message.CreateProcess:
				session.agentConn.WriteMessage(websocket.BinaryMessage, m)
			case *message.ExecuteCommand:
				session.agentConn.WriteMessage(websocket.BinaryMessage, m)
			default:
				common.Logger.Debugf("Received Master an unknown message type: %v", v)
			}
		}
	}()
}
