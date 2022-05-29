package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/derhabicht/metoc/orchestrations"
)

var (
	tzOffset int

	rootCmd = &cobra.Command{
		Use:   "metoc",
		Short: "Generate METOC reports for operational planning",
		Long:  "",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rg, err := orchestrations.NewReportGenerator(args[0])
			if err != nil {
				return err
			}

			out, err := rg.Generate()
			if err != nil {
				return err
			}

			fmt.Println(out)
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
