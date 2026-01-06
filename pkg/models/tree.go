package models

import (
	"fmt"
	"sort"
	"strings"
)

type FileMode uint32

const (
	ModeBlob       FileMode = 0100644 // normal file 
	ModeExecutable FileMode = 0100755 // executable file 
	ModeTree       FileMode = 0040000 // directory/tree
)

type TreeEntry struct {
	Mode FileMode
	Name string 
	Hash string 
	Type ObjectType
}

// the tree represents a complete tree structure
type Tree struct {
	Entries []TreeEntry
}

func NewTree() *Tree {
	return &Tree{
		Entries: make([]TreeEntry, 0),
	}
}

func (t *Tree) AddEntry(mode FileMode, name string, hash string, objType ObjectType) {
	t.Entries = append(t.Entries, TreeEntry{
		Mode: mode,
		Name: name,
		Hash: hash,
		Type: objType,
	})
}

func (t *Tree) Sort() {
	sort.Slice(t.Entries, func(i, j int) bool {
		return t.Entries[i].Name < t.Entries[j].Name
	})
}

// serialize converts the tree to the normal git tree format
// <mode> <name>\0<20-byte-sha>
// we can use a full hex hash instead of binary
func (t *Tree) Serialize() []byte {
	var builder strings.Builder 

	t.Sort()

	for _, entry := range t.Entries {
		builder.WriteString(fmt.Sprintf("%o %s\x00%s\n", entry.Mode, entry.Name, entry.Hash))
	}

	return []byte(builder.String())
}

func (t *Tree) String() string {
	var builder strings.Builder 
	builder.WriteString("Tree with entries:\n")

	for _, entry := range t.Entries {
		builder.WriteString(fmt.Sprintf("  %o %s %s (%s)\n",
			entry.Mode, entry.Type, entry.Name, entry.Hash[:8]))
	}

	return builder.String()
}

// GetMode determines the file mode from os.FileMode 
func GetMode(isExecutable bool) FileMode {
	if isExecutable {
		return ModeExecutable
	}

	return ModeBlob
}