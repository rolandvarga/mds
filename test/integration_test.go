package test

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github/rolandvarga/mds/internal"
	"net"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

const SERVER_PORT = 7654

func TestServer(t *testing.T) {
	// -- --------- SETUP VARS ---------
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	serverAddr := fmt.Sprintf("localhost:%d", SERVER_PORT)

	// -- --------- SETUP SERVER ---------
	srv := internal.NewServer()

	// TODO don't just fire & forget
	go srv.Run()

	go func() {
		for !srv.Done {
			select {
			case err := <-srv.ErrChan:
				t.Errorf("server produced an error: %v\n", err)
			}
		}
	}()

	// -- --------- CLIENT TESTS ---------
	// test identity
	count := 3
	log.Infof("starting identity test with %d clients\n", count)

	var clients = make([]internal.Client, count)

	// establish connections first
	for id := 0; id < count; id++ {
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			t.Errorf("error connecting to server: %v\n", err)
		}
		clients[id] = internal.Client{ID: uint8(id), Conn: conn}
	}

	// range over clients, confirm they have the expected IDs & close connections
	for _, c := range clients {
		log.Infof("checking client with id '%d'", c.ID)

		_, err := c.Conn.Write([]byte("0"))
		if err != nil {
			t.Errorf("client couldn't request identity: %v\n", err)
		}

		buff := make([]byte, 100)
		reader := bufio.NewReader(c.Conn)

		log.Info("before reading first byte")
		size, err := reader.ReadByte()
		if err != nil {
			t.Errorf("error reading response size: %v\n", err)
		}

		log.Infof("size: %v", size)

		log.Info("before reading all byte")
		_, err = reader.Read(buff)
		// _, err = io.ReadFull(reader, buff[1:size])
		if err != nil {
			t.Errorf("error reading full response: %v\n", err)
		}

		clientId := uint8(binary.BigEndian.Uint16(buff[:size]))

		log.Info("before comparing bytes")
		if clientId != c.ID {
			t.Errorf("got client ID '%d' want '%d'\n", clientId, c.ID)
		}

		c.Conn.Close()

		log.Infof("no errors found for client id '%d'\n", c.ID)
	}

	srv.Stop()

	// TODO how do we receive messages from existing clients? Threadpool?

	// TODO test concurrent requests

	// TODO test list messages

	// TODO test message clients; but limit number of clients between 5-10

	// TODO make sure can't connect more than max clients
	// for cIdx := range internal.MAX_CLIENTS {
	// 	client.connectToServeR()
	// }
}

func connectToServer() {}
