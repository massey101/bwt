package compresslib

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

func encodeRunLength(runLength int, runLengthBytes int) []byte {
	runLengthBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(runLengthBuffer, uint64(runLength))

	return runLengthBuffer[:runLengthBytes]
}

func getMaxRunLength(minRun int, runLengthBytes int) int {
	return 1<<(8*runLengthBytes) + minRun
}

func encodeRun(runLength int, runChar byte, minRun int, runLengthBytes int) []byte {
	if runLength < minRun {
		return bytes.Repeat([]byte{runChar}, runLength)
	}

	return append(bytes.Repeat([]byte{runChar}, minRun), encodeRunLength(runLength-minRun, runLengthBytes)...)
}

// Compress will compress a byte stream and output to another byte stream.
// Returns the number of bytes read, the number of bytes output and any errors.
func Compress(input io.ByteReader, output io.Writer, minRun int, runLengthBytes int) (int, int, error) {
	maxRunLength := getMaxRunLength(minRun, runLengthBytes)

	inBytes := 0
	outBytes := 0
	runChar := byte(0)
	runLength := 0
	for {
		readByte, err := input.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return inBytes, outBytes, err
		}
		inBytes++

		if readByte != runChar || runLength >= maxRunLength-1 {
			if runLength > 0 {
				encodedRun := encodeRun(runLength, runChar, minRun, runLengthBytes)
				n, err := output.Write(encodedRun)
				outBytes += n
				if err != nil {
					return inBytes, outBytes, err
				}
			}
			runLength = 0
			runChar = readByte
		}

		runLength++
	}

	if runLength > 0 {
		encodedRun := encodeRun(runLength, runChar, minRun, runLengthBytes)
		n, err := output.Write(encodedRun)
		outBytes += n
		if err != nil {
			return inBytes, outBytes, err
		}
	}

	return inBytes, outBytes, nil
}

func readRunLength(input io.ByteReader, runLengthBytes int) (int, error) {
	runLengthBinary := make([]byte, 8)

	for i := 0; i < runLengthBytes; i++ {
		runLengthBinaryByte, err := input.ReadByte()
		if err != nil {
			return 0, err
		}

		runLengthBinary[i] = runLengthBinaryByte
	}

	runLength := binary.LittleEndian.Uint64(runLengthBinary)

	return int(runLength), nil
}

// Decompress will decompress the input stream and write to the output.
func Decompress(input io.ByteReader, output io.Writer, minRun int, runLengthBytes int) (int, int, error) {

	inBytes := 0
	outBytes := 0
	runChar := byte(0x00)
	runLength := 0

	for {
		readChar, err := input.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return inBytes, outBytes, err
		}
		inBytes++
		n, err := output.Write([]byte{readChar})
		outBytes += n
		if err != nil {
			return inBytes, outBytes, err
		}

		if readChar != runChar {
			runLength = 1
			runChar = readChar
			continue
		}

		runLength++
		if runLength >= minRun {
			repeatedRunLength, err := readRunLength(input, runLengthBytes)
			inBytes += runLengthBytes
			if err != nil {
				if errors.Is(err, io.EOF) {
					return inBytes, outBytes, errors.New("expected encoded run length, got eof")
				}

				return inBytes, outBytes, err
			}

			n, err := output.Write(bytes.Repeat([]byte{runChar}, repeatedRunLength))
			outBytes += n
			if err != nil {
				return inBytes, outBytes, err
			}
			runLength = 0
			runChar = 0x00
		}
	}

	return inBytes, outBytes, nil
}
