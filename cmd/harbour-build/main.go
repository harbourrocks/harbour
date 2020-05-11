package main

import (
	"github.com/harbourrocks/harbour/cmd/harbour-build/app"
	"os"
)

func main() {
	cmd := app.NewBuildServerCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
