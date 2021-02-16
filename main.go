package main

import (
	"os"
	"runtime/debug"

	"github.com/ONSdigital/dp-zebedee-content/commands"
	"github.com/daiLlew/cli-fmt/log"
)

func main() {
	log.Init("dp-zebedee-content")

	err := commands.GetRootCommand().Execute()
	if err != nil {
		log.Err("unexpected error executing cli command :violin:: %s", err.Error())
		debug.PrintStack()
		os.Exit(1)
	}
}