package main

import (
	"math/rand"

	"github.com/hamcha/clessy/tg"
)

var actions []string = []string{
	"Puppami", "Degustami", "Lucidami", "Manipolami", "Disidratami", "Irritami", "Martorizzami",
	"Lustrami", "Osannami", "Sorseggiami", "Assaporami", "Apostrofami", "Spremimi", "Dimenami",
	"Agitami", "Stimolami", "Suonami", "Strimpellami", "Stuzzicami", "Spintonami", "Sguinzagliami",
	"Modellami", "Sgrullami", "Cavalcami", "Perquotimi", "Misurami", "Sventolami", "Induriscimi",
	"Accordami", "Debuggami",
}

var objects []string = []string{
	"il birillo", "il bastone", "l'ombrello", "il malloppo", "il manico", "il manganello",
	"il ferro", "la mazza", "l'archibugio", "il timone", "l'arpione", "il flauto", "la reliquia",
	"il fiorino", "lo scettro", "il campanile", "la proboscide", "il pino", "il maritozzo", "il perno",
	"il tubo da 100", "la verga", "l'idrante", "il pendolo", "la torre di Pisa", "la lancia",
	"il cilindro", "il lampione", "il joystick", "il Wiimote", "il PSMove", "l'albero maestro",
	"il trenino",
}

func metafora(broker *tg.Broker, update tg.APIMessage) {
	if isCommand(update, "metafora") {
		n := rand.Intn(len(actions))
		m := rand.Intn(len(objects))
		broker.SendTextMessage(update.Chat, actions[n]+" "+objects[m], nil)
		return
	}
}
