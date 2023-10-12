package mtflib

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StreamFunc func(io.ByteReader, io.ByteWriter) error

func runStreamTest(t *testing.T, f StreamFunc, input, expected []byte) {
	output := make([]byte, 0, len(expected))
	inputReader := bytes.NewReader(input)
	outputWriter := bytes.NewBuffer(output)
	err := f(inputReader, outputWriter)
	assert.NoError(t, err)
	assert.Equal(t, expected, outputWriter.Bytes())
}

func TestMTFStreamBasic(t *testing.T) {
	input := []byte("BANANA")
	expected := []byte("\x42\x42\x4e\x01\x01\x01")

	runStreamTest(t, MTF, input, expected)
}

func TestIMTFStreamBasic(t *testing.T) {
	input := []byte("\x42\x42\x4e\x01\x01\x01")
	expected := []byte("BANANA")

	runStreamTest(t, IMTF, input, expected)
}
