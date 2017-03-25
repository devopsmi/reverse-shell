package message

import "strings"

type ProcessTerminated struct {
	Id     string
	Reason []byte
}

func (m ProcessTerminated) ToBinary() []byte {
	session := make([]byte, SessionIdentifierMaxLen)
	copy(session, m.Id[:])

	return append(session, m.Reason...)
}

func ProcessTerminatedFromBinary(b []byte) *ProcessTerminated {
	return &ProcessTerminated{
		Id:     strings.TrimRight(string(b[0:SessionIdentifierMaxLen-1]), "\x00"),
		Reason: b[SessionIdentifierMaxLen:],
	}
}
