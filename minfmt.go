package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: minfmt <file_or_directory_path>")
		return
	}

	path := os.Args[1]
	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Could not find path: %v\n", err)
		return
	}

	if info.IsDir() {
		err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				FormatFile(filePath)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error walking the path: %v\n", err)
		}
	} else {
		// Format the single file
		FormatFile(path)
	}
}
