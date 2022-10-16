package main

import (
	"fmt"
	"github/rolandvarga/mds/internal"
)

func main() {
	srv := internal.NewServer()

	fmt.Printf("got a server: %v\n", srv)
}
