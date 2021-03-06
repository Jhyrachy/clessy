package main

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"strconv"
	"strings"
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
	ByUserCount map[string]uint64
	ByWeekday   [7]uint64
	ByHour      [24]uint64
	ByType      [MessageTypeMax]uint64
	ByDay       map[string]uint64
	TodayDate   time.Time
	Today       uint64
	TotalCount  uint64
}

var stats Stats

type UserCount map[string]uint64

var words map[string]UserCount

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
		stats.TotalCount = MakeUint(b.Get([]byte("count")), "global", "count")

		// Load hour counters
		b, err = tx.CreateBucketIfNotExists([]byte("hour"))
		if err != nil {
			return err
		}

		for i := 0; i < 24; i++ {
			stats.ByHour[i] = MakeUint(b.Get([]byte{byte(i)}), "hour", strconv.Itoa(i))
		}

		// Load weekday counters
		b, err = tx.CreateBucketIfNotExists([]byte("weekday"))
		if err != nil {
			return err
		}

		for i := 0; i < 7; i++ {
			stats.ByWeekday[i] = MakeUint(b.Get([]byte{byte(i)}), "weekday", strconv.Itoa(i))
		}

		// Load day counters
		stats.ByDay = make(map[string]uint64)
		b, err = tx.CreateBucketIfNotExists([]byte("date"))
		if err != nil {
			return err
		}

		b.ForEach(func(day, messages []byte) error {
			stats.ByDay[string(day)] = MakeUint(messages, "date", string(day))
			return nil
		})

		todayKey := stats.TodayDate.Format("2006-1-2")
		stats.Today = MakeUint(b.Get([]byte(todayKey)), "date", todayKey)

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

		// Load type counters
		b, err = tx.CreateBucketIfNotExists([]byte("types"))
		if err != nil {
			return err
		}
		for i := 0; i < MessageTypeMax; i++ {
			stats.ByType[i] = MakeUint(b.Get([]byte{byte(i)}), "types", strconv.Itoa(i))
		}

		// Load dictionary
		b, err = tx.CreateBucketIfNotExists([]byte("words"))
		if err != nil {
			return err
		}
		words = make(map[string]UserCount)
		b.ForEach(func(word, ucount []byte) error {
			var val UserCount
			err := json.Unmarshal(ucount, &val)
			if err != nil {
				return err
			}
			words[string(word)] = val
			return nil
		})

		return nil
	})
	assert(err)
}

func updateDate() {
	dateKey := stats.TodayDate.Format("2006-1-2")
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("date"))
		err := b.Put([]byte(dateKey), PutUint(stats.Today))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("[updateDate] Couldn't save last day stats: " + err.Error())
	}
	stats.ByDay[dateKey] = stats.Today
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
	updatetype := 0

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
		updatetype = MessageTypeText

		// Process words
		processWords(message)
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

		// Update time counters
		b = tx.Bucket([]byte("hour"))
		err = b.Put([]byte{byte(hour)}, PutUint(stats.ByHour[hour]))
		if err != nil {
			return err
		}

		b = tx.Bucket([]byte("weekday"))
		err = b.Put([]byte{byte(wday)}, PutUint(stats.ByWeekday[wday]))
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

func processWords(message tg.APIMessage) {
	if len(*(message).Text) < 3 {
		return
	}

	wordList := strings.Split(*(message.Text), " ")
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("words"))
		for _, word := range wordList {
			if len(word) < 3 {
				continue
			}

			word = strings.ToLower(word)

			if strings.HasPrefix(word, "http") {
				continue
			}

			word = strings.Trim(word, " ?!.,:;/-_()[]{}'\"+=*^\n")
			count, ok := words[word]
			if !ok {
				count = make(UserCount)
			}
			val, ok := count[message.User.Username]
			if !ok {
				val = 0
			}
			count[message.User.Username] = val + 1
			words[word] = count

			j, err := json.Marshal(count)
			if err != nil {
				return err
			}
			b.Put([]byte(word), j)
		}
		return nil
	})
	if err != nil {
		log.Println("[processWords] Error encountered: " + err.Error())
	}
}

var FILTER = []string{
	"100", "abbastanza", "abbia", "abbiamo", "adesso", "again", "agli", "ah", "alcune",
	"alcuni", "all", "all'inizio", "alla", "alle", "allo", "allora", "almeno", "also",
	"alto", "altra", "altre", "altri", "altrimenti", "altro", "amici", "amico", "amo",
	"anche", "ancora", "and", "andare", "andato", "anime", "anni", "anzi", "appena", "apposta",
	"are", "assieme", "avanti", "aver", "avere", "avete", "aveva", "avevano", "avevi", "avevo",
	"avrebbe", "avrei", "avuto", "base", "bel", "bella", "belle", "belli", "bellissimo",
	"bello", "ben", "bene", "benissimo", "bisogno", "bravo", "brutta", "brutto", "cambia",
	"che", "chi", "cioe", "cioè", "ciò", "coi", "col", "com'è", "come", "con", "cos'è", "cosa",
	"così", "cui", "dai", "dal", "dalla", "dalle", "danno", "dare", "degli", "dei", "del",
	"della", "delle", "dello", "deve", "devi", "devo", "dove", "e", "era", "erano", "eri",
	"ero", "fa", "fai", "fanno", "finché", "gia", "già", "giù", "gli", "hai", "han", "hanno",
	"have", "il", "in", "io", "l'altro", "l'avevo", "l'ha", "l'hai", "l'hanno", "l'ho",
	"la", "lei", "lui", "lì", "ma", "me", "meno", "mentre", "mia", "mie", "miei", "mio", "molti",
	"molto", "negli", "nei", "nel", "nella", "nelle", "nello", "no", "noi", "non", "not", "nuovi",
	"nuovo", "ok", "oltre", "oppure", "ora", "per", "perche", "perchè", "perché", "però",
	"piu", "più", "po", "poi", "puoi", "pure", "può", "qua", "qualche", "quale", "quando",
	"quanti", "quanto", "quasi", "quei", "quel", "quella", "quelle", "quelli", "quello",
	"questa", "queste", "questi", "questo", "qui", "quindi", "sai", "sarei", "sarà", "se",
	"sei", "sempre", "sennò", "senza", "si", "sia", "siamo", "siano", "siete", "son", "sono",
	"sopra", "sta", "stai", "ste", "sti", "stiamo", "sto", "sua", "sue", "sui", "sul", "sulla",
	"sulle", "suo", "suoi", "sì", "tanta", "tante", "tanti", "tanto", "te", "that", "the",
	"then", "too", "tra", "troppi", "troppo", "tua", "tuo", "tuoi", "tutta", "tutte",
	"tutti", "tutto", "un'altra", "una", "uno", "usa", "usi", "uso", "vai", "verso", "via",
	"voglia", "voglio", "vogliono", "voi", "volete", "voleva", "volevo", "volta", "volte",
	"vorrei", "vuoi", "vuol", "vuole", "was",
}

const USAGE_THRESHOLD = 10

func filteredWords() map[string]UserCount {
	filtered := make(map[string]UserCount)
	for word, usage := range words {
		// Check for too common

		isfilter := false
		for _, filter := range FILTER {
			if word == filter {
				isfilter = true
				break
			}
		}
		if isfilter {
			continue
		}

		// Check for not common enough
		good := false
		ucount := make(UserCount)
		for user, count := range usage {
			if count < USAGE_THRESHOLD {
				continue
			}
			if !good {
				good = true
			}
			ucount[user] = count
		}
		if !good {
			continue
		}

		filtered[word] = usage
	}
	return filtered
}
