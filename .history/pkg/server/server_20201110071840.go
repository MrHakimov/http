package server

import (
	"bytes"
	"io"
	"strconv"
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

		go s.handle((Request{
			Conn: conn,
		}))
	}
}

func (s *Server) handle(req Request) {
	defer func() {
		if closeErr := req.Conn.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	buf := make([]byte, 4096)
	n, err := req.Conn.Read(buf)
	if err == io.EOF {
		log.Printf("%s", buf[:n])
	}

	data := buf[:n]
	requestLineDelim := []byte{'\r', '\n'}
	requestLineEnd := bytes.Index(data, requestLineDelim)
	if requestLineEnd == -1 {
		log.Print("requestLineEndErr: ", requestLineEnd)
	}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		log.Print("partsErr: ", parts)
	}

	path := parts[1]
	if strings.Contains(path, "payments") {
		uri, err := url.ParseRequestURI(path)
		if err != nil {
			log.Println("url query parse err: ", err)
		}

		if uri.RawQuery != "" {
			req.QueryParams = uri.Query()
			log.Println(req.QueryParams["id"])
		} else {
			split := strings.Split(uri.Path, "/payments/")
			m := make(map[string]string)
			m["id"] = split[1]
		}

		path = uri.Path
	}
	if err != nil {
		log.Print(err)
	}

	s.mu.RLock()
	if handler, ok := s.handlers[path]; ok {
		s.mu.RUnlock()
		handler(&req)
	}
	return
}

func (s *Server) Response(body string) string {
	return "HTTP/1.1 200 OK\r\n" +
		"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
		"Content-Type: text/html\r\n" +
		"Connection: close\r\n" +
		"\r\n" + body
}
