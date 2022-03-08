package cmd

import (
	"github.com/ali-a-a/gophermq/config"
	"github.com/ali-a-a/gophermq/internal/app/gophermq/cmd/broker"
	"github.com/ali-a-a/gophermq/pkg/log"

	"github.com/spf13/cobra"
)

// NewRootCommand creates a new gophermq root command.
func NewRootCommand() *cobra.Command {
	var root = &cobra.Command{
		Use: "gophermq",
	}

	cfg := config.Init()

	log.SetupLogger(log.AppLogger{
		Level:  cfg.Logger.Level,
		StdOut: true,
	})

	broker.Register(root, cfg)

	return root
}
