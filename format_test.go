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
	if len(expected) == 0 && len(input) == 0 {
		return true
	}

	// TODO inspect a bool returned, e.g. can_save, that confirms the buffer has/has-not changed
	
	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected result to be\n%v, got:\n%v", expected, found)
	}

	return true
}

func expectedCleanBufferError(t *testing.T, input []byte, expected_error error) bool {
	result, err := CleanBuffer(input)
	if err == nil {
		t.Fatal("Expected an error but got nil")
	}

	if result != nil {
		t.Errorf("Expected result to be nil, found: %v", result)
	}
	
	// TODO inspect a bool returned, e.g. can_save, that should be false when erroring

	// Check for the specific error
	if !errors.Is(err, expected_error) {
		t.Errorf("Expected error to be %v, got: %v", expected_error, err)
	}

	return true
}

func TestCleanBuffer_MissingNewLine(t *testing.T) {
	input := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
	expected := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_TooManyNewLines(t *testing.T) {
	input := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine}
	expected := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_MissingNewLineAndHasATraillingSpace(t *testing.T) {
	input := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteSpace}
	expected := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_MissingNewLineAndHasATraillingTab(t *testing.T) {
	input := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteTab}
	expected := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_MissingNewLineAndHasTraillingTabAndSpace(t *testing.T) {
	input := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteTab, ByteSpace}
	expected := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_MissingNewLineAndFourTraillingSpacesRemoved(t *testing.T) {
	input := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteSpace, ByteSpace, ByteSpace, ByteSpace}
	expected := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_TwoLinesWithTraillingSpaces(t *testing.T) {
	input := []byte{'g', 'o', ByteSpace, ByteNewLine, 'l', 'a', ByteSpace, ByteSpace, ByteNewLine, 'n', 'g', ByteSpace}
	expected := []byte{'g', 'o', ByteNewLine, 'l', 'a', ByteNewLine, 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_OnlyRemoveMultipleNewLinesIfAtEndOfFile(t *testing.T) {
	input := []byte{'g', 'o', ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine, 'l', 'a', ByteNewLine, ByteNewLine, ByteNewLine, 'n', 'g', ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine}
	expected := []byte{'g', 'o', ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine, ByteNewLine, 'l', 'a', ByteNewLine, ByteNewLine, ByteNewLine, 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_NullByteFound(t *testing.T) {
	input := []byte{'g', 'o', ByteSpace, ByteNewLine, 'l', NullByte, 'a', ByteSpace, ByteSpace, ByteNewLine, 'n', 'g', ByteSpace}
	expectedCleanBufferError(t, input, ErrNullByte)
}

func TestCleanBuffer_MultipleNullBytes(t *testing.T) {
	input := []byte{'g', 'o', NullByte, 'l', 'a', NullByte, 'n', 'g'}
	expectedCleanBufferError(t, input, ErrNullByte)
}

func TestCleanBuffer_BufferNeedsNoChanges(t *testing.T) {
	input := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine}
	expected := []byte{'g', 'o', 'l', 'a', 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_MultiLineBufferNeedsNoChanges(t *testing.T) {
	input := []byte{'g', 'o', ByteNewLine, 'l', 'a', 'n', 'g', ByteNewLine}
	expected := []byte{'g', 'o', ByteNewLine, 'l', 'a', 'n', 'g', ByteNewLine}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_EmptyInput(t *testing.T) {
	input := []byte{}
	expected := []byte{}
	compareCleanBufferOutput(t, expected, input)
}

func TestCleanBuffer_OnlyWhitespace(t *testing.T) {
	input := []byte{ByteSpace, ByteSpace, ByteTab}
	expected := []byte{}
	compareCleanBufferOutput(t, expected, input)
}
