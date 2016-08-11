package network

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"net"
)

// A Connection represents a TCP connection
type Connection struct {
	socket  net.Conn
	channel chan string
}

// NewConnection creates and setups new transmission connection
func NewConnection(socket net.Conn) (conn *Connection) {
	conn = new(Connection)
	conn.socket = socket
	conn.channel = make(chan string)

	go conn.receive()
	return
}

// Close closes the connection
func (c *Connection) Close() error {
	return c.socket.Close()
}

// Send sends an array of bytes
func (c *Connection) Send(b []byte) (int, error) {
	return c.socket.Write(b)
}

func (c *Connection) receive() {
	defer close(c.channel)
	for {
		reader := bufio.NewReader(c.socket)
		buffer := make([]byte, 2)

		// read the size of the string
		_, err := reader.Read(buffer)
		if err != nil {
			log.Error(err.Error())
			return
		}

		size := binary.BigEndian.Uint16(buffer)
		buffer = make([]byte, size)

		// read the data string
		_, err = reader.Read(buffer)
		if err != nil {
			log.Error(err.Error())
			return
		}
		var dat map[string]interface{}
		if err = json.Unmarshal(buffer, &dat); err != nil {
			return
		}
		println(string(buffer))
	}
}
