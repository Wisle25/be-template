package main

import (
	"fmt"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/infrastructures/server"
)

func main() {
	config := commons.LoadConfig(".")
	app := server.CreateServer(config)
	
	_ = app.Listen(fmt.Sprintf(":%s", config.ServerPort))
}
