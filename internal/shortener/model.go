package shortener

// Redirect has info about shortened link.
// Application is going to use
// Json, MongoDb, MsgPack as well as Redis.
// Tag "json" uses to send via http protocol
// Tag "bson" uses to store Redirection in MongoDb
// Tag "msgpack" uses to exchange data via Redis.
type Redirect struct {
	Code      string `json:"code" bson:"code" msgpack:"code"`
	URL       string `json:"url" bson:"url" msgpack:"url" validate:"empty=false & format=url"`
	CreatedAt int64  `json:"created_at" bson:"created_at" msgpack:"created_at"`
}

const (
	CodeField      = "code"
	URLField       = "url"
	CreatedAtField = "created_at"
)
