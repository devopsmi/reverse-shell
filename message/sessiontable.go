package message

import "strings"

type SessionTable struct {
	Sessions []string
}

func (m SessionTable) ToBinary() []byte {
	return []byte(strings.Join(m.Sessions, ","))
}

func SessionTableFromBinary(b []byte) *SessionTable {
	if len(b) == 0 {
		return &SessionTable{}
	}
	return &SessionTable{Sessions: strings.Split(strings.TrimRight(string(b), "\x00"), ",")}
}
