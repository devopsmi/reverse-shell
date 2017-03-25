package message

import "strings"

type ExecuteCommand struct {
	Id      string
	Command []byte
}

func (m ExecuteCommand) ToBinary() []byte {
	session := make([]byte, SessionIdentifierMaxLen)
	copy(session, m.Id[:])

	return append(session, m.Command...)
}

func ExecuteCommandFromBinary(b []byte) *ExecuteCommand {
	return &ExecuteCommand{
		Id:      strings.TrimRight(string(b[0:SessionIdentifierMaxLen-1]), "\x00"),
		Command: b[SessionIdentifierMaxLen:],
	}
}
