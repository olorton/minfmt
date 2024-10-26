package formatter

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// TODO list
// - Read file/dir as argument
// - Iterate over all sub dirs and files if a directory
// - If null byte, this is a binary file, so cancel

const NullByte = byte(0x00)
const ByteNewLine = byte(0x0A)
const CharTab = string(0x09)
const CharSpace = string(0x20)

func Format() {
	src_path := getFullPath()
	buffer_output := CleanBuffer(getInputBuffer(src_path))
	writeFile(buffer_output, src_path)
}

func CleanBuffer(buffer_input []byte) (buffer_output []byte) {
	var b byte
	for {
		line := []byte{}
		// loop over the buffer_output to create a line, stop at ByteNewLine
		for {
			if len(buffer_input) == 0 {
				break
			}
			b, buffer_input = buffer_input[0], buffer_input[1:]
			if b == ByteNewLine {
				break
			}
			if b == NullByte {
				os.Exit(0)
			}
			line = append(line, b)
		}

		line = bytes.TrimRight(line, CharTab+CharSpace)
		line = append(line, ByteNewLine)
		buffer_output = append(buffer_output, line...)
		if len(buffer_input) == 0 {
			break
		}
	}
	for buffer_output[len(buffer_output)-1] == ByteNewLine &&
		buffer_output[len(buffer_output)-2] == ByteNewLine {
		buffer_output = buffer_output[:len(buffer_output)-1]
	}
	return
}
func writeFile(buffer_output []byte, src_path string) {
	// Open the file for writing
	file_tmp, err := os.Create(src_path + "~")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file_tmp.Close()

	// Write to temp file
	writer := bufio.NewWriter(file_tmp)
	err = nil
	_, err = writer.Write(buffer_output)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = nil
	err = writer.Flush()
	if err != nil {
		fmt.Println(err)
		return
	}
	file_tmp.Close()

	// Replace file with the temp file
	err = nil
	err = os.Rename(src_path+"~", src_path)
	if err != nil {
		fmt.Println(err)
	}
}

func getInputBuffer(src_path string) []byte {
	// Read the file to an input buffer
	buffer_input, err := os.ReadFile(src_path)
	if err != nil {
		fmt.Print(err)
	}

	return buffer_input
}

func getFullPath() string {
	filename := os.Args[1]
	regexForTilde := regexp.MustCompile(`^~`)
	m, err := regexp.MatchString("^~", filename)
	if m {
		homedir, err := os.UserHomeDir()
		if err != nil {
			fmt.Print(err)
		}
		filename = regexForTilde.ReplaceAllLiteralString(filename, homedir)
	}
	full_filepath, err := filepath.Abs(filename)
	if err != nil {
		fmt.Print(err)
	}
	return full_filepath
}
