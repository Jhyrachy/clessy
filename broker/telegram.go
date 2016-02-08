package main

import (
	"encoding/json"
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
		var result tg.APIResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Println("Could not read reply: " + err.Error())
			return
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
