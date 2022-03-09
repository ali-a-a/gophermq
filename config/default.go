package config

// nolint:gomnd
func Default() Config {
	return Config{
		Logger: Logger{
			Level: "info",
		},
		Broker: Broker{
			Port: 8082,
			MaxPending: 10,
			Subjects: []string{"test.a", "test.b", "test.a", "test.c"},
		},
	}
}
