package compresslib

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StreamFunc func(io.ByteReader, io.Writer, int) error

func runStreamTest(t *testing.T, f StreamFunc, input, expected []byte, runLengthBytes int) {
	output := make([]byte, 0, len(expected))
	inputReader := bytes.NewReader(input)
	outputWriter := bytes.NewBuffer(output)
	err := f(inputReader, outputWriter, 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, outputWriter.Bytes())
}

func TestCompressBasic(t *testing.T) {
	input := []byte("Test    !!!")
	expected := []byte("Test\x84 \x83!")
	runStreamTest(t, Compress, input, expected, 1)
}

func TestDecompressBasic(t *testing.T) {
	input := []byte{
		0x01, 'T',
		0x01, 'e',
		0x01, 's',
		0x01, 't',
		0x04, ' ',
		0x03, '!',
	}
	expected := []byte("Test    !!!")
	output := make([]byte, 0, len(expected))
	inputReader := bytes.NewReader(input)
	outputWriter := bytes.NewBuffer(output)
	err := Decompress(inputReader, outputWriter, 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, outputWriter.Bytes())
}

func TestCompressLargeRunLength(t *testing.T) {
	input := []byte("A" + strings.Repeat("B", 127) + "C")
	expected := []byte("A\xffBC")
	runStreamTest(t, Compress, input, expected, 1)
}

func TestDecompressLargeRunLength(t *testing.T) {
	input := []byte{
		0x01, byte('A'),
		0xff, byte('B'),
		0x01, byte('C'),
	}

	expected := make([]byte, 0, 260)
	expected = append(expected, 'A')
	for i := 0; i < 255; i++ {
		expected = append(expected, 'B')
	}
	expected = append(expected, 'C')

	output := make([]byte, 0, len(expected))
	inputReader := bytes.NewReader(input)
	outputWriter := bytes.NewBuffer(output)
	err := Decompress(inputReader, outputWriter, 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, outputWriter.Bytes())
}

func TestCompressMaxRunLength(t *testing.T) {
	input := []byte("A" + strings.Repeat("B", 128) + "C")
	expected := []byte("A\xffBBC")
	runStreamTest(t, Compress, input, expected, 1)
}

func TestCompressDoubleMaxRunLength(t *testing.T) {
	input := []byte("A" + strings.Repeat("B", 256) + "C")
	expected := []byte("A\xffB\xffBBBC")
	runStreamTest(t, Compress, input, expected, 1)
}

func TestDecompressMaxRunLength(t *testing.T) {
	input := []byte{
		0x01, byte('A'),
		0xff, byte('B'),
		0x01, byte('B'),
		0x01, byte('C'),
	}

	expected := make([]byte, 0, 260)
	expected = append(expected, 'A')
	for i := 0; i < 256; i++ {
		expected = append(expected, 'B')
	}
	expected = append(expected, 'C')

	output := make([]byte, 0, len(expected))
	inputReader := bytes.NewReader(input)
	outputWriter := bytes.NewBuffer(output)
	err := Decompress(inputReader, outputWriter, 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, outputWriter.Bytes())
}

func TestCompressEmpty(t *testing.T) {
	input := []byte{}
	expected := []byte{}
	runStreamTest(t, Compress, input, expected, 1)
}

func TestDecompressEmpty(t *testing.T) {
	input := []byte{}
	expected := []byte{}

	output := make([]byte, 0, len(expected))
	inputReader := bytes.NewReader(input)
	outputWriter := bytes.NewBuffer(output)
	err := Decompress(inputReader, outputWriter, 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, outputWriter.Bytes())
}
