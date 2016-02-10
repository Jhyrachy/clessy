package main

import (
	"log"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/hamcha/clessy/tg"
)

var users map[string]string

func getNick(apiuser tg.APIUser) {
	if _, ok := users[apiuser.Username]; ok && strings.HasPrefix(users[apiuser.Username], apiuser.FirstName) {
		// It's updated, don't bother
		return
	}

	users[apiuser.Username] = apiuser.FirstName
	if apiuser.LastName != "" {
		users[apiuser.Username] += apiuser.LastName
	}

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("usernames"))
		return b.Put([]byte(apiuser.Username), []byte(users[apiuser.Username]))
	})
	if err != nil {
		log.Printf("[getNick] Could not update %s name: %s\n", apiuser.Username, err.Error())
	}
}

func loadUsers() {
	users = make(map[string]string)
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("usernames"))
		if err != nil {
			return err
		}
		b.ForEach(func(user, name []byte) error {
			users[string(user)] = string(name)
			return nil
		})
		return nil
	})
	assert(err)
}
