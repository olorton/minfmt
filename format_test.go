package main

import (
	"errors"
	"reflect"
	"testing"
)

func compareCleanBufferOutput(t *testing.T, expected []byte, input []byte) bool {
	found, err := CleanBuffer(input)
	if err != nil {
		t.Errorf("Error calling CleanBuffer: %v", err)
		return false
	}
	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected result to be\n%v, got:\n%v", expected, found)
	}

	return true
}

func expectedCleanBufferError(t *testing.T, input []byte, expected_error error) bool {
	_, err := CleanBuffer(input)
	if err == nil {
		t.Fatal("Expected an error but got nil")
	}

	// Check for the specific error
	if !errors.Is(err, expected_error) {
		t.Errorf("Expected error to be %v, got: %v", expected_error, err)
	}

	return true
}

func TestCleanBuffer_EndsWithNewLineWhenMissing(t *testing.T) {
	input := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
	expected := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_EndsWithASingleNewLine(t *testing.T) {
	input := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine}
	expected := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_TrailingSpaceRemoved(t *testing.T) {
	input := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteSpace}
	expected := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}
func TestCleanBuffer_TrailingTabRemoved(t *testing.T) {
	input := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteTab}
	expected := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_FourTrailingSpacesRemoved(t *testing.T) {
	input := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteSpace, ByteSpace, ByteSpace, ByteSpace}
	expected := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_MultiLineTwoFixes(t *testing.T) {
	input := []byte{'g', 'o', ByteSpace, ByteNewLine, 'l', 'a', ByteSpace, ByteSpace, ByteNewLine, 'n', 'g', ByteSpace}
	expected := []byte{'g', 'o', ByteNewLine, 'l', 'a', ByteNewLine, 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_OnlyRemoveMultipleNewLinesIfAtEndOfFile(t *testing.T) {
	input := []byte{'g', 'o', ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine, 'l', 'a', ByteNewLine, ByteNewLine, ByteNewLine, 'n', 'g', ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine}
	expected := []byte{'g', 'o', ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine, 'l', 'a', ByteNewLine, ByteNewLine, ByteNewLine, 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_RejectFileIfNullByteFound(t *testing.T) {
	input := []byte{'g', 'o', ByteSpace, ByteNewLine, 'l', NullByte, 'a', ByteSpace, ByteSpace, ByteNewLine, 'n', 'g', ByteSpace}
	expectedCleanBufferError(t, input, ErrNullByte)
}
