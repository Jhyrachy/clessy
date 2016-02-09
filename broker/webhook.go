package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func webhook(rw http.ResponseWriter, req *http.Request) {
	log.Println("Received request! Details follow:")
	defer req.Body.Close()
	/*
		var update tg.APIUpdate

		err := json.NewDecoder(req.Body).Decode(&update)
		if err != nil {
			log.Println("ERR: Not JSON!")
			return
		}

		jenc, _ := json.Marshal(update)
		log.Println(jenc)
	*/
	io.Copy(os.Stdout, req.Body)
}