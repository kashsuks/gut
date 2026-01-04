package cli

import (
	"fmt"
	"gut/internal/store"
	"github.com/spf13/cobra"
	"os"
)

func NewStartCommand() *cobra.Command {
	return &cobra.Command {
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {
			cwd, _ := os.Getwd()
			fs := store.NewFileStore(cwd)
			if err := fs.InitLayout(); err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Gut started!")
		},
	}
}
