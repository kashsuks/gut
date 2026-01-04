package cli

import "github.com/spf13/cobra"

func Execute() error {
	rootCmd := &cobra.Command{Use: "gut"}
	rootCmd.AddCommand(NewStartCommand(), NewSnapCommand())
	return rootCmd.Execute()
}
