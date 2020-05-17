package app

import (
	"github.com/harbourrocks/harbour/pkg/harbouriam"
	"github.com/harbourrocks/harbour/pkg/harbouriam/configuration"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewIAMServerCommand creates a *cobra.Command object with default parameters
func NewIAMServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "harbour-iam",
		Long: `The harbour.rocks IAM server manages
authentication and authorization for the harbour environment.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// load IAM config
			s := configuration.ParseViperConfig()

			// configure logging
			l := logconfig.ParseViperConfig()
			logconfig.ConfigureLog(l)

			logrus.Info("Harbour IAM configured")

			// test redis connection
			redisconfig.TestConnection(s.Redis)

			return harbouriam.RunIAMServer(s)
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
