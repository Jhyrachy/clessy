package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const APIEndpoint = "https://api.telegram.org/"

type Telegram struct {
	Token string
}

func mkAPI(token string) *Telegram {
	tg := new(Telegram)
	tg.Token = token
	return tg
}

func (t Telegram) setWebhook(webhook string) {
	resp, err := http.PostForm(t.apiURL("setWebhook"), url.Values{"url": {webhook}})
	if !checkerr("setWebhook", err) {
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("setWebhook result: %s\n", string(data))
		}
	}
}

func (t Telegram) apiURL(method string) string {
	return APIEndpoint + "bot" + t.Token + "/" + method
}

func checkerr(method string, err error) bool {
	if err != nil {
		log.Printf("Received error with call to %s: %s\n", method, err.Error())
		return true
	}
	return false
}
