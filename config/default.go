package config

// nolint:gomnd
func Default() Config {
	return Config{
		Logger: Logger{
			Level: "info",
		},
	}
}
