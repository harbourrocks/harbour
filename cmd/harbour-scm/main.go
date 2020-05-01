package main

import (
	"github.com/harbourrocks/harbour/cmd/harbour-scm/app"
	"os"
)

func main() {
	cmd := app.NewSCMServerCommmand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
