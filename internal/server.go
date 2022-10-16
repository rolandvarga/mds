package internal

import (
	"errors"
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

const MAX_CLIENTS = 256

var (
	ErrStartServer      = errors.New("unable to start server")
	ErrAcceptConnection = errors.New("unable to accept connection")
)

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
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	slots := make([]bool, MAX_CLIENTS)
	for i := 0; i < MAX_CLIENTS; i++ {
		slots[i] = true
	}
	return Server{slots: slots}
}

func (srv *Server) Run() error {
	listener, err := net.Listen("tcp", "0:7654")
	if err != nil {
		return fmt.Errorf("%s: %s\n", ErrStartServer, err)
	}

	running := true
	for running {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("%s: %s\n", ErrAcceptConnection, err)
		}
		log.Infof("received a new connection: %v\n", conn)
	}

	listener.Close()
	return nil
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
