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

// Server class
type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
}

// NewServer creates new server
func NewServer(add string) *Server {
	return &Server{
		addr:     add,
		handlers: make(map[string]HandlerFunc),
	}
}

// Register new server
func (s *Server) Register(path string, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

// Start is main function
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
