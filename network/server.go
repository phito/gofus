package network

import (
	"net"

	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("server")

// A Server represents a TCP listener
type Server struct {
	socket  net.Listener
	running bool
}

// Open opens a TCP listener
func (server *Server) Open(address string) (err error) {
	// open the socket
	server.socket, err = net.Listen("tcp", address)
	if err != nil {
		return
	}

	// start the worker
	go server.run()

	return
}

func (server *Server) run() {
	defer server.socket.Close()

	server.running = true
	for server.running {
		_, err := server.socket.Accept()
		if err != nil {
			log.Error("Accept failed: ", err)
			return
		}
		log.Info("New client connected")
	}
}
