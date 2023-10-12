package bwtlib

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sort"

	"golang.org/x/exp/slices"
)

// IBWT will compress a byte array and output the full bytearray.
func IBWT(bwt []byte) ([]byte, error) {
	if !slices.Contains(bwt, 0x02) || !slices.Contains(bwt, 0x03) {
		return nil, errors.New("Need EOF character in input")
	}

	sortedBwt := append([]byte{}, bwt...)
	sort.Slice(
		sortedBwt,
		func(i, j int) bool {
			return sortedBwt[i] < sortedBwt[j]
		},
	)
	lShift := make([]int, len(bwt))
	lShiftUsed := make([]bool, len(bwt))

	for i := range sortedBwt {
		for j := range bwt {
			if lShiftUsed[j] == true {
				continue
			}

			if sortedBwt[i] == bwt[j] {
				lShift[i] = j
				lShiftUsed[j] = true
				break
			}
		}
	}

	original := make([]byte, 0, len(bwt))

	x := 0
	for x = range bwt {
		if bwt[x] == 0x03 {
			break
		}
	}

	for range bwt {
		x = lShift[x]
		original = append(original, bwt[x])
	}

	return original[1 : len(original)-1], nil
}

// readBlockSize will read and decode the block size from the input buffer. A
// block size of 0 indicates that the buffer has finished.
func readBlockSize(input io.Reader) (int, error) {
	blockSizeBuffer := make([]byte, 4)
	n, err := io.ReadAtLeast(input, blockSizeBuffer, 4)

	if err != nil {
		if errors.Is(err, io.EOF) {
			return 0, nil
		}

		if errors.Is(err, io.ErrUnexpectedEOF) {
			return 0, errors.New("malformed block size")
		}

		return 0, err
	}

	if n == 0 {
		return 0, nil
	}

	if n != 4 {
		return 0, errors.New("malformed block size")
	}

	return int(binary.LittleEndian.Uint32(blockSizeBuffer)), nil
}

// IBWTStream performs a BWT operation on a byte stream.
func IBWTStream(input io.Reader, output io.Writer) error {
	defaultBlockSize := 32 * 1024
	block := make([]byte, 0, defaultBlockSize)
	for {
		blockSize, err := readBlockSize(input)
		if err != nil {
			return err
		}

		if blockSize == 0 {
			break
		}

		block = block[:blockSize]
		n, err := io.ReadAtLeast(input, block, blockSize)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				return fmt.Errorf("malformed block of %v bytes, expected %v", n, blockSize)
			}

			return err
		}

		if n != blockSize {
			return fmt.Errorf("malformed block of %v bytes, expected %v", n, blockSize)
		}

		originalBlock, err := IBWT(block)
		if err != nil {
			return err
		}

		n, err = output.Write(originalBlock)
		if err != nil {
			return err
		}

		if n == 0 {
			return fmt.Errorf("failed to write to output")
		}
	}

	return nil
}
