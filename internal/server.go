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
	ErrAddClient        = errors.New("unable to add client")
)

type slot int

const (
	available slot = iota
	taken
)

type client struct {
	id   uint8
	conn net.Conn
}

func newClient(id uint8, conn net.Conn) client {
	return client{id: id, conn: conn}
}

type Server struct {
	clients []client
	slots   []slot
	ErrChan chan error
	Done    bool
}

func NewServer() Server {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	errChan := make(chan error)

	slots := make([]slot, MAX_CLIENTS)
	for i := 0; i < MAX_CLIENTS; i++ {
		slots[i] = available
	}
	return Server{slots: slots, ErrChan: errChan, Done: false}
}

func (srv *Server) Run() {
	listener, err := net.Listen("tcp", "0:7654")
	if err != nil {
		srv.ErrChan <- fmt.Errorf("%s: %s\n", ErrStartServer, err)
	}

	running := true
	for running {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("%s: %s\n", ErrAcceptConnection, err)
		}
		log.Infof("received a new connection: %v %v\n", conn, conn.RemoteAddr())

		client, err := srv.addClient(conn)
		if err != nil {
			log.Errorf("%s: %s\n", ErrAddClient, err)
			conn.Write([]byte("server is unable to add client add this time\n"))
			conn.Close()
		}

		// respond to client
		client.conn.Write([]byte(fmt.Sprintf("welcome! Your client id is %d\n", client.id)))
	}

	listener.Close()
	srv.Done = true
}

func (srv *Server) Stop() {
	// TODO goodbye message to clients?
}

// TODO if we hit max clients then we should check for dead clients
func (srv *Server) addClient(conn net.Conn) (client, error) {
	for idx, slot := range srv.slots {
		if slot == available {
			client := newClient(uint8(idx), conn)

			srv.clients = append(srv.clients, client)
			srv.slots[idx] = taken

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
	srv.slots[removeIdx] = available

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
