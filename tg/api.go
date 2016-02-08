package tg

type APIUser struct {
	UserID    int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

type ChatType string

const (
	ChatTypePrivate    ChatType = "private"
	ChatTypeGroup      ChatType = "group"
	ChatTypeSupergroup ChatType = "supergroup"
	ChatTypeChannel    ChatType = "channel"
)

type APIChat struct {
	UserID    int      `json:"id"`
	Type      ChatType `json:"type"`
	Title     string   `json:"title,omitempty"`
	Username  string   `json:"username,omitempty"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
}

type APIMessage struct {
	MessageID         int            `json:"message_id"`
	User              APIUser        `json:"from"`
	Time              int            `json:"date"`
	Chat              APIChat        `json:"chat"`
	FwdUser           APIUpdate      `json:"forward_from,omitempty"`
	FwdTime           int            `json:"forward_date,omitempty"`
	ReplyTo           APIMessage     `json:"reply_to_message,omitempty"`
	Text              string         `json:"text,omitempty"`
	Audio             APIAudio       `json:"audio,omitempty"`
	Document          APIDocument    `json:"document,omitempty"`
	Photo             []APIPhotoSize `json:"photo,omitempty"`
	Sticker           APISticker     `json:"sticker,omitempty"`
	Video             APIVideo       `json:"video,omitempty"`
	Voice             APIVoice       `json:"voice,omitempty"`
	Caption           string         `json:"caption,omitempty"`
	Contact           APIContact     `json:"contact,omitempty"`
	Location          APILocation    `json:"location,omitempty"`
	NewUser           APIUser        `json:"new_chat_partecipant",omitempty"`
	LeftUser          APIUser        `json:"left_chat_partecipant,omitempty"`
	PhotoDeleted      bool           `json:"delete_chat_photo,omitempty"`
	GroupCreated      bool           `json:"group_chat_created,omitempty"`
	SupergroupCreated bool           `json:"supergroup_chat_created,omitempty"`
	ChannelCreated    bool           `json:"channel_chat_created,omitempty"`
	GroupToSuper      int            `json:"migrate_to_chat_id,omitempty"`
	GroupFromSuper    int            `json:"migrate_from_chat_id,omitempty"`
}

type APIPhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size,omitempty"`
}

type APIAudio struct {
	FileID    string `json:"file_id"`
	Duration  int    `json:"duration"`
	Performer string `json:"performer,omitempty"`
	Title     string `json:"title,omitempty"`
	MimeType  string `json:"mime_type,omitempty"`
	FileSize  int    `json:"file_size,omitempty"`
}

type APIDocument struct {
	FileID    string       `json:"file_id"`
	Thumbnail APIPhotoSize `json:"thumb,omitempty"`
	Filename  string       `json:"file_name"`
	MimeType  string       `json:"mime_type,omitempty"`
	FileSize  int          `json:"file_size,omitempty"`
}

type APISticker struct {
	FileID    string       `json:"file_id"`
	Width     int          `json:"width"`
	Height    int          `json:"height"`
	Thumbnail APIPhotoSize `json:"thumb,omitempty"`
	FileSize  int          `json:"file_size,omitempty"`
}

type APIVideo struct {
	FileID    string       `json:"file_id"`
	Width     int          `json:"width"`
	Height    int          `json:"height"`
	Duration  int          `json:"duration"`
	Thumbnail APIPhotoSize `json:"thumb,omitempty"`
	MimeType  string       `json:"mime_type,omitempty"`
	FileSize  int          `json:"file_size,omitempty"`
}

type APIVoice struct {
	FileID   string `json:"file_id"`
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type,omitempty"`
	FileSize int    `json:"file_size,omitempty"`
}

type APIContact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	UserID      int    `json:"user_id,omitempty"`
}

type APILocation struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type APIUpdate struct {
	UpdateID int        `json:"update_id"`
	Message  APIMessage `json:"message"`
}
