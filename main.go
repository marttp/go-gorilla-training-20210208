package main

import (
	"./config"
)

func main() {
	configData := config.GetConfig()
	app := &App{}
	app.Initialize(configData)
	app.Run(":3000")
}
