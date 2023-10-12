package mtflib

import (
	"errors"
	"io"
)

func applyAndUpdateTransform(transform []byte, x byte) ([]byte, byte) {
	var i int
	for i = range transform {
		if transform[i] == x {
			break
		}
	}

	newTransform := append([]byte{x}, transform[:i]...)
	newTransform = append(newTransform, transform[i+1:]...)

	return newTransform, byte(i & 0xff)
}

// MTF will apply the MoveToFront transform to an io stream.
func MTF(input io.ByteReader, output io.ByteWriter) error {
	transform := make([]byte, 0xff)

	for i := range transform {
		transform[i] = byte(i & 0xff)
	}

	for {
		x, err := input.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}

		transform, x = applyAndUpdateTransform(transform, x)
		if err = output.WriteByte(x); err != nil {
			return err
		}
	}

	return nil
}
