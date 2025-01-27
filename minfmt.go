package main

import (
	"fmt"
	"os"
	"path/filepath"
)

/*
TODOs
- When writing the file, we should not change permissions! (maybe creating a tmp file in the same location is not great, could use /tmp instead)
- Windows support
- Although this should be a silent cli tool, add an option to give a verdict on how many files have changed, and list them.
- Add ability to use in --lint mode, so that no file changes are made, and could be used in CI
- When walking a directory, silently ignore binary files
- Find the best way to ignore things like .git dirs. Should this only change files that are in git, or, behave differently?
- Needs to work with symbolic links; I'm fairly sure Walk (below) won't work currently
- Add support for a more comprehensive list of line break bytes (0x0A 0x0B 0x0D 0x0E 0x85)
- Add support for a more comprehensive list of whitespace bytes (0x20 0xA0)
*/

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
