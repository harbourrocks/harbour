package main

import (
	"os"

	"github.com/harbourrocks/harbour/cmd/harbour-gateway/app"
)

func main() {
	cmd := app.NewGatewayServerCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
