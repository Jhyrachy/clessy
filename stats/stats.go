package main

import (
	"encoding/binary"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
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
	MessageTypeMax      int = 8
)

type Stats struct {
	ByUserCount    map[string]uint64
	ByUserAvgLen   map[string]uint64
	ByWeekday      [7]uint64
	ByHour         [24]uint64
	ByType         [MessageTypeMax]uint64
	TodayDate      time.Time
	Today          uint64
	TotalCount     uint64
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

func loadStats() {
	// Load today
	stats.TodayDate = time.Now()

	err := db.View(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("global"))
		if err != nil {
			return err
		}

		// Load total messages counter
		bval := bucket.Get([]byte("count"))
		stats.TotalCount = MakeUint(bval, "global", "count")

		// Load total messages counter
		bval := bucket.Get([]byte("avg"))
		stats.TotalAvgLength = MakeUint(bval, "global", "avg")

		// Load total replies counter
		bval = bucket.Get([]byte("replies"))
		stats.Replies = MakeUint(bval, "global", "replies")

		// Load total replies counter
		bval = bucket.Get([]byte("forward"))
		stats.Forward = MakeUint(bval, "global", "forward")

		// Load hour counters
		b, err = tx.CreateBucketIfNotExists([]byte("hour"))
		if err != nil {
			return err
		}

		for i := 0; i < 24; i++ {
			bval = bucket.Get([]byte(i))
			stats.ByHour[i] = MakeUint(bval, "hour", strconv.Itoa(i))
		}

		// Load weekday counters
		b, err = tx.CreateBucketIfNotExists([]byte("weekday"))
		if err != nil {
			return err
		}

		for i := 0; i < 7; i++ {
			bval = bucket.Get([]byte(i))
			stats.ByWeekday[i] = MakeUint(bval, "weekday", strconv.Itoa(i))
		}

		// Load today's message counter, if possible
		b, err = tx.CreateBucketIfNotExists([]byte("date"))
		if err != nil {
			return err
		}

		todayKey := stats.TodayDate.Format("2006-1-2")
		bval = bucket.Get([]byte(todayKey))
		stats.Today = MakeUint(bval, "date", todayKey)

		// Load user counters
		stats.ByUserCount = make(map[string]uint64)
		b, err = tx.CreateBucketIfNotExists([]byte("users-count"))
		if err != nil {
			return err
		}
		b.ForEach(func(user, messages []byte) error {
			stats.ByUserCount[string(user)] = MakeUint(messages, "users-count", string(user))
		})

		stats.ByUserAvgLen = make(map[string]uint64)
		b, err = tx.CreateBucketIfNotExists([]byte("users-avg"))
		if err != nil {
			return err
		}
		b.ForEach(func(user, messages []byte) error {
			stats.ByUserAvgLen[string(user)] = MakeUint(messages, "users-avg", string(user))
		})

		// Load type counters
		b, err = tx.CreateBucketIfNotExists([]byte("types"))
		if err != nil {
			return err
		}
		for i := 0; i < MessageTypeMax; i++ {
			bval = bucket.Get([]byte(i))
			stats.ByType[i] = MakeUint(bval, "types", strconv.Itoa(i))
		}

		return nil
	})
	assert(err)
}
