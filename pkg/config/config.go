package config

import (
	"github.com/julienschmidt/httprouter"
)

type HttpConfig struct {
	ListenPort string
	Router     *httprouter.Router
}

type NatsConfig struct {
	ConsumerId string
	NatsURL    string
	Subject    string
	Stream     string
}
