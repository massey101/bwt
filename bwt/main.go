package main

import (
	"bufio"
	"fmt"
	"os"

	"git.neds.sh/jack.massey/bwt/bwtlib"
)

// DefaultBlockSize is the default block size to use for each bwt block.
const DefaultBlockSize = 1 * 1024

func main() {
	writer := bufio.NewWriter(os.Stdout)
	err := bwtlib.BWTStream(bufio.NewReader(os.Stdin), writer, DefaultBlockSize)
	if err != nil {
		fmt.Printf("Error! %v\n", err.Error())
		os.Exit(1)
	}
	err = writer.Flush()
	if err != nil {
		fmt.Printf("Error! %v\n", err.Error())
		os.Exit(1)
	}
}
