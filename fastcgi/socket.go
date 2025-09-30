package fastcgi

import (
	"log"
	"net"
)

func NewSocket() (*Socket, error) {
	socketPath := DefaultSocketPath
	listener, err := net.Listen("unix", socketPath)

	if err != nil {
		return nil, err
	}

	return &Socket{
		socketPath: socketPath,
		listener:   listener,
	}, nil
}

func MustSocket() *Socket {
	socket, err := NewSocket()

	if err != nil {
		panic(err)
	}

	return socket
}

func (socket *Socket) Listen() {
	defer socket.listener.Close()

	for {
		conn, err := socket.listener.Accept()

		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}

		go handleConncetion(conn)
	}
}

func handleConncetion(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 4096)

	_, err := conn.Read(buffer)

	if err != nil {
		log.Println("Error reading from connection", err)
	}

	_, _ = conn.Write(buffer)
}
