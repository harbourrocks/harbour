package app

import (
	server "github.com/harbourrocks/harbour/pkg/harbourgateway"
	"github.com/harbourrocks/harbour/pkg/harbourgateway/configuration"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewGatewayServerCommand creates a *cobra.Command object with default parameters
func NewGatewayServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "harbour-gateway",
		Long: `The harbour.rocks Gateway server manages
all incoming http requests.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// load OIDC  config
			s := configuration.ParseViperConfig()

			// configure logging
			l := logconfig.ParseViperConfig()
			logconfig.ConfigureLog(l)

			logrus.Info("Harbour Gateway configured")

			// test redis connection
			redisconfig.TestConnection(s.Redis)

			return server.RunGatewayServer(s)
		},
	}

	return cmd
}

func init() {
	cobra.OnInitialize(initCobra)
}

func initCobra() {
	viper.AutomaticEnv()
}
