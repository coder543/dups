package main

import "github.com/spf13/cobra"

func addCmd(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
	cmd.Flags().BoolP("single-core", "s", false, "use single cpu core")
	cmd.Flags().Int64("min-size", 1024, "minimum file size to scan in bytes")
}
