package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	app := newApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal("An error occurred running the application!")
	}
}
