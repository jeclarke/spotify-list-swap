package main

import (
	"fmt"
	"os"

	"james-clarke.co.uk/listswap"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: listswapcmd <playlist-id>")
		os.Exit(1)
	}
	listswap.Run(os.Args[1])
}
