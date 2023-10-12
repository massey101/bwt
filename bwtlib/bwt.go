package bwtlib

import (
	"encoding/binary"
	"errors"
	"io"

	"golang.org/x/exp/slices"
)

// BWT will compress a byte array and output the full bytearray.
func BWT(input []byte) ([]byte, error) {
	for i := range input {
		if input[i] == 0x02 || input[i] == 0x03 {
			return nil, errors.New("Found EOF character in input")
		}
	}

	input = append([]byte{0x02}, input...)
	input = append(input, 0x03)

	table := make([]int, 0, len(input)+2)
	for i := range input {
		table = append(table, i)
	}

	slices.SortFunc(
		table,
		func(a, b int) int {
			for i := range input {
				if input[(a+i)%len(input)] > input[(b+i)%len(input)] {
					return 1
				}
				if input[(a+i)%len(input)] < input[(b+i)%len(input)] {
					return -1
				}
			}

			return 0
		},
	)

	output := make([]byte, 0, len(input)+2)
	for _, index := range table {
		wrappedIndex := index - 1
		if wrappedIndex < 0 {
			wrappedIndex = len(input) + wrappedIndex
		}
		output = append(output, input[wrappedIndex])
	}

	return output, nil
}

func encodeBlockSize(blockSize int) []byte {
	blockSizeEncoded := make([]byte, 4)
	binary.LittleEndian.PutUint32(blockSizeEncoded, uint32(blockSize))

	return blockSizeEncoded
}

// BWTStream performs a BWT operation on a byte stream.
func BWTStream(input io.Reader, output io.Writer, blockSize int) error {
	block := make([]byte, blockSize)
	for {
		readN, readErr := io.ReadAtLeast(input, block, blockSize)
		if readN == 0 {
			if readErr != nil {
				if errors.Is(readErr, io.ErrUnexpectedEOF) || errors.Is(readErr, io.EOF) {
					break
				}

				return readErr
			}
		}

		bwtBlock, err := BWT(block[:readN])
		if err != nil {
			return err
		}

		encodedBlockSize := encodeBlockSize(len(bwtBlock))

		n, writeErr := output.Write(encodedBlockSize)
		if writeErr != nil {
			return err
		}

		if n == 0 {
			break
		}

		n, writeErr = output.Write(bwtBlock)
		if err != nil {
			return err
		}

		if n == 0 {
			break
		}

		if readErr != nil {
			if errors.Is(readErr, io.ErrUnexpectedEOF) || errors.Is(readErr, io.EOF) {
				break
			}

			return readErr
		}
	}

	return nil
}
