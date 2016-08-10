package network

import "net"

// A Connection represents a TCP connection
type Connection struct {
	socket *net.TCPConn
}

// Open opens a TCP connection
func (c *Connection) Open(address string) error {
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}

	c.socket, err = net.DialTCP("tcp", nil, tcpAddress)
	return err
}

// Close closes the connection
func (c *Connection) Close() error {
	return c.socket.Close()
}

// Send sends an array of bytes
func (c *Connection) Send(b []byte) (int, error) {
	return c.socket.Write(b)
}
