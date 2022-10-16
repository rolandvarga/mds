package internal

import "fmt"

const MAX_CLIENTS = 256

type client struct {
	id uint8
}

func newClient(id uint8) client {
	return client{id: id}
}

type Server struct {
	clients []client
	slots   []bool
}

func NewServer() Server {
	slots := make([]bool, MAX_CLIENTS)
	for i := 0; i < MAX_CLIENTS; i++ {
		// for i := 0; i < 256; i++ {
		slots[i] = true
	}
	return Server{slots: slots}
}

func (srv *Server) Run() {
	// TODO all the logic of handling requests should go here
}

func (srv *Server) Stop() {
	// TODO goodbye message to clients?
}

// TODO if we hit max clients then we should check for dead clients
func (srv *Server) addClient() (client, error) {
	for idx, slot := range srv.slots {
		if slot == true {
			client := newClient(uint8(idx))

			srv.clients = append(srv.clients, client)
			srv.slots[idx] = false

			return client, nil
		}
	}
	return client{}, fmt.Errorf("unable to add new client")
}

// BUG there's a bug when removing non existent clients
func (srv *Server) removeClient(client uint8) error {
	found := false
	removeIdx := 0

	for i, c := range srv.clients {
		if c.id == client {
			found = true
			removeIdx = i
			break
		}
	}

	if !found {
		return fmt.Errorf("couldn't find client with ID '%d' in client list", client)
	}

	srv.clients = append(srv.clients[:removeIdx], srv.clients[removeIdx+1:]...)
	srv.slots[removeIdx] = true

	return nil
}

func (srv *Server) listClientIDsExcept(id uint8) []uint8 {
	out := []uint8{}
	for _, client := range srv.clients {
		if id != client.id {
			out = append(out, client.id)
		}
	}
	return out
}

func (srv *Server) messageClients(message string, clients []client) error {
	return nil
}
