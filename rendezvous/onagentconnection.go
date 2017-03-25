package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/maxlaverse/reverse-shell/common"
	"github.com/maxlaverse/reverse-shell/message"
)

type onAgentConnection struct{}

func (h onAgentConnection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := common.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		common.Logger.Debugf("Err:%s", err)
		return
	}
	common.Logger.Debugf("New agent %s", conn.RemoteAddr().String())

	agentTable.AddAgent(conn)

	go func() {
		defer conn.Close()
		//	defer close(done)
		for {
			_, m, err := conn.ReadMessage()
			if err != nil {
				common.Logger.Debugf("Agent '%s' disconnected. Clearing sessions! Reason: %s", conn.RemoteAddr().String(), err)

				ses := sessionTable.FindSessionByAgent(conn)
				for _, masterConn := range ses {
					common.Logger.Debugf("Session '%s' was lost due to agent failure", masterConn.Id)
					a := message.ProcessTerminated{
						Id: masterConn.Id,
					}
					masterConn.State = SESSION_LOST
					if masterConn != nil {
						for _, c := range masterConn.masterConn {
							c.WriteMessage(websocket.BinaryMessage, message.ToBinary(a))
						}
					}
				}
				agentTable.RemoveAgent(conn)
				return
			}
			b := message.FromBinary(m)
			switch v := b.(type) {
			case *message.ProcessOutput:
				common.Logger.Debugf("New Agent ProcessOutput for: %s (%d)", v.Id, len(v.Id))
				masterConn := sessionTable.FindSession(v.Id)
				if masterConn == nil {
					common.Logger.Debugf("That's bad session was lost %s", v.Id)
				} else {
					for _, c := range masterConn.masterConn {
						c.WriteMessage(websocket.BinaryMessage, m)
					}
				}

			case *message.ProcessCreated:
				common.Logger.Debugf("New Agent ProcessCreated for: %s (%d), %s\n", v.Id, len(v.Id), v.WantedId)

				s := Session{
					Id:        v.Id,
					agentConn: conn,
					State:     SESSION_OPEN,
				}
				sessionTable.AddSession(&s)

				if responseTable[v.WantedId] != nil {
					responseTable[v.WantedId] <- v.Id
					close(responseTable[v.WantedId])
					responseTable[v.WantedId] = nil
				}
			case *message.ProcessTerminated:
				common.Logger.Debugf("Session ended for: %s (%d), %s\n", v.Id, len(v.Id))
				//Create just session, sent back id, wait for attachement

				masterConn := sessionTable.FindSession(v.Id)
				masterConn.State = SESSION_CLOSED
				if masterConn == nil {
					common.Logger.Debugf("That's bad session was lost %s", v.Id)
				} else {
					for _, c := range masterConn.masterConn {
						c.WriteMessage(websocket.BinaryMessage, m)
					}
				}
			case *message.SessionTable:
				common.Logger.Debugf("Restoring session table")
				//Create just session, sent back id, wait for attachement
				for _, v2 := range v.Sessions {
					s := Session{
						Id:        v2,
						agentConn: conn,
						State:     SESSION_OPEN,
					}
					common.Logger.Debugf("Adding session: %s", v2)
					sessionTable.AddSession(&s)
				}

			default:
				common.Logger.Debugf("Received Agent an unknown message type: %v", v)
			}
		}
	}()
}
