package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var ErrNullByte = errors.New("Null byte found, cannot format a binary file")

const NullByte = byte(0x00)
const ByteNewLine = byte(0x0A)
const ByteTab = byte(0x09)
const ByteSpace = byte(0x20)
const CharTab = string(rune(0x09))
const CharSpace = string(rune(0x20))

// TODO refactor this, handle all sorts of errors.
func FormatFile(path string) {
	src_path := getFullPath(path)
	buffer_output, err := CleanBuffer(getInputBuffer(src_path))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	writeFile(buffer_output, src_path)
}

// TODO also return a bool so that when a buffer is unchanged, we don't have to write to the new file
func CleanBuffer(buffer_input []byte) ([]byte, error) {
	var b byte
	var buffer_output []byte
	if len(buffer_input) == 0 {
		return buffer_output, nil
	}
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
				return nil, ErrNullByte
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
	if len(buffer_output) > 1 {
		// Remove trailing newlines from the output buffer if there are multiple
		// consecutive newlines
		for buffer_output[len(buffer_output)-1] == ByteNewLine &&
			buffer_output[len(buffer_output)-2] == ByteNewLine {
			buffer_output = buffer_output[:len(buffer_output)-1]
		}
	}
	if len(buffer_output) == 1 {
		// We should not ever return simply a single byte. Either there is a
		// non-line-break byte followed by a new-line, or it should be empty.
		return []byte{}, nil
	}
	return buffer_output, nil
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

func getFullPath(filename string) string {
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
