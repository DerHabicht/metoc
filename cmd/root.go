package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/derhabicht/metoc/orchestration"
)

var (
	rootCmd = &cobra.Command{
		Use:   "metoc <PLAN YAML> <OUTPUT TEX>",
		Short: "Generate METOC reports for operational planning",
		Long:  "METOC v1.0.0",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			rg, err := orchestration.NewReportGenerator(args[0])
			if err != nil {
				return err
			}

			err = rg.Generate(args[1])
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".metoc")

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to read config %s", viper.ConfigFileUsed()))
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
