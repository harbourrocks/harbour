package app

import (
	"github.com/harbourrocks/harbour/pkg/harbourbuild"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/configuration"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewBuildServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "harbour-builder",
		Long: "The habour.rocks build server manages the building of projects for the registry",
		RunE: func(cmd *cobra.Command, args []string) error {

			// load builder config
			s := configuration.ParseViperConfig()

			// configure logging
			l := logconfig.ParseViperConfig()
			logconfig.ConfigureLog(l)

			logrus.Info("Harbour Builder configured")

			// test redis connection
			redisconfig.TestConnection(s.Redis)

			return harbourbuild.RunBuildServer(s)
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
