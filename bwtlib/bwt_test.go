package bwtlib

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StreamFunc func(io.Reader, io.Writer, int) error
type IStreamFunc func(io.Reader, io.Writer) error

func runStreamTest(t *testing.T, f StreamFunc, input, expected []byte, blockSize int) {
	output := make([]byte, 0, len(expected))
	inputReader := bytes.NewReader(input)
	outputWriter := bytes.NewBuffer(output)
	err := f(inputReader, outputWriter, blockSize)
	assert.NoError(t, err)
	assert.Equal(t, expected, outputWriter.Bytes())
}

func runIStreamTest(t *testing.T, f IStreamFunc, input, expected []byte) {
	output := make([]byte, 0, len(expected))
	inputReader := bytes.NewReader(input)
	outputWriter := bytes.NewBuffer(output)
	err := f(inputReader, outputWriter)
	assert.NoError(t, err)
	assert.Equal(t, expected, outputWriter.Bytes())
}

func TestBWTBasic(t *testing.T) {
	input := []byte("BANANA")
	expected := []byte("\x03ANNB\x02AA")

	output, err := BWT(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
}

func TestIBWTBasic(t *testing.T) {
	input := []byte("\x03ANNB\x02AA")
	expected := []byte("BANANA")

	output, err := IBWT(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
}

func TestBWTLonger(t *testing.T) {
	input := []byte("SIX.MIXED.PIXIES.SIFT.SIXTY.PIXIE.DUST.BOXES")
	expected := []byte("\x03STEXYDST.E.IXXIIXXSSMPPS.B..EE.\x02.USFXDIIOIIIT")

	output, err := BWT(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
}

func TestIBWTLonger(t *testing.T) {
	input := []byte("\x03STEXYDST.E.IXXIIXXSSMPPS.B..EE.\x02.USFXDIIOIIIT")
	expected := []byte("SIX.MIXED.PIXIES.SIFT.SIXTY.PIXIE.DUST.BOXES")

	output, err := IBWT(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
}

func TestBWTEmpty(t *testing.T) {
	input := []byte("")
	expected := []byte("\x03\x02")

	output, err := BWT(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
}

func TestIBWTEmpty(t *testing.T) {
	input := []byte("\x03\x02")
	expected := []byte("")

	output, err := IBWT(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
}

func TestBWTStreamBasic(t *testing.T) {
	input := []byte("BANANA")
	expected := []byte("\x08\x00\x00\x00\x03ANNB\x02AA")

	runStreamTest(t, BWTStream, input, expected, 256)
}

func TestIBWTStreamBasic(t *testing.T) {
	input := []byte("\x08\x00\x00\x00\x03ANNB\x02AA")
	expected := []byte("BANANA")

	runIStreamTest(t, IBWTStream, input, expected)
}

func TestBWTStreamBlockSize(t *testing.T) {
	input := []byte("A" + strings.Repeat("B", 256-2) + "C")
	expected := []byte("\x02\x01\x00\x00\x03C\x02A" + strings.Repeat("B", 256-2))

	runStreamTest(t, BWTStream, input, expected, 256)
}

func TestIBWTStreamBlockSize(t *testing.T) {
	input := []byte("\x02\x01\x00\x00\x03C\x02A" + strings.Repeat("B", 256-2))
	expected := []byte("A" + strings.Repeat("B", 256-2) + "C")

	runIStreamTest(t, IBWTStream, input, expected)
}

func TestBWTStreamTwoBlocks(t *testing.T) {
	input := []byte("A" + strings.Repeat("B", 256-1) + "C")
	expected := []byte(
		"\x02\x01\x00\x00\x03B\x02B" + strings.Repeat("B", 256-3) + "A" + "\x03\x00\x00\x00\x03C\x02",
	)

	runStreamTest(t, BWTStream, input, expected, 256)
}

func TestIBWTStreamTwoBlocks(t *testing.T) {
	input := []byte(
		"\x02\x01\x00\x00\x03B\x02B" + strings.Repeat("B", 256-3) + "A" + "\x03\x00\x00\x00\x03C\x02",
	)
	expected := []byte("A" + strings.Repeat("B", 256-1) + "C")

	runIStreamTest(t, IBWTStream, input, expected)
}
