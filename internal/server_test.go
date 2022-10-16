package internal

import (
	"reflect"
	"testing"
)

func TestRemoveClient(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		client     uint8
		clientList []client
		want       []client
	}{
		"Remove client Start": {
			client: 0,
			clientList: []client{
				client{id: 0},
				client{id: 1},
				client{id: 2},
			},
			want: []client{
				client{id: 1},
				client{id: 2},
			},
		},
		"Remove client Middle": {
			client: 1,
			clientList: []client{
				client{id: 0},
				client{id: 1},
				client{id: 2},
			},
			want: []client{
				client{id: 0},
				client{id: 2},
			},
		},
		"Remove client End": {
			client: 2,
			clientList: []client{
				client{id: 0},
				client{id: 1},
				client{id: 2},
			},
			want: []client{
				client{id: 0},
				client{id: 1},
			},
		},
		"Client not found": {
			client: 0,
			clientList: []client{
				client{id: 1},
				client{id: 2},
			},
			want: []client{
				client{id: 1},
				client{id: 2},
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
				srv.slots[c.id] = false
			}

			err := srv.removeClient(c.client)
			if err != nil {
				if name == "Client not found" {
					// error is expected
				} else {
					t.Errorf("encountered an error while removing client: %v\n", err)
				}
			}

			if !reflect.DeepEqual(srv.clients, c.want) {
				t.Errorf("removing client didn't produce expected client list;\ngot: %v\nwant: %v\n", srv.clients, c.want)
			}
		})
	}
}

func TestListClientsExcept(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		client     uint8
		clientList []client
		want       []uint8
	}{
		"List clients Ok": {
			client: 0,
			clientList: []client{
				client{id: 0},
				client{id: 1},
				client{id: 2},
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
				srv.slots[c.id] = false
			}

			clients := srv.listClientIDsExcept(c.client)

			if !reflect.DeepEqual(clients, c.want) {
				t.Errorf("didn't receive expected client list;\ngot: %v\nwant: %v\n", clients, c.want)
			}
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
				srv.clients = append(srv.clients, client{id: uint8(id)})
			}

			client, err := srv.addClient()
			if err != nil && c.wantError != true {
				t.Errorf("received an unexpected error: %v\nclient: %v\nclients: %v", err, client, srv.clients)
			}
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
		wantClient  client
		wantSlot    int
	}{
		"Add client at idx '0' Ok": {
			clientCount: 10,
			emptySlots:  []uint8{0, 1},
			wantClient:  client{id: 0},
			wantSlot:    0,
		},
		"Add client at idx '1' Ok": {
			clientCount: 10,
			emptySlots:  []uint8{1, 2, 3},
			wantClient:  client{id: 1},
			wantSlot:    1,
		},
		"Add client at idx '50' Ok": {
			clientCount: 100,
			emptySlots:  []uint8{50, 91},
			wantClient:  client{id: 50},
			wantSlot:    50,
		},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv := NewServer()

			for id := 0; id <= int(c.clientCount); id++ {
				if !sliceContains(uint8(id), c.emptySlots) {
					srv.clients = append(srv.clients, client{id: uint8(id)})
					srv.slots[id] = false
				}
			}

			client, err := srv.addClient()
			if err != nil {
				t.Errorf("received an unexpected error while adding client: %v\n", err)
			}

			if client.id != c.wantClient.id {
				// if client.id != c.wantClient.id || srv.slots[c.wantSlot] != false {
				t.Errorf("received client has an unexpected ID;\ngot: %d\nwant: %d\nclients: %v", client.id, c.wantClient.id, srv.clients)
			}
		})
	}
}
