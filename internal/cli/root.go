package cli

import "github.com/spf13/cobra"

func Execute() error {
	rootCmd := &cobra.Command{
		Use: "gut",
		Short: "Gut - git but better",
	}
	rootCmd.AddCommand(
		NewStartCommand(),
		NewSnapCommand(),
		NewShowCommand(),
		NewListCommand(),
	)
	return rootCmd.Execute()
}
