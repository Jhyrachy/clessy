package main

import (
	"flag"
	"strings"

	"github.com/hamcha/clessy/tg"
)

func initmods() {
	initviaggi()
}

func dispatch(broker *tg.Broker, update tg.APIMessage) {
	metafora(broker, update)
	viaggi(broker, update)
}

func isCommand(update tg.APIMessage, cmdname string) bool {
	if update.Text == nil {
		return false
	}

	text := *(update.Text)
	return strings.HasPrefix(text, "/"+cmdname+"@"+*botname) || (strings.HasPrefix(text, "/"+cmdname) && !strings.Contains(text, "@"))
}

var botname *string

func main() {
	brokerAddr := flag.String("broker", "localhost:7314", "Broker address:port")
	botname = flag.String("botname", "maudbot", "Bot name for /targetet@commands")
	flag.Parse()

	initmods()

	err := tg.CreateBrokerClient(*brokerAddr, dispatch)
	if err != nil {
		panic(err)
	}
}
