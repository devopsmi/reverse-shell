package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/maxlaverse/reverse-shell/common"
	"github.com/maxlaverse/reverse-shell/message"
	"github.com/maxlaverse/reverse-shell/rendezvous/api"
)

type onSessionCreate struct{}

func (h onSessionCreate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	common.Logger.Debugf("On session create")
	decoder := json.NewDecoder(r.Body)
	var m api.CreateSession
	err := decoder.Decode(&m)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	agent := agentTable.FindAgent(m.Agent)
	if agent == nil {
		common.Logger.Debugf("Agent not found %s", m.Agent)
		return
	}

	common.Logger.Debugf("Agent found %s, creating session", m.Agent)
	responseTable["generated-token"] = make(chan string)
	m2 := message.CreateProcess{
		CommandLine: m.Command,
		Id:          "generated-token",
	}
	agent.WriteMessage(websocket.BinaryMessage, message.ToBinary(m2))

	t := <-responseTable["generated-token"]
	common.Logger.Debugf("New session, answering %s", t)
	w.Write([]byte(t))
}
