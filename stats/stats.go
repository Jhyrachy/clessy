package main

import (
	"encoding/binary"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hamcha/clessy/tg"
)

const (
	MessageTypeText     int = 0
	MessageTypeAudio    int = 1
	MessageTypePhoto    int = 2
	MessageTypeSticker  int = 3
	MessageTypeVideo    int = 4
	MessageTypeVoice    int = 5
	MessageTypeContact  int = 6
	MessageTypeLocation int = 7
	MessageTypeDocument int = 8
	MessageTypeMax      int = 9
)

type Stats struct {
	ByUserCount    map[string]uint64
	ByUserAvgLen   map[string]uint64
	ByUserAvgCount map[string]uint64
	ByWeekday      [7]uint64
	ByHour         [24]uint64
	ByType         [MessageTypeMax]uint64
	TodayDate      time.Time
	Today          uint64
	TotalCount     uint64
	TotalTxtCount  uint64
	TotalAvgLength uint64
	Replies        uint64
	Forward        uint64
}

var stats Stats

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

func loadStats() {
	// Load today
	stats.TodayDate = time.Now()

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("global"))
		if err != nil {
			return err
		}

		// Load total messages counter
		bval := b.Get([]byte("count"))
		stats.TotalCount = MakeUint(bval, "global", "count")

		bval = b.Get([]byte("avg"))
		stats.TotalAvgLength = MakeUint(bval, "global", "avg")

		bval = b.Get([]byte("avgcount"))
		stats.TotalTxtCount = MakeUint(bval, "global", "avgcount")

		// Load total replies counter
		bval = b.Get([]byte("replies"))
		stats.Replies = MakeUint(bval, "global", "replies")

		// Load total replies counter
		bval = b.Get([]byte("forward"))
		stats.Forward = MakeUint(bval, "global", "forward")

		// Load hour counters
		b, err = tx.CreateBucketIfNotExists([]byte("hour"))
		if err != nil {
			return err
		}

		for i := 0; i < 24; i++ {
			bval = b.Get([]byte{byte(i)})
			stats.ByHour[i] = MakeUint(bval, "hour", strconv.Itoa(i))
		}

		// Load weekday counters
		b, err = tx.CreateBucketIfNotExists([]byte("weekday"))
		if err != nil {
			return err
		}

		for i := 0; i < 7; i++ {
			bval = b.Get([]byte{byte(i)})
			stats.ByWeekday[i] = MakeUint(bval, "weekday", strconv.Itoa(i))
		}

		// Load today's message counter, if possible
		b, err = tx.CreateBucketIfNotExists([]byte("date"))
		if err != nil {
			return err
		}

		todayKey := stats.TodayDate.Format("2006-1-2")
		bval = b.Get([]byte(todayKey))
		stats.Today = MakeUint(bval, "date", todayKey)

		// Load user counters
		stats.ByUserCount = make(map[string]uint64)
		b, err = tx.CreateBucketIfNotExists([]byte("users-count"))
		if err != nil {
			return err
		}
		b.ForEach(func(user, messages []byte) error {
			stats.ByUserCount[string(user)] = MakeUint(messages, "users-count", string(user))
			return nil
		})

		stats.ByUserAvgLen = make(map[string]uint64)
		b, err = tx.CreateBucketIfNotExists([]byte("users-avg"))
		if err != nil {
			return err
		}
		b.ForEach(func(user, messages []byte) error {
			stats.ByUserAvgLen[string(user)] = MakeUint(messages, "users-avg", string(user))
			return nil
		})

		stats.ByUserAvgCount = make(map[string]uint64)
		b, err = tx.CreateBucketIfNotExists([]byte("users-avgcount"))
		if err != nil {
			return err
		}
		b.ForEach(func(user, messages []byte) error {
			stats.ByUserAvgCount[string(user)] = MakeUint(messages, "users-avgcount", string(user))
			return nil
		})

		// Load type counters
		b, err = tx.CreateBucketIfNotExists([]byte("types"))
		if err != nil {
			return err
		}
		for i := 0; i < MessageTypeMax; i++ {
			bval = b.Get([]byte{byte(i)})
			stats.ByType[i] = MakeUint(bval, "types", strconv.Itoa(i))
		}

		return nil
	})
	assert(err)
}

func updateDate() {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("date"))
		todayKey := stats.TodayDate.Format("2006-1-2")
		err := b.Put([]byte(todayKey), PutUint(stats.Today))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("[updateDate] Couldn't save last day stats: " + err.Error())
	}
	stats.TodayDate = time.Now()
	stats.Today = 0
}

func updateMean(currentMean, meanCount, newValue uint64) uint64 {
	return ((currentMean * meanCount) + newValue) / (meanCount + 1)
}

func updateStats(message tg.APIMessage) {
	//
	// Local update
	//

	// DB Update flags
	updatemean := false
	updatetype := 0
	updatereplies := false
	updateforward := false

	// Update total count
	stats.TotalCount++

	// Update individual user's count
	username := message.User.Username
	val, exists := stats.ByUserCount[username]
	if !exists {
		val = 0
	}
	stats.ByUserCount[username] = val + 1

	// Update time counters
	now := time.Now()
	hour := now.Hour()
	wday := now.Weekday()
	stats.ByHour[hour]++
	stats.ByWeekday[wday]++

	// Check for day reset
	if now.Day() != stats.TodayDate.Day() {
		updateDate()
	}
	stats.Today++

	// Text message
	if message.Text != nil {
		stats.ByType[MessageTypeText]++

		// Update total and individual average
		msglen := uint64(len(*(message.Text)))
		if stats.TotalTxtCount > 0 {
			stats.TotalAvgLength = updateMean(stats.TotalAvgLength, stats.TotalTxtCount, msglen)
			stats.TotalTxtCount++
		} else {
			stats.TotalAvgLength = msglen
			stats.TotalTxtCount = 1
		}
		val, exists = stats.ByUserAvgCount[username]
		if exists {
			stats.ByUserAvgLen[username] = updateMean(stats.ByUserAvgLen[username], val, msglen)
			stats.ByUserAvgCount[username]++
		} else {
			stats.ByUserAvgLen[username] = msglen
			stats.ByUserAvgCount[username] = 1
		}
		updatemean = true
		updatetype = MessageTypeText
	}
	// Audio message
	if message.Audio != nil {
		stats.ByType[MessageTypeAudio]++
		updatetype = MessageTypeAudio
	}
	// Photo
	if message.Photo != nil {
		stats.ByType[MessageTypePhoto]++
		updatetype = MessageTypePhoto
	}
	// Sticker
	if message.Sticker != nil {
		stats.ByType[MessageTypeSticker]++
		updatetype = MessageTypeSticker
	}
	// Video
	if message.Video != nil {
		stats.ByType[MessageTypeVideo]++
		updatetype = MessageTypeVideo
	}
	// Voice message
	if message.Voice != nil {
		stats.ByType[MessageTypeVoice]++
		updatetype = MessageTypeVoice
	}
	// Contact
	if message.Contact != nil {
		stats.ByType[MessageTypeContact]++
		updatetype = MessageTypeContact
	}
	// Location
	if message.Location != nil {
		stats.ByType[MessageTypeLocation]++
		updatetype = MessageTypeLocation
	}
	// Document
	if message.Document != nil {
		stats.ByType[MessageTypeDocument]++
		updatetype = MessageTypeDocument
	}
	// Reply
	if message.ReplyTo != nil {
		stats.Replies++
		updatereplies = true
	}
	// Forwarded message
	if message.FwdUser != nil {
		stats.Forward++
		updateforward = true
	}

	//
	// DB Update
	//

	err := db.Update(func(tx *bolt.Tx) error {
		// Update total counters
		b := tx.Bucket([]byte("global"))

		err := b.Put([]byte("count"), PutUint(stats.TotalCount))
		if err != nil {
			return err
		}

		if updatemean {
			err = b.Put([]byte("avg"), PutUint(stats.TotalAvgLength))
			if err != nil {
				return err
			}
			err = b.Put([]byte("avgcount"), PutUint(stats.TotalTxtCount))
			if err != nil {
				return err
			}
		}

		if updatereplies {
			err = b.Put([]byte("replies"), PutUint(stats.Replies))
			if err != nil {
				return err
			}
		}
		if updateforward {
			err = b.Put([]byte("forward"), PutUint(stats.Forward))
			if err != nil {
				return err
			}
		}

		// Update time counters
		b = tx.Bucket([]byte("hour"))
		err = b.Put([]byte{byte(hour)}, PutUint(stats.ByHour[hour]))
		if err != nil {
			return err
		}

		b = tx.Bucket([]byte("weekday"))
		err = b.Put([]byte{byte(wday)}, PutUint(stats.ByHour[wday]))
		if err != nil {
			return err
		}

		b = tx.Bucket([]byte("date"))
		todayKey := stats.TodayDate.Format("2006-1-2")
		err = b.Put([]byte(todayKey), PutUint(stats.Today))
		if err != nil {
			return err
		}

		// Update user counters
		b = tx.Bucket([]byte("users-count"))
		err = b.Put([]byte(username), PutUint(stats.ByUserCount[username]))
		if err != nil {
			return err
		}

		if updatemean {
			b = tx.Bucket([]byte("users-avg"))
			err = b.Put([]byte(username), PutUint(stats.ByUserAvgLen[username]))
			if err != nil {
				return err
			}
			b = tx.Bucket([]byte("users-avgcount"))
			err = b.Put([]byte(username), PutUint(stats.ByUserAvgCount[username]))
			if err != nil {
				return err
			}
		}

		// Update type counter
		b = tx.Bucket([]byte("types"))
		err = b.Put([]byte{byte(updatetype)}, PutUint(stats.ByType[updatetype]))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("[updateStats] Got error while updating DB: " + err.Error())
	}
}
