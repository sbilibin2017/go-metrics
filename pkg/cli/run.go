package cli

import "github.com/spf13/cobra"

func Run(cmd *cobra.Command) int {
	if err := cmd.Execute(); err != nil {
		return 1
	}
	return 0
}
