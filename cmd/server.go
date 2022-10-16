package main

import (
	"github/rolandvarga/mds/internal"
)

func main() {
	srv := internal.NewServer()

	srv.Run()
}
