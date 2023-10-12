package main

import (
	"bufio"
	"fmt"
	"os"

	"git.neds.sh/jack.massey/bwt/compresslib"
)

func main() {
	writer := bufio.NewWriter(os.Stdout)
	err := compresslib.Decompress(bufio.NewReader(os.Stdin), writer, 1)
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
