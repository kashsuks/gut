package cli

import (
	"fmt"
	"gut/internal/core"
	"gut/internal/store"
	"gut/pkg/models"
	"os"
	"github.com/spf13/cobra"
)

func NewSnapCommand() *cobra.Command {
	return &cobra.Command{
		Use: "snap [file]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			data, _ := os.ReadFile(args[0])
			obj := core.HashObject(data, models.BlobObject)

			cwd, _ := os.Getwd()
			fs := store.NewFileStore(cwd)
			fs.WriteObject(obj)

			fmt.Printf("Snapped: %s\n", obj.HashSum)
		},
	}
}
