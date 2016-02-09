package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hamcha/clessy/tg"
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

func (t Telegram) SetWebhook(webhook string) {
	resp, err := http.PostForm(t.apiURL("setWebhook"), url.Values{"url": {webhook}})
	if !checkerr("SetWebhook", err) {
		defer resp.Body.Close()
		var result tg.APIResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Println("[SetWebhook] Could not read reply: " + err.Error())
			return
		}
		if result.Ok {
			log.Println("Webhook successfully set!")
		} else {
			log.Printf("[SetWebhook] Error setting webhook (errcode %d): %s\n", *(result.ErrCode), *(result.Description))
			panic(errors.New("Cannot set webhook"))
		}
	}
}

func (t Telegram) SendTextMessage(data tg.ClientTextMessageData) {
	postdata := url.Values{
		"chat_id":    {strconv.Itoa(data.ChatID)},
		"text":       {data.Text},
		"parse_mode": {"HTML"},
	}
	if data.ReplyID != nil {
		postdata["reply_to_message_id"] = []string{strconv.Itoa(*(data.ReplyID))}
	}

	_, err := http.PostForm(t.apiURL("sendMessage"), postdata)
	checkerr("SendTextMessage", err)
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
