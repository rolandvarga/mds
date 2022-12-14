package test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"testing"

	"github/rolandvarga/mds/internal"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		require.NoError(t, err)

		clients[id] = internal.Client{ID: uint8(id), Conn: conn}
	}

	// range over clients, confirm they have the expected IDs & close connections
	for _, c := range clients {
		log.Infof("checking client with id '%d'", c.ID)

		_, err := c.Conn.Write([]byte("0"))
		require.NoError(t, err, "client couldn't request identity: %v\n", err)

		// create a temp buffer
		tmp := make([]byte, 500)
		c.Conn.Read(tmp)

		// convert bytes into Buffer (which implements io.Reader/io.Writer)
		tmpBuff := bytes.NewBuffer(tmp)
		resp := new(internal.Response)

		// creates a decoder object
		decoder := gob.NewDecoder(tmpBuff)
		decoder.Decode(resp)

		errMsg := fmt.Sprintf("msg type '%v' client id '%s'\nwant type '%v' with id '%d'\n",
			resp.Type, resp.Msg, internal.Identity, c.ID,
		)

		assert.Equal(t, internal.Identity, resp.Type, errMsg)
		assert.Equal(t, fmt.Sprint(c.ID), resp.Msg, errMsg)

		c.Conn.Close()
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
