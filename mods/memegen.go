package main

import (
	"strings"

	"github.com/hamcha/clessy/tg"
)

func memegen(broker *tg.Broker, update tg.APIMessage) {
	if update.Caption != nil {
		if strings.HasPrefix(*(update.Caption), "/meme") {
		}
	}
}
