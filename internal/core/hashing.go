package core

import (
	"crypto/sha256"
	"fmt"
	"gut/pkg/models"
)

func HashObject(data []byte, objType models.ObjectType) models.GutObject {
	header := fmt.Sprintf("%s %d\x00", objType, len(data))

	fullContent := append([]byte(header), data...)

	hash := sha256.Sum256(fullContent)
	hashString := fmt.Sprintf("%x", hash)

	return models.GutObject {
		Type: objType,
		Size: int64(len(data)),
		Content: fullContent,
		HashSum: hashString,
	}
}
