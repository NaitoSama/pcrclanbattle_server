package config

import (
	"github.com/BurntSushi/toml"
)

var Config config

func ConfigInit() {
	_, err := toml.DecodeFile("./config/config.toml", &Config)
	if err != nil {
		panic(err)
		return
	}
}

type config struct {
	General General
}

type General struct {
	HttpPort string
}
