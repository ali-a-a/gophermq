package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"strings"

	"github.com/knadh/koanf"
	"github.com/sirupsen/logrus"
)

const Prefix = "GOPHERMQ_"

type (
	// Config represents application configuration struct.
	Config struct {
		Logger Logger `koanf:"logger"`
		Broker Broker `koanf:"broker"`
	}

	// Logger represents logger configuration struct.
	Logger struct {
		Level string `koanf:"level"`
	}

	// Broker represents Broker configuration struct.
	Broker struct {
		Port       int
		MaxPending int
		Subjects   []string
	}
)

func Init() Config {
	var cfg Config

	k := koanf.New(".")

	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		logrus.Fatalf("error loading default: %s", err)
	}

	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		logrus.Errorf("error loading config.yml: %s", err)
	}

	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		parsedEnv := strings.Replace(strings.ToLower(strings.TrimPrefix(s, Prefix)), "__", "-", -1)
		return strings.Replace(parsedEnv, "_", ".", -1)
	}), nil); err != nil {
		logrus.Errorf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &cfg); err != nil {
		logrus.Fatalf("error unmarshalling config: %s", err)
	}

	return cfg
}
