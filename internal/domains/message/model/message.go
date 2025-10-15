package model

type Pos struct {
	X  float64 `json:"x" bson:"x"`   // Longitude
	Y  float64 `json:"y" bson:"y"`   // Latitude
	Z  int     `json:"z" bson:"z"`   // Height
	A  int     `json:"a" bson:"a"`   // Azimuth
	S  int     `json:"s" bson:"s"`   // Speed
	Sl int     `json:"sl" bson:"sl"` // Satellites
}

type Message struct {
	ID int64 `json:"id"`           // Internal object ID
	DT int64 `json:"dt" bson:"dt"` // Device Time
	ST int64 `json:"st" bson:"st"` // Server Time

	Pos     Pos                    `json:"pos" bson:"pos"`
	Params  map[string]interface{} `json:"p" bson:"p"`
	Address string                 `json:"address" bson:"address"`

	//T time.Time `json:"t" bson:"-"` // Время отправки в ISO 8601 формате (RFC 3339 с миллисекундами)
}
