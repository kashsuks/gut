package cli

import (
	"fmt"
	"gut/internal/core"
	"gut/internal/store"
	"gut/pkg/models"
	"os"
	"path/filepath"
	"github.com/spf13/cobra"
)

func NewSnapCommand() *cobra.Command {
	return &cobra.Command{
		Use: "snap [path]",
		Short: "Snapshot a file or directory",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			targetPath := args[0]

			// get the file or directory info
			info, err := os.Stat(targetPath)
			if err != nil {
				fmt.Printf("Error: path does not exist: %s\n", targetPath)
				return
			}

			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error: failed to get working directory: %v\n", err)
				return
			}

			fs := store.NewFileStore(cwd)

			if info.IsDir() {
				if err := snapDirectory(fs, targetPath); err != nil {
					fmt.Printf("Error: %v\n", err)
					return
				}
			} else {
				// handle single file Snapshot
				if err := snapFile(fs, targetPath); err != nil {
					fmt.Printf("Error: %v\n", err)
					return
				}
			}
		},
	}
}

func snapFile(fs *store.FileStore, filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	obj := core.HashObject(data, models.BlobObject)

	if err := fs.WriteObject(obj); err != nil {
		return fmt.Errorf("failed to write object: %w", err)
	}

	fileName := filepath.Base(filepath)
	fmt.Printf("Snapped file: %s\n", fileName)
	fmt.Printf("	Type: %s\n", obj.Type)
	fmt.Printf("	Hash: %s\n", obj.HashSum)
	fmt.Printf("	Size: %d bytes\n", obj.Size)

	return nil
}

func snapDirectory(fs *store.FileStore, dirPath string) error {
	builder := core.NewTreeBuilder(fs)

	treeHash, err := builder.BuildTreeFromDirectory(dirPath)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err)
	}

	absPath, _ := filepath.Abs(dirPath)
	dirName := filepath.Base(absPath)

	fmt.Printf("Snapped directory: %s\n", dirName)
	fmt.Printf("	Type: tree\n")
	fmt.Printf("	Hash: %s\n", treeHash)
	fmt.Println("\nAll files and subdirectories have been stored in .gut/objects/")
	fmt.Println("Use this hash to recreate the entire directory structure later.")

	return nil
}
