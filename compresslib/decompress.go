package compresslib

import (
	"encoding/binary"
	"errors"
	"io"
)

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
func Decompress(input io.ByteReader, output io.ByteWriter, runLengthBytes int) error {
	for {
		runLength, err := readRunLength(input, runLengthBytes)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}
		readByte, err := input.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return errors.New("Unexpected EOF")
			}

			return err
		}
		for i := 0; i < runLength; i++ {
			output.WriteByte(readByte)
		}
	}

	return nil
}
