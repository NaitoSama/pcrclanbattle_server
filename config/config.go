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
	General general
	DB      db
}

type general struct {
	HttpPort string
}

type db struct {
	AdminName   string
	AdminPasswd string
}

type boss struct {
	Stage            int
	StageOne         []int64
	StageTwo         []int64
	StageThree       []int64
	StageFour        []int64
	StageFive        []int64
	StageSix         []int64
	StageOneToTwo    int
	StageTwoToThree  int
	StageThreeToFour int
	StageFourToFive  int
	StageFiveToSix   int
}
