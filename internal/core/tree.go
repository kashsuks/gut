package core

import (
	"fmt"
	"gut/internal/store"
	"gut/pkg/models"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type TreeBuilder struct {
	store *store.FileStore
}

func NewTreeBuilder(fileStore *store.FileStore) *TreeBuilder {
	return &TreeBuilder{
		store: fileStore,
	}
}

func (tb *TreeBuilder) BuildTreeFromDirectory(dirPath string) (string, error) {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	tree, err := tb.buildTree(absPath, absPath)
	if err != nil {
		return "", err
	}

	treeData := tree.Serialize()
	treeObj := HashObject([]byte(treeData), models.TreeObject)

	if err := tb.store.WriteObject(treeObj); err != nil {
		return "", fmt.Errorf("failed to write tree object: %w", err)
	}

	return treeObj.HashSum, nil
}

func (tb *TreeBuilder) buildTree(rootPath, currentPath string) (*models.Tree, error) {
	tree := models.NewTree()

	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(currentPath, entry.Name())

		if tb.shouldIgnore(fullPath, rootPath) {
			continue
		}

		if entry.IsDir() {
			subTree, err := tb.buildTree(rootPath, fullPath)
			if err != nil {
				return nil, err
			}

			subTreeData := subTree.Serialize()
			subTreeObj := HashObject(subTreeData, models.TreeObject)

			if err := tb.store.WriteObject(subTreeObj); err != nil {
				return nil, fmt.Errorf("failed to write subtree: %w", err)
			}

			tree.AddEntry(models.ModeTree, entry.Name(), subTreeObj.HashSum, models.TreeObject)
		} else {
			hash, err := tb.hashFile(fullPath)
			if err != nil {
				return nil, err
			}

			info, err := entry.Info()
			if err != nil {
				return nil, err
			}

			mode := tb.getFileMode(info.Mode())
			tree.AddEntry(mode, entry.Name(), hash, models.BlobObject)
		}
	}

	return tree, nil
}

func (tb *TreeBuilder) hashFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	obj := HashObject(data, models.BlobObject)

	if err := tb.store.WriteObject(obj); err != nil {
		return "", fmt.Errorf("failed to write blob: %w", err)
	}

	return obj.HashSum, nil
}

func (tb *TreeBuilder) getFileMode(mode fs.FileMode) models.FileMode {
	if mode&0111 != 0 {
		return models.ModeExecutable
	}
	return models.ModeBlob
}

func (tb *TreeBuilder) shouldIgnore(path, rootPath string) bool {
	if strings.Contains(path, string(filepath.Separator)+".gut") ||
		strings.HasSuffix(path, ".gut") {
		return true
	}

	return false
}
