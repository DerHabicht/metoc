package cmd

import (
	"fmt"
	"github.com/derhabicht/metoc/orchestration"
	"os"

	"github.com/spf13/cobra"
)

var (
	tzOffset int

	rootCmd = &cobra.Command{
		Use:   "metoc <PLAN YAML> <OUTPUT TEX>",
		Short: "Generate METOC reports for operational planning",
		Long:  "",
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

func init() {
	rootCmd.Flags().IntVar(&tzOffset, "tzoffset", 0, "Timezone offset for the report (default: 0)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
