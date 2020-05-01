package app

import (
	"github.com/harbourrocks/harbour/pkg/harbourscm"
	"github.com/harbourrocks/harbour/pkg/harbourscm/configuration"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewSCMServerCommmand creates a *cobra.Command object with default parameters
func NewSCMServerCommmand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "harbour-scm",
		Long: `The harbour.rocks SCM server manages
external version control repositories.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// load SCM config
			s := configuration.ParseViperConfig()

			// load redis config
			s.Redis = redisconfig.ParseViperConfig()

			// configure logging
			l := logconfig.ParseViperConfig()
			logconfig.ConfigureLog(l)

			logrus.Info("Harbour SCM configured")

			// test redis connection
			redisconfig.TestConnection(s.Redis)

			return harbourscm.RunSCMServer(s)
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
