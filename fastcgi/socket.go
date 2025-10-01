package fastcgi

import (
	"context"
	"log"
	"net"
	"sync"
)

const DefaultSocketPath = "/var/run/php-fubar.sock"

type Socket struct {
	socketPath       string
	listener         net.Listener
	handleConnection func(net.Conn, context.Context)
	wg               sync.WaitGroup
	ctx              context.Context
	cancel           context.CancelFunc
}

func NewSocket() (*Socket, error) {
	socketPath := DefaultSocketPath
	listener, err := net.Listen("unix", socketPath)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Socket{
		socketPath: socketPath,
		listener:   listener,
		wg:         sync.WaitGroup{},
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

func MustSocket() *Socket {
	socket, err := NewSocket()

	if err != nil {
		panic(err)
	}

	return socket
}

func (s *Socket) Handle(handler func(net.Conn, context.Context)) {
	s.handleConnection = handler
}

func (s *Socket) Close() error {
	if err := s.listener.Close(); err != nil {
		return err
	}
	s.cancel()
	s.wg.Wait()

	return nil
}

func (s *Socket) Listen() {

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			select {
			case <-s.ctx.Done():
				return
			default:
				log.Println("Error accepting connection", err)
				continue
			}
		}

		s.wg.Add(1)
		go func(conn net.Conn) {
			defer s.wg.Done()
			defer conn.Close()

			s.handleConnection(conn, s.ctx)
		}(conn)
	}
}
