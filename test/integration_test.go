package test

import (
	"github/rolandvarga/mds/internal"
	"testing"
)

const SERVER_PORT = 7654

func TestServer(t *testing.T) {
	srv := internal.NewServer()

	err := srv.Run()
	if err != nil {
		t.Errorf("server produced an error: %v\n", err)
	}
}
