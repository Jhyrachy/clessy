package main

import (
	"io/ioutil"
	"net/http"
)

func webhook(rw http.ResponseWriter, req *http.Request) {
	// Read entire request and broadcast to everyone
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	defer req.Body.Close()

	broadcast(string(data))
}
