package tg

type ClientCommandType uint

const (
	CmdSendTextMessage ClientCommandType = 1
)

type ClientTextMessageData struct {
	Text string
}

type ClientCommand struct {
	Type            ClientCommandType
	TextMessageData *ClientTextMessageData
}
