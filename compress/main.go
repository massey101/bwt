package main

import (
	"bufio"
	"fmt"
	"os"

	"git.neds.sh/jack.massey/bwt/compresslib"
)

func main() {
	writer := bufio.NewWriter(os.Stdout)
	_, _, err := compresslib.Compress(bufio.NewReader(os.Stdin), writer, 4, 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error! %v\n", err.Error())
		os.Exit(1)
	}
	err = writer.Flush()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error! %v\n", err.Error())
		os.Exit(1)
	}
}
