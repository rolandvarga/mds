package test

import (
	"github/rolandvarga/mds/internal"
	"testing"
)

const SERVER_PORT = 7654

func TestServer(t *testing.T) {
	srv := internal.NewServer()

	// TODO this won't block though
	srv.Run()

	go func() {
		for !srv.Done {
			select {
			case err := <-srv.ErrChan:
				t.Errorf("server produced an error: %v\n", err)
			}
		}
	}()

	// for internal.MAX_CLIENTS: client.connectToServeR()

	// TODO return error from server if max_clients & trying to add new
}
