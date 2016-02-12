package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

type User struct {
	Username string `json:"username"`
}

type Message struct {
	Event   string `json:"event"`
	ReplyID *int   `json:"reply_id,omitempty"`
	FwdUser *User  `json:"fwd_from"`
	Text    string `json:"text"`
	From    User   `json:"from"`
	Date    int    `json:"date"`
}

type Stats struct {
	Total     uint64
	ByUser    map[string]uint64
	ByHour    [24]uint64
	ByWeekday [7]uint64
	ByDate    map[string]uint64
	Replies   uint64
	Forward   uint64
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	logFile := flag.String("logfile", "tl.log", "Telegram CLI dump")
	boltdbFile := flag.String("boltdb", "stats.db", "BoltDB database file")
	flag.Parse()

	file, err := os.Open(*logFile)
	assert(err)

	var data Stats
	data.ByUser = make(map[string]uint64)
	data.ByDate = make(map[string]uint64)

	scanner := bufio.NewScanner(file)
	lines := 0
	log.Printf("Started processing %s...\n\n", *logFile)
	for scanner.Scan() {
		var msg Message
		err := json.Unmarshal([]byte(scanner.Text()), &msg)
		assert(err)

		processMessage(msg, &data)

		lines++
		if lines%10000 == 0 {
			log.Printf("Processed %d lines...\n", lines)
		}
	}

	log.Printf("\nFinished with success, processed %d lines.\n", lines)

	log.Printf("Opening database %s for writing...\n", *boltdbFile)
	db, err := bolt.Open(*boltdbFile, 0600, nil)
	err = update(db, data)
	assert(err)

	log.Printf("All done! Bye!\n")
}

func processMessage(msg Message, data *Stats) {
	data.Total++

	if msg.ReplyID != nil {
		data.Replies++
	}

	if msg.FwdUser != nil {
		data.Forward++
	}

	date := time.Unix(int64(msg.Date), 0)

	data.ByHour[date.Hour()]++
	data.ByWeekday[date.Weekday()]++

	val, exists := data.ByUser[msg.From.Username]
	if !exists {
		val = 0
	}
	data.ByUser[msg.From.Username] = val + 1

	datekey := date.Format("2006-1-2")
	val, exists = data.ByDate[datekey]
	if !exists {
		val = 0
	}
	data.ByDate[datekey] = val + 1
}

func MakeUint(bval []byte, bucketName string, key string) uint64 {
	if bval != nil {
		intval, bts := binary.Uvarint(bval)
		if bts > 0 {
			return intval
		} else {
			log.Printf("[%s] Value of key \"%s\" is NaN: %v\r\n", bucketName, key, bval)
			return 0
		}
	} else {
		log.Printf("[%s] Key \"%s\" does not exist, set to 0\n", bucketName, key)
		return 0
	}
}

func PutUint(value uint64) []byte {
	bytes := make([]byte, 10)
	n := binary.PutUvarint(bytes, value)
	return bytes[:n]
}

func update(db *bolt.DB, data Stats) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("global"))
		if err != nil {
			return err
		}

		// Update total
		total := MakeUint(b.Get([]byte("count")), "global", "count")
		total += data.Total
		err = b.Put([]byte("count"), PutUint(total))
		if err != nil {
			return err
		}

		// Update replies
		replies := MakeUint(b.Get([]byte("replies")), "global", "replies")
		replies += data.Replies
		err = b.Put([]byte("replies"), PutUint(total))
		if err != nil {
			return err
		}

		// Update forward
		forward := MakeUint(b.Get([]byte("forward")), "global", "forward")
		forward += data.Forward
		err = b.Put([]byte("forward"), PutUint(total))
		if err != nil {
			return err
		}

		// Update hour counters
		b, err = tx.CreateBucketIfNotExists([]byte("hour"))
		if err != nil {
			return err
		}

		for i := 0; i < 24; i++ {
			curhour := MakeUint(b.Get([]byte{byte(i)}), "hour", strconv.Itoa(i))
			curhour += data.ByHour[i]
			err = b.Put([]byte{byte(i)}, PutUint(curhour))
			if err != nil {
				return err
			}
		}

		// Update weekday counters
		b, err = tx.CreateBucketIfNotExists([]byte("weekday"))
		if err != nil {
			return err
		}

		for i := 0; i < 7; i++ {
			curwday := MakeUint(b.Get([]byte{byte(i)}), "weekday", strconv.Itoa(i))
			curwday += data.ByWeekday[i]
			err = b.Put([]byte{byte(i)}, PutUint(curwday))
			if err != nil {
				return err
			}
		}

		// Update date counters
		b, err = tx.CreateBucketIfNotExists([]byte("date"))
		if err != nil {
			return err
		}

		for day, count := range data.ByDate {
			count += MakeUint(b.Get([]byte(day)), "date", day)
			err = b.Put([]byte(day), PutUint(count))
			if err != nil {
				return err
			}
		}

		// Update user counters
		b, err = tx.CreateBucketIfNotExists([]byte("users-count"))
		if err != nil {
			return err
		}

		for user, count := range data.ByUser {
			count += MakeUint(b.Get([]byte(user)), "users-count", user)
			err = b.Put([]byte(user), PutUint(count))
			if err != nil {
				return err
			}
		}

		return nil
	})
}
