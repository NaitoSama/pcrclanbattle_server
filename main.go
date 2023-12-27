package main

import (
	"pcrclanbattle_server/config"
	"pcrclanbattle_server/router"
	"pcrclanbattle_server/server"
)

func main() {
	config.ConfigInit()
	server.WSInit()
	router.RouterInit()
}
