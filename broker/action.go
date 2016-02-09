package main

import (
	"github.com/hamcha/clessy/tg"
)

func executeClientCommand(action tg.ClientCommand) {
	switch action.Type {
	case tg.CmdSendTextMessage:
		data := *(action.TextMessageData)
		api.SendTextMessage(data)
	}
}
