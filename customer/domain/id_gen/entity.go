package idgen

type IDGen struct {
	Key   string `bson:"key"`
	Value int64  `bson:"value"`
}
