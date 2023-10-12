package mtflib

import (
	"errors"
	"io"
)

func applyAndUpdateITransform(transform []byte, i byte) ([]byte, byte) {
	x := transform[i]

	newTransform := append([]byte{x}, transform[:i]...)
	newTransform = append(newTransform, transform[i+1:]...)

	return newTransform, x
}

// IMTF will apply the MoveToFront transform to an io stream.
func IMTF(input io.ByteReader, output io.ByteWriter) error {
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

		transform, x = applyAndUpdateITransform(transform, x)
		if err = output.WriteByte(x); err != nil {
			return err
		}
	}

	return nil
}
