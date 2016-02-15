package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/hamcha/clessy/tg"
)

const viaggiurl = "http://free.rome2rio.com/api/1.2/json/Search?key=X5JMLHNc&languageCode=IT&currencyCode=EUR"

var reg *regexp.Regexp

func initviaggi() {
	reg = regexp.MustCompile("([^-]+) -> (.+)")
}

func viaggi(broker *tg.Broker, update tg.APIMessage) {
	if isCommand(update, "viaggi") {
		usage := func() {
			broker.SendTextMessage(update.Chat, "Formato: /viaggi <i>&lt;PARTENZA&gt;</i> -> <i>&lt;DESTINAZIONE&gt;</i>", &update.MessageID)
		}
		oops := func(err error) {
			log.Println("[viaggi] GET error:" + err.Error())
			broker.SendTextMessage(update.Chat, "<b>ERRORE!</b> @hamcha controlla la console!", &update.MessageID)
		}

		parts := strings.SplitN(*(update.Text), " ", 2)
		if len(parts) < 2 {
			usage()
			return
		}
		text := parts[1]
		msgs := reg.FindStringSubmatch(text)
		if len(msgs) <= 2 {
			usage()
			return
		}

		src := url.QueryEscape(msgs[1])
		dst := url.QueryEscape(msgs[2])
		url := viaggiurl + "&oName=" + src + "&dName=" + dst
		resp, err := http.Get(url)
		if err != nil {
			oops(err)
			return
		}
		defer resp.Body.Close()

		var outjson Romejson
		err = json.NewDecoder(resp.Body).Decode(&outjson)
		if err != nil {
			oops(err)
			return
		}

		var moreeco Romeroute
		var lesstim Romeroute
		if len(outjson.Routes) < 1 {
			// Should never happen
			log.Println("[viaggi] No routes found (??)")
			broker.SendTextMessage(update.Chat, "<b>ERRORE!</b> @hamcha controlla la console!", &update.MessageID)
			return
		}

		// Calculate cheapest and fastest
		moreeco = outjson.Routes[0]
		lesstim = outjson.Routes[0]
		for _, v := range outjson.Routes {
			if v.IndicativePrice.Price < moreeco.IndicativePrice.Price {
				moreeco = v
			}
			if v.Duration < lesstim.Duration {
				lesstim = v
			}
		}

		broker.SendTextMessage(update.Chat,
			"Viaggio da <b>"+outjson.Places[0].Name+
				"</b> a <b>"+outjson.Places[1].Name+"</b>"+
				"\n\n"+
				"Piu economico: <b>"+moreeco.Name+"</b> ("+parseData(moreeco)+")"+
				"\n"+
				"Piu veloce: <b>"+lesstim.Name+"</b> ("+parseData(lesstim)+")"+
				"\n\n"+
				"<a href=\"http://www.rome2rio.com/it/s/"+src+"/"+dst+"\">Maggiori informazioni</a>",
			&update.MessageID)
	}
}

func parseData(route Romeroute) string {
	// Get time
	minutes := int(route.Duration)
	hours := minutes / 60
	minutes -= hours * 60
	days := hours / 24
	hours -= days * 24
	timestamp := ""
	if days > 0 {
		timestamp += strconv.Itoa(days) + "d "
	}
	if hours > 0 {
		timestamp += strconv.Itoa(hours) + "h "
	}
	if minutes > 0 {
		timestamp += strconv.Itoa(minutes) + "m"
	}

	return strconv.Itoa(int(route.IndicativePrice.Price)) + " " + route.IndicativePrice.Currency + " - " + strconv.Itoa(int(route.Distance)) + " Km - " + timestamp
}

type Romeplace struct {
	Name string
}

type Romeprice struct {
	Price    float64
	Currency string
}

type Romeroute struct {
	Name            string
	Distance        float64
	Duration        float64
	IndicativePrice Romeprice
}

type Romejson struct {
	Places []Romeplace
	Routes []Romeroute
}
