package domain

type Pos struct {
	X  float64 `json:"x" bson:"x"`
	Y  float64 `json:"y" bson:"y"`
	Z  int     `json:"z" bson:"z"`
	A  int     `json:"a" bson:"a"`
	S  int     `json:"s" bson:"s"`
	St int     `json:"st" bson:"st"`
}

type Message struct {
	ID      int64                  `json:"id" bson:"id"`
	T       int                    `json:"t" bson:"t"`
	ST      int                    `json:"st" bson:"st"`
	Pos     Pos                    `json:"pos" bson:"pos"`
	Params  map[string]interface{} `json:"p" bson:"p"`
	Address string                 `json:"address,omitempty"` // <--- добавили сюда
}
