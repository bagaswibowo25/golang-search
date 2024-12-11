package config

type HttpConfig struct {
	ListenPort string
}

type NatsConfig struct {
	ConsumerId string
	NatsURL    string
	Subject    string
	Stream     string
}
