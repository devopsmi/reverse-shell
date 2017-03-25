package message

import "strings"

type ProcessOutput struct {
	Id   string
	Data []byte
}

func (m ProcessOutput) ToBinary() []byte {
	session := make([]byte, SessionIdentifierMaxLen)
	copy(session, m.Id[:])

	return append(session, m.Data...)
}

func ProcessOutputFromBinary(b []byte) *ProcessOutput {
	return &ProcessOutput{
		Id:   strings.TrimRight(string(b[0:SessionIdentifierMaxLen-1]), "\x00"),
		Data: b[SessionIdentifierMaxLen:],
	}
}
