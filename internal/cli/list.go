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

func NewListCommand() *cobra.Command {
	var typeFilter string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all objects in the database",
		Long: `List all objects stored in the .gut/objects directory.

Examples:
  gut list              # List all objects
  gut list --type blob  # List only blob objects
  gut list --type tree  # List only tree objects`,
		Run: func(cmd *cobra.Command, args []string) {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error: failed to get working directory: %v\n", err)
				return
			}

			fs := store.NewFileStore(cwd)

			if err := listObjects(fs, typeFilter); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		},
	}

	cmd.Flags().StringVarP(&typeFilter, "type", "t", "", "Filter by object type (blob, tree, commit)")

	return cmd
}

func listObjects(fs *store.FileStore, typeFilter string) error {
	objectsPath := filepath.Join(fs.RootPath, ".gut", "objects")

	if _, err := os.Stat(objectsPath); os.IsNotExist(err) {
		fmt.Println("No objects found. Have you run 'gut start' and 'gut snap' yet?")
		return nil
	}

	objects := make([]ObjectInfo, 0)

	err := filepath.Walk(objectsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(objectsPath, path)
		parts := strings.Split(rel, string(filepath.Separator))
		if len(parts) != 2 {
			return nil
		}
		hash := parts[0] + parts[1]

		objType, size, err := readObjectType(path)
		if err != nil {
			return nil // skip objects we can't read
		}

		if typeFilter != "" && string(objType) != typeFilter {
			return nil
		}

		objects = append(objects, ObjectInfo{
			Hash: hash,
			Type: objType,
			Size: size,
		})

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk objects directory: %w", err)
	}

	if len(objects) == 0 {
		if typeFilter != "" {
			fmt.Printf("No objects of type '%s' found.\n", typeFilter)
		} else {
			fmt.Println("No objects found.")
		}
		return nil
	}

	fmt.Printf("Found %d object(s):\n\n", len(objects))
	fmt.Printf("%-10s  %-64s  %s\n", "TYPE", "HASH", "SIZE")
	fmt.Println(strings.Repeat("─", 85))

	for _, obj := range objects {
		icon := getTypeIcon(obj.Type)
		fmt.Printf("%s %-6s  %s  %d bytes\n", icon, obj.Type, obj.Hash, obj.Size)
	}

	return nil
}

type ObjectInfo struct {
	Hash string
	Type models.ObjectType
	Size int
}

func readObjectType(path string) (models.ObjectType, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()

	zr, err := zlib.NewReader(f)
	if err != nil {
		return "", 0, err
	}
	defer zr.Close()

	content, err := io.ReadAll(zr)
	if err != nil {
		return "", 0, err
	}

	nullIdx := -1
	for i, b := range content {
		if b == 0 {
			nullIdx = i
			break
		}
	}

	if nullIdx == -1 {
		return "", 0, fmt.Errorf("invalid object format")
	}

	header := string(content[:nullIdx])
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid header format")
	}

	var size int
	fmt.Sscanf(parts[1], "%d", &size)

	return models.ObjectType(parts[0]), size, nil
}

func getTypeIcon(objType models.ObjectType) string {
	switch objType {
	case models.BlobObject:
		return "📄"
	case models.TreeObject:
		return "📁"
	case models.CommitObject:
		return "💾"
	default:
		return "❓"
	}
}
