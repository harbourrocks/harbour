package logconfig

// LoggingOptions configures some options of the global logger.
type LoggingOptions struct {
	Level            string
	ShowFunctionName bool
}

// NewDefaultLoggingOptions returns the default logging options.
func NewDefaultLoggingOptions() *LoggingOptions {
	s := LoggingOptions{
		Level:            "info",
		ShowFunctionName: false,
	}

	return &s
}
