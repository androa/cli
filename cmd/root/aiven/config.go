package aiven

import (
	"github.com/nais/cli/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AivenConfig struct {
	aiven  *cobra.Command
	create *cobra.Command
	get    *cobra.Command
	tidy   *cobra.Command
}

func NewAivenConfig() *AivenConfig {
	return &AivenConfig{
		aiven:  aivenCommand,
		create: createCmd,
		get:    getCmd,
		tidy:   tidyCmd,
	}
}

func (a AivenConfig) InitCmds(root *cobra.Command) {
	a.create.Flags().StringP(cmd.PoolFlag, "p", "nav-dev", "Preferred kafka pool to connect (optional)")
	viper.BindPFlag(cmd.PoolFlag, a.create.Flags().Lookup(cmd.PoolFlag))

	a.create.Flags().IntP(cmd.ExpireFlag, "e", 1, "Time in days the created secret should be valid (optional)")
	viper.BindPFlag(cmd.ExpireFlag, a.create.Flags().Lookup(cmd.ExpireFlag))

	a.create.Flags().StringP(cmd.SecretNameFlag, "s", "", "Preferred secret-name instead of generated (optional)")
	viper.BindPFlag(cmd.SecretNameFlag, a.create.Flags().Lookup(cmd.SecretNameFlag))

	a.get.Flags().StringP(cmd.ConfigFlag, "c", "all", "Type of config to generate. Supported values: .env, kcat, all (optional)")
	viper.BindPFlag(cmd.ConfigFlag, a.get.Flags().Lookup(cmd.ConfigFlag))

	root.AddCommand(a.aiven)
	a.aiven.AddCommand(a.create)
	a.aiven.AddCommand(a.get)
	a.aiven.AddCommand(a.tidy)
}