package tg

type ClientCommandType uint

const (
	CmdSendMessage ClientCommandType = 1
)

type ClientCommandMessageData struct {
	MessageText string
}

type ClientCommand struct {
	Type        ClientCommandType
	MessageData *ClientCommandMessageData
}
