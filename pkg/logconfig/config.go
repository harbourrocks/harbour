package logconfig

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ParseViperConfig tries to map a viper configuration to LoggingOptions
func ParseViperConfig() *LoggingOptions {
	l := NewDefaultLoggingOptions()

	if v := viper.GetString("LOG_LEVEL"); v != "" {
		l.Level = v
	}

	l.ShowFunctionName = viper.GetBool("SHOW_FUNCTION_NAME")

	return l
}

// ConfigureLog configures a global logging solution.
func ConfigureLog(o *LoggingOptions) {
	logrus.SetReportCaller(o.ShowFunctionName)

	logrus.SetOutput(os.Stdout)

	if l, err := logrus.ParseLevel(o.Level); err != nil {
		logrus.Fatal("Failed to parse log level: %v", err)
	} else {
		logrus.SetLevel(l)
	}
}
