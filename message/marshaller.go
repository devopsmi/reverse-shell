package message

import (
	"os"
)

const messageTypeLength = 1

type Serializable interface {
	ToBinary() []byte
}

func FromBinary(b []byte) interface{} {
	if len(b) < messageTypeLength {
		return b
	}

	messageType := b[0]
	switch messageType {
	case 1:
		return CreateProcessFromBinary(b[messageTypeLength:])
	case 2:
		return ExecuteCommandFromBinary(b[messageTypeLength:])
	case 3:
		return ProcessOutputFromBinary(b[messageTypeLength:])
	case 4:
		return ProcessCreatedFromBinary(b[messageTypeLength:])
	case 5:
		return ProcessTerminatedFromBinary(b[messageTypeLength:])
	case 6:
		return SessionTableFromBinary(b[messageTypeLength:])
	default:
		return b
	}
}

func ToBinary(b Serializable) []byte {
	var flag []byte
	switch b.(type) {
	case ExecuteCommand:
		flag = []byte{2}
	case CreateProcess:
		flag = []byte{1}
	case ProcessOutput:
		flag = []byte{3}
	case ProcessCreated:
		flag = []byte{4}
	case ProcessTerminated:
		flag = []byte{5}
	case SessionTable:
		flag = []byte{6}
	default:
		os.Exit(1)
	}

	s := append(flag, b.ToBinary()...)
	return s
}
