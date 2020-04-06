package app

import (
	"github.com/harbourrocks/harbour/pkg/harbouriam"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewIAMServerCommmand creates a *cobra.Command object with default parameters
func NewIAMServerCommmand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "harbour-iam",
		Long: `The harbour.rocks IAM server manages
authentication and authorization for the harbour environment.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			s := harbouriam.NewDefaultOptions()
			s.OIDCClientID = viper.GetString("OIDC_CLIENT_ID")
			s.OIDCClientSecret = viper.GetString("OIDC_CLIENT_SECRET")
			s.OIDCURL = viper.GetString("OIDC_URL")
			s.IAMBaseURL = viper.GetString("IAM_BASE_URL")

			// load redis config
			s.Redis = redisconfig.ParseViperConfig()

			// configure logging
			l := logconfig.ParseViperConfig()
			logconfig.ConfigureLog(l)

			logrus.Info("Harbour IAM configured")

			// test redis connection
			redisClient := redisconfig.OpenClient(s.Redis)
			if pong, err := redisClient.Ping().Result(); err != nil {
				logrus.Fatal("Failed to open redis connection: ", err)
			} else {
				logrus.Info("Redis connection ok: ", pong)
				redisClient.Close()
			}

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
