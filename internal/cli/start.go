package cli

import (
	"fmt"
	"gut/internal/store"
)

func HandleStart() {
	fs := store.NewFileStore(".")

	err := fs.InitLayout()
	if err != nil {
		fmt.Printf("Failed to start gut: %v\n", err)
		return
	}

	fmt.Println("Gut project started successfully.")
}
