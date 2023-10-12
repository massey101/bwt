package compresslib

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StreamFunc func(io.ByteReader, io.Writer, int, int) (int, int, error)

func runStreamTest(t *testing.T, f StreamFunc, input, expected []byte, minRun int, runLengthBytes int) {
	output := make([]byte, 0, len(expected))
	inputReader := bytes.NewReader(input)
	outputWriter := bytes.NewBuffer(output)
	inByte, outByte, err := f(inputReader, outputWriter, minRun, runLengthBytes)
	assert.NoError(t, err)
	assert.Equal(t, expected, outputWriter.Bytes())
	assert.Equal(t, len(input), inByte)
	assert.Equal(t, len(expected), outByte)
}

func TestCompressBasic(t *testing.T) {
	input := []byte("Test    !!")
	expected := []byte("Test  \x02!!\x00")
	runStreamTest(t, Compress, input, expected, 2, 1)
}

func TestDecompressBasic(t *testing.T) {
	input := []byte("Test  \x02!!\x00")
	expected := []byte("Test    !!")
	runStreamTest(t, Decompress, input, expected, 2, 1)
}

func TestCompressLargeRunLength(t *testing.T) {
	input := []byte("A" + strings.Repeat("B", 257) + "C")
	expected := []byte("ABB\xffC")
	runStreamTest(t, Compress, input, expected, 2, 1)
}

func TestDecompressLargeRunLength(t *testing.T) {
	input := []byte("ABB\xffC")
	expected := []byte("A" + strings.Repeat("B", 257) + "C")
	runStreamTest(t, Decompress, input, expected, 2, 1)
}

func TestCompressMaxRunLength(t *testing.T) {
	input := []byte("A" + strings.Repeat("B", 258) + "C")
	expected := []byte("ABB\xffBC")
	runStreamTest(t, Compress, input, expected, 2, 1)
}

func TestDecompressMaxRunLength(t *testing.T) {
	input := []byte("ABB\xffBC")
	expected := []byte("A" + strings.Repeat("B", 258) + "C")
	runStreamTest(t, Decompress, input, expected, 2, 1)
}

func TestCompressDoubleMaxRunLength(t *testing.T) {
	input := []byte("A" + strings.Repeat("B", 514) + "C")
	expected := []byte("ABB\xffBB\xffC")
	runStreamTest(t, Compress, input, expected, 2, 1)
}

func TestDecompressDoubleMaxRunLength(t *testing.T) {
	input := []byte("ABB\xffBB\xffC")
	expected := []byte("A" + strings.Repeat("B", 514) + "C")
	runStreamTest(t, Decompress, input, expected, 2, 1)
}

func TestCompressEmpty(t *testing.T) {
	input := []byte{}
	expected := []byte{}
	runStreamTest(t, Compress, input, expected, 2, 1)
}

func TestDecompressEmpty(t *testing.T) {
	input := []byte{}
	expected := []byte{}
	runStreamTest(t, Compress, input, expected, 2, 1)
}

func TestCompress31Basic(t *testing.T) {
	input := []byte("Test    !!!")
	expected := []byte("Test   \x01!!!\x00")
	runStreamTest(t, Compress, input, expected, 3, 1)
}

func TestDecompress31Basic31(t *testing.T) {
	input := []byte("Test   \x01!!!\x00")
	expected := []byte("Test    !!!")
	runStreamTest(t, Decompress, input, expected, 3, 1)
}

func TestCompress31LargeRunLength(t *testing.T) {
	input := []byte("A" + strings.Repeat("B", 258) + "C")
	expected := []byte("ABBB\xffC")
	runStreamTest(t, Compress, input, expected, 3, 1)
}

func TestDecompress31LargeRunLength(t *testing.T) {
	input := []byte("ABBB\xffC")
	expected := []byte("A" + strings.Repeat("B", 258) + "C")
	runStreamTest(t, Decompress, input, expected, 3, 1)
}

func TestCompress31MaxRunLength(t *testing.T) {
	input := []byte("A" + strings.Repeat("B", 259) + "C")
	expected := []byte("ABBB\xffBC")
	runStreamTest(t, Compress, input, expected, 3, 1)
}

func TestDecompress31MaxRunLength(t *testing.T) {
	input := []byte("ABBB\xffBC")
	expected := []byte("A" + strings.Repeat("B", 259) + "C")
	runStreamTest(t, Decompress, input, expected, 3, 1)
}

func TestCompress31DoubleMaxRunLength(t *testing.T) {
	input := []byte("A" + strings.Repeat("B", 516) + "C")
	expected := []byte("ABBB\xffBBB\xffC")
	runStreamTest(t, Compress, input, expected, 3, 1)
}

func TestDecompress31DoubleMaxRunLength(t *testing.T) {
	input := []byte("ABBB\xffBBB\xffC")
	expected := []byte("A" + strings.Repeat("B", 516) + "C")
	runStreamTest(t, Decompress, input, expected, 3, 1)
}

func TestCompress31Empty(t *testing.T) {
	input := []byte{}
	expected := []byte{}
	runStreamTest(t, Compress, input, expected, 3, 1)
}

func TestDecompress31Empty(t *testing.T) {
	input := []byte{}
	expected := []byte{}
	runStreamTest(t, Compress, input, expected, 3, 1)
}

func TestCompress22Basic(t *testing.T) {
	input := []byte("Test    !!")
	expected := []byte("Test  \x02\x00!!\x00\x00")
	runStreamTest(t, Compress, input, expected, 2, 2)
}

func TestDecompress22Basic22(t *testing.T) {
	input := []byte("Test  \x02\x00!!\x00\x00")
	expected := []byte("Test    !!")
	runStreamTest(t, Decompress, input, expected, 2, 2)
}

func TestCompress22LargeRunLength(t *testing.T) {
	input := []byte("ABB" + strings.Repeat("B", 1<<16-1) + "C")
	expected := []byte("ABB\xff\xffC")
	runStreamTest(t, Compress, input, expected, 2, 2)
}

func TestDecompress22LargeRunLength(t *testing.T) {
	input := []byte("ABB\xff\xffC")
	expected := []byte("ABB" + strings.Repeat("B", 1<<16-1) + "C")
	runStreamTest(t, Decompress, input, expected, 2, 2)
}

func TestCompress22MaxRunLength(t *testing.T) {
	input := []byte("ABB" + strings.Repeat("B", 1<<16-1) + "BC")
	expected := []byte("ABB\xff\xffBC")
	runStreamTest(t, Compress, input, expected, 2, 2)
}

func TestDecompress22MaxRunLength(t *testing.T) {
	input := []byte("ABB\xff\xffBC")
	expected := []byte("ABB" + strings.Repeat("B", 1<<16-1) + "BC")
	runStreamTest(t, Decompress, input, expected, 2, 2)
}

func TestCompress22Empty(t *testing.T) {
	input := []byte{}
	expected := []byte{}
	runStreamTest(t, Compress, input, expected, 2, 2)
}

func TestDecompress22Empty(t *testing.T) {
	input := []byte{}
	expected := []byte{}
	runStreamTest(t, Compress, input, expected, 2, 2)
}
