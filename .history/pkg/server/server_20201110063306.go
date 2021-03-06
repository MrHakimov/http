package server

import (
	"log"
	"net"
	"sync"
)

// Constants
const (
	HOST = "0.0.0.0"
	PORT = "9999"
)

// HandlerFunc connection
type HandlerFunc func(conn net.Conn)

type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
}

func NewServer(add string) *Server {
	return &Server{
		addr:     add,
		handlers: make(map[string]HandlerFunc),
	}
}

func (s *Server) Register(path string, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

func (s *Server) Start() error {
	listner, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
		return err
	}

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.handle(Request{
			Conn: conn,
		})
	}
}
