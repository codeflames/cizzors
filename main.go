package main

import (
	"github.com/codeflames/cizzors/models"
	"github.com/codeflames/cizzors/server"
)

func main() {
	models.Setup()
	server.SetupAndListen()
}
