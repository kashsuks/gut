package models

import "fmt"

type ObjectType string

const (
	BlobObject	ObjectType = "blob"
	TreeObject	ObjectType = "tree"
	CommitObject	ObjectType = "commit"
)

type GutObject struct {
	Type     ObjectType
	Size     int64
	Content  []byte
	HashSum  string //used for the sha256
}

func (o GutObject) String() string {
	return fmt.Sprintf("[%s] %s (%d bytes)", o.Type, o.HashSum, o.Size)
}
