package cli

import (
	"compress/zlib"
	"fmt"
	"gut/internal/store"
	"gut/pkg/models"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	return &cobra.Command{
		Use: "show [hash]",
		Short: "Display the contents of an object",
		Long: `Show the contents of a blob or tree objects by its hash.

		Examples:
		gut show ab123def456`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hash := args[0]

			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error: failed to get working directory: %v\n", err)
				return
			}

			fs := store.NewFileStore(cwd)

			if err := showObject(fs, hash); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		},
	}
}

func showObject(fs *store.FileStore, hash string) error {
	objPath := filepath.Join(fs.RootPath, ".gut", "objects", hash[:2], hash[2:])

	f, err := os.Open(objPath)
	if err != nil {
		return fmt.Errorf("object not found: %s", hash)
	}

	defer f.Close()

	zr, err := zlib.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to decompress object: %w", err)
	}

	defer zr.Close()

	content, err := io.ReadAll(zr)
	if err != nil {
		return fmt.Errorf("failed to read object: %w", err)
	}

	//parse the header
	objType, size, data := parseObject(content)

	fmt.Printf("Object: %s\n", hash)
	fmt.Printf("Type: %s\n", objType)
	fmt.Printf("Size: %d bytes\n\n", size)

	// display based on data type 
	switch objType {
	case models.BlobObject:
		fmt.Println("Content:")
		fmt.Println("─────────────────────────────────────")
		fmt.Println(string(data))
	case models.TreeObject:
		fmt.Println("Tree entries:")
		fmt.Println("─────────────────────────────────────")
		displayTree(data)
	default:
		fmt.Printf("Unknown object type: %s\n", objType)
	}

	return nil
}

func parseObject(content []byte) (models.ObjectType, int, []byte) {
	nullIdx := -1
	for i, b := range content {
		if b == 0 {
			nullIdx = i 
			break
		}
	}

	if nullIdx == -1 {
		return "", 0, content 	
	}

	header := string(content[:nullIdx])
	data := content[nullIdx+1:]

	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return "", 0, data
	}

	var size int
	fmt.Sscanf(parts[1], "%d", &size)

	return models.ObjectType(parts[0]), size, data
}

func displayTree(data []byte) {
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line[:nullIdx], " ", 2)
		if len(parts) != 2 {
			continue
		}

		mode := parts[0]
		name := parts[1]
		hash := line[nullIdx+1:]

		icon := "📄"
		typeStr := "blob"
		if mode == "40000" {
			icon = "📁"
			typeStr = "tree"
		} else if mode == "100755" {
			icon = "⚙️"
		}

		fmt.Printf("	%s %s %-6s %s %s\n", icon, mode, typeStr, name, hash[:8]+"...")
	}
}
