package main

import (
	"context"
	"os"

	"github.com/ONSdigital/dp-zebedee-content/commands"
	"github.com/ONSdigital/log.go/log"
)

func main() {
	log.Namespace = "dp-zebedee-content"

	err := commands.GetRootCommand().Execute()
	if err != nil {
		log.Event(context.Background(), "unexpected error executing cli command", log.Error(err))
		os.Exit(1)
	}
}
