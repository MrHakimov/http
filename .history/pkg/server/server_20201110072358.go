package server

import (
	"bytes"
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
	PathParams  map[string]string
	Headers     map[string]string
	Body        []byte
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

	buf := make([]byte, (1024 * 50))

	for {
		readBytes, err := conn.Read(buf)
		if err != nil {
			return
		}

		data := buf[:readBytes]
		ldelim := []byte{'\r', '\n'}
		index := bytes.Index(data, ldelim)
		if index == -1 {
			log.Println("delim chars not found :(")
			return
		}
		var req Request
		var good bool = true
		var path1 string = ""
		rline := string(data[:index])
		parts := strings.Split(rline, " ")
		var header []byte = data[index+2:]
		req.Headers = make(map[string]string)
		req.PathParams = make(map[string]string)
		if len(parts) == 3 {
			_, path, version := parts[0], parts[1], parts[2]
			decode, err := url.PathUnescape(path)
			if err != nil {
				log.Println(err)
				return
			}
			if version != "HTTP/1.1" {
				log.Println("version is not valid")
				return
			}
			url, err := url.ParseRequestURI(decode)
			if err != nil {
				log.Println(err)
				return
			}
			req.Conn = conn
			req.QueryParams = url.Query()
			partsPath := strings.Split(url.Path, "/")
			for cur := range s.handlers {
				partsCur := strings.Split(cur, "/")
				if len(partsPath) != len(partsCur) {
					continue
				}
				var n int = len(partsPath)
				for i := 0; i < n && good == true; i++ {
					var l int = strings.Index(partsCur[i], "{")
					var r int = strings.LastIndex(partsCur[i], "}")
					var cnt int = strings.Count(partsCur[i], "{") +
						strings.Count(partsCur[i], "}")
					if cnt == 0 {
						if partsCur[i] != partsPath[i] {
							good = false
						}
					} else if cnt == 2 {
						req.PathParams[partsCur[i][l+1:r]] = partsPath[i][l:]
					} else {
						good = false
					}
				}
				if good == false {
					req.PathParams = make(map[string]string)
				} else {
					path1 = cur
					break
				}
			}
			log.Println("decode(path:)", decode)
			log.Println("url.Query():", url.Query())
			log.Println("url.Path:", url.Path)
			log.Println("path1:", path1)
			log.Println("req.PathParams:", req.PathParams)
		}
		/// Headers...
		var body []byte
		if len(header) > 0 {
			ldelim := []byte{'\r', '\n', '\r', '\n'}
			index := bytes.Index(header, ldelim)
			if index == -1 {
				log.Println("delim ^ 2 chars not found :(")
				return
			}
			body = header[index+4:]
			data := string(header[:index])
			log.Println("data(header):", data)
			lheader := strings.Split(data, "\r\n")
			for _, header := range lheader {
				index := strings.Index(header, ":")
				if index == -1 {
					log.Println("index for seperating key and value not found")
					return
				}
				key, value := header[:index], header[index+2:]
				req.Headers[key] = value // join them
			}
			log.Println("Headers: ", req.Headers)
		}
		// Body...
		req.Body = body
		log.Println("Body:", string(body))

		log.Println()
		var f = func(req *Request) {}

		s.mu.RLock()
		f, good = s.handlers[path1]
		s.mu.RUnlock()

		if good == false {
			conn.Close()
		} else {
			f(&req)
		}
	}
}

// Response common answer
func (s *Server) Response(body string) string {
	return "HTTP/1.1 200 OK\r\n" +
		"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
		"Content-Type: text/html\r\n" +
		"Connection: close\r\n" +
		"\r\n" + body
}
