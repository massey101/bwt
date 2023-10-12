package main

import (
	"bufio"
	"fmt"
	"os"

	"git.neds.sh/jack.massey/bwt/mtflib"
)

func main() {
	writer := bufio.NewWriter(os.Stdout)
	err := mtflib.IMTF(bufio.NewReader(os.Stdin), writer)
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
