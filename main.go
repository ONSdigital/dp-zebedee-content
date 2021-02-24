package main

import (
	"os"

	"github.com/ONSdigital/dp-zebedee-content/commands"
	log "github.com/daiLlew/funkylog"
)

func main() {
	log.Init("dp-zebedee-content")

	err := commands.GetRootCommand().Execute()
	if err != nil {
		log.Err("unexpected error executing cli command :violin:: %+v", err)
		os.Exit(1)
	}
}