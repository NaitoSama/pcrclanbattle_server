package config

import (
	"github.com/BurntSushi/toml"
	"pcrclanbattle_server/common"
)

var Config config

func ConfigInit() {
	_, err := toml.DecodeFile("./config/config.toml", &Config)
	if err != nil {
		panic(err)
		return
	}
	common.Logln(0, "config init")
}

type config struct {
	General general
	DB      db
	Boss    boss
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
