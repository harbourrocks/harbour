package main

import (
	"os"

	"github.com/harbourrocks/harbour/cmd/harbour-iam/app"
)

func main() {
	cmd := app.NewIAMServerCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
