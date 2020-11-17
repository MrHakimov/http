package server

import (
	"bytes"
	"io"
	"strings"

	"log"
	"net"
	"net/url"
	"sync"
)

// HandlerFunc handler
type HandlerFunc func(req *Request)

// Server class
type Server struct {
	addr string

	mu sync.RWMutex

	handlers map[string]HandlerFunc
}

// Request class
type Request struct {
	Conn        net.Conn
	QueryParams url.Values
}

// NewServer can create new servers
func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandlerFunc)}

}

// Register path
func (s *Server) Register(path string, handler HandlerFunc) {
	s.mu.RLock()
	s.handlers[path] = handler
	s.mu.RUnlock()
}

// Start is main function
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		if cerr := listener.Close(); cerr != nil {

			if err == nil {
				err = cerr
				return
			}
			log.Print(cerr)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, (1024 * 8))
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			log.Printf("%s", buf[:n])
		}
		if err != nil {
			log.Println(err)
			return
		}

		var req Request
		data := buf[:n]
		rLD := []byte{'\r', '\n'}
		rLE := bytes.Index(data, rLD)
		if rLE == -1 {
			log.Printf("Bad Request")
			return
		}

		hLD := []byte{'\r', '\n', '\r', '\n'}
		hLE := bytes.Index(data, hLD)
		if rLE == -1 {
			return
		}

		headersLine := string(data[rLE:hLE])
		headers := strings.Split(headersLine, "\r\n")[1:]
		mp := make(map[string]string)
		for _, v := range headers {
			headerLine := strings.Split(v, ": ")
			mp[headerLine[0]] = headerLine[1]
		}

		b := string(data[hLE:])
		b = strings.Trim(b, "\r\n")

		reqLine := string(data[:rLE])
		parts := strings.Split(reqLine, " ")

		if len(parts) != 3 {
			log.Println("Error request")
			return
		}

		path, version := parts[1], parts[2]
		if version != "HTTP/1.1" {
			log.Println("Error version")
			return
		}

		decode, err := url.PathUnescape(path)
		if err != nil {
			log.Println(err)
			return
		}

		uri, err := url.ParseRequestURI(decode)
		if err != nil {
			log.Println(err)
			return
		}

		req.Conn = conn
		req.QueryParams = uri.Query()

		var handler = func(req *Request) { conn.Close() }

		s.mu.RLock()

		handler(&req)
		s.mu.RUnlock()

	}
}
