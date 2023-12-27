package main

import (
	"pcrclanbattle_server/config"
	"pcrclanbattle_server/db"
	"pcrclanbattle_server/router"
	"pcrclanbattle_server/server"
)

func main() {
	config.ConfigInit()
	db.DBInit()
	server.WSInit()
	router.RouterInit()
}
