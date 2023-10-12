package compresslib

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

func encodeRunLength(runLength int, runLengthBytes int) []byte {
	runLengthBuffer := make([]byte, 8)
	binary.BigEndian.PutUint64(runLengthBuffer, uint64(runLength))

	return runLengthBuffer[8-runLengthBytes:]
}

func getMaxRunLength(runLengthBytes int) int {
	return 1 << (7 + (8 * (runLengthBytes - 1)))
}

func encodeRun(runLength int, runChar byte, runLengthBytes int) []byte {
	maxRunLength := getMaxRunLength(runLengthBytes)

	if runLength <= runLengthBytes+1 {
		return bytes.Repeat([]byte{runChar}, runLength)
	}

	runLength += maxRunLength
	return append(encodeRunLength(runLength, runLengthBytes), runChar)
}

// Compress will compress a byte stream and output to another byte stream.
func Compress(input io.ByteReader, output io.Writer, runLengthBytes int) error {
	maxRunLength := getMaxRunLength(runLengthBytes)

	runChar := byte(0)
	runLength := 0
	for {
		readByte, err := input.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}

		if readByte != runChar || runLength >= maxRunLength-1 {
			if runLength > 0 {
				encodedRun := encodeRun(runLength, runChar, runLengthBytes)
				_, err := output.Write(encodedRun)
				if err != nil {
					return err
				}
			}
			runLength = 0
			runChar = readByte
		}

		runLength++
	}

	if runLength > 0 {
		encodedRun := encodeRun(runLength, runChar, runLengthBytes)
		_, err := output.Write(encodedRun)
		if err != nil {
			return err
		}
	}

	return nil
}
