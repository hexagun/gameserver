package gameserver

type Action interface {
	GetType() string
	GetPayload() interface{}
}
