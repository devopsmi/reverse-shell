package message

import "strings"

type ProcessCreated struct {
	Id       string
	WantedId string
}

func (m ProcessCreated) ToBinary() []byte {
	session := make([]byte, SessionIdentifierMaxLen*2)
	copy(session, m.Id[:])
	copy(session[SessionIdentifierMaxLen:], m.WantedId[:])
	return session
}

func ProcessCreatedFromBinary(b []byte) *ProcessCreated {
	return &ProcessCreated{
		Id:       strings.TrimRight(string(b[0:SessionIdentifierMaxLen-1]), "\x00"),
		WantedId: strings.TrimRight(string(b[SessionIdentifierMaxLen:SessionIdentifierMaxLen*2-1]), "\x00"),
	}
}
