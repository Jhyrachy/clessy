package tg

type ClientCommandType uint

const (
	CmdSendTextMessage ClientCommandType = 1
)

type ClientTextMessageData struct {
	ChatID  int
	Text    string
	ReplyID *int
}

type ClientCommand struct {
	Type            ClientCommandType
	TextMessageData *ClientTextMessageData
}
