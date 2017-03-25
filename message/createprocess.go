package message

import "strings"

type CreateProcess struct {
	CommandLine string
	Id          string
}

func (m CreateProcess) ToBinary() []byte {
	session := make([]byte, SessionIdentifierMaxLen)
	copy(session, m.Id[:])

	return append(session, m.CommandLine...)
}

func CreateProcessFromBinary(b []byte) *CreateProcess {
	return &CreateProcess{
		Id:          strings.TrimRight(string(b[0:SessionIdentifierMaxLen-1]), "\x00"),
		CommandLine: strings.TrimRight(string(b[SessionIdentifierMaxLen:]), "\x00"),
	}
}
