package main

import (
	"github/rolandvarga/mds/internal"
)

func main() {
	srv := internal.NewServer()

	err := srv.Run()
	if err != nil {
		panic(err)
	}
}
