package internal

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveClient(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		client     uint8
		clientList []Client
		want       []Client
	}{
		"Remove client Start": {
			client: 0,
			clientList: []Client{
				Client{ID: 0},
				Client{ID: 1},
				Client{ID: 2},
			},
			want: []Client{
				Client{ID: 1},
				Client{ID: 2},
			},
		},
		"Remove client Middle": {
			client: 1,
			clientList: []Client{
				Client{ID: 0},
				Client{ID: 1},
				Client{ID: 2},
			},
			want: []Client{
				Client{ID: 0},
				Client{ID: 2},
			},
		},
		"Remove client End": {
			client: 2,
			clientList: []Client{
				Client{ID: 0},
				Client{ID: 1},
				Client{ID: 2},
			},
			want: []Client{
				Client{ID: 0},
				Client{ID: 1},
			},
		},
		"Client not found": {
			client: 0,
			clientList: []Client{
				Client{ID: 1},
				Client{ID: 2},
			},
			want: []Client{
				Client{ID: 1},
				Client{ID: 2},
			},
		},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			srv := NewServer()
			srv.clients = c.clientList
			for _, c := range c.clientList {
				srv.slots[c.ID] = taken
			}

			err := srv.removeClient(c.client)
			if err != nil {
				if name == "Client not found" {
					// error is expected
				} else {
					t.Errorf("encountered an error while removing client: %v\n", err)
				}
			}

			assert.Equal(t, c.want, srv.clients, "got: %v\nwant: %v", srv.clients, c.want)
		})
	}
}

func TestListClientsExcept(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		client     uint8
		clientList []Client
		want       []uint8
	}{
		"List clients Ok": {
			client: 0,
			clientList: []Client{
				Client{ID: 0},
				Client{ID: 1},
				Client{ID: 2},
			},
			want: []uint8{1, 2},
		},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv := NewServer()
			srv.clients = c.clientList
			for _, c := range c.clientList {
				srv.slots[c.ID] = taken
			}

			clients := srv.listClientIDsExcept(c.client)

			assert.Equal(t, c.want, clients, "got: %v\nwant: %v", clients, c.want)
		})
	}
}

func TestAddClientOverflow(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		clientCount uint8
		wantError   bool
	}{
		"Prevent client limit overflow": {
			clientCount: 255,
			wantError:   true,
		},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv := NewServer()

			for id := 0; id <= int(c.clientCount); id++ {
				srv.clients = append(srv.clients, Client{ID: uint8(id)})
				srv.slots[id] = taken
			}

			// create a pipe so that a connection can be passed to the client
			serverConn, clientConn := net.Pipe()

			client, err := srv.addClient(clientConn)
			if err != nil && c.wantError != true {
				t.Errorf("received an unexpected error: %v\nclient: %v\nclients: %v", err, client, srv.clients)
			}
			serverConn.Close()
			clientConn.Close()
		})
	}
}

func sliceContains(id uint8, slice []uint8) bool {
	for _, s := range slice {
		if id == s {
			return true
		}
	}
	return false
}

func TestAddClientAtExpectedIndex(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		clientCount uint8
		emptySlots  []uint8
		wantClient  Client
	}{
		"Add client at idx '0' Ok": {
			clientCount: 10,
			emptySlots:  []uint8{0, 1},
			wantClient:  Client{ID: 0},
		},
		"Add client at idx '1' Ok": {
			clientCount: 10,
			emptySlots:  []uint8{1, 2, 3},
			wantClient:  Client{ID: 1},
		},
		"Add client at idx '50' Ok": {
			clientCount: 100,
			emptySlots:  []uint8{50, 91},
			wantClient:  Client{ID: 50},
		},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv := NewServer()

			for id := 0; id <= int(c.clientCount); id++ {
				if !sliceContains(uint8(id), c.emptySlots) {
					srv.clients = append(srv.clients, Client{ID: uint8(id)})
					srv.slots[id] = taken
				}
			}

			// create a pipe so that a connection can be passed to the client
			serverConn, clientConn := net.Pipe()

			client, err := srv.addClient(clientConn)
			require.NoError(t, err)

			assert.Equal(t, c.wantClient.ID, client.ID, "got: %v\nwant: %v", client.ID, c.wantClient.ID)

			serverConn.Close()
			clientConn.Close()
		})
	}
}
