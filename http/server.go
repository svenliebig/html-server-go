package http

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
)

// HTTP Protocol: https://www.rfc-editor.org/rfc/rfc7231
// Content Type:  https://www.rfc-editor.org/rfc/rfc2046

type Server struct {
}

type Handler func(func (req *Request, res *Response, n string))

const max = uint16(1 << 15)

var subscribers = map[uint16]chan Handler{}
var routes = map[string]uint16{}

func (s *Server) AddRoute(method, route string) (uint16, error) {
  var (
    ok bool
    r uint16
  )

  for !ok {
    ok = true
    r = uint16(rand.Intn(int(max)))
    for vr, v := range routes {
      if route == vr {
        return 0, fmt.Errorf("route '%s' already exists", route)
      }

      if v == r {
        ok = false
      }
    }
  }

  routes[route] = r

  return r, nil
}

func (s *Server) SubscribeTo(n uint16, c chan Handler) {
  subscribers[n] = c
}

func (s *Server) Listen() {
  listener, err := net.Listen("tcp", ":8080")

  if err != nil {
    log.Fatal("Error listening:", err.Error())
  }

  defer listener.Close()
  
  fmt.Println("Listening on :8080")

  for {
    fmt.Println("Waiting for connection")
    conn, err := listener.Accept()

    if err != nil {
      fmt.Println("Error accepting: ", err.Error())
      continue
    }

    req, err := readConnection(conn)

    if err != nil {
      fmt.Println("Error parsing request: ", err.Error())
      continue
    }

    fmt.Println("Connection established")

    subscribers[routes["/hello"]] <- func(h func (req *Request, res *Response, n string)) {
      h(req, &Response{ conn }, "world")
    }
  }
}

func readConnection(conn net.Conn) (*Request, error) {
  buf := make([]byte, 1024)
  n, err := conn.Read(buf)

  if err != nil {
    return nil, fmt.Errorf("%w: Error reading", err)
  }

  rs := string(buf[:n])

  lines := strings.Split(rs, "\n")
  head, headers := pop(lines)

  r := &Request{}

  split := strings.Split(head, " ")
  method := split[0]
  path := split[1]

  switch method {
    // https://www.rfc-editor.org/rfc/rfc7231#section-4.1
    case "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE", "CONNECT":
      r.Method = method
    default:
      return nil, fmt.Errorf("invalid method: %s", method)
  }

  r.Path = path

  for _, v := range headers {
    n := strings.IndexByte(v, ':')
    if n == -1 {
      fmt.Println("Invalid header:", v)
      continue
    }

    switch v[:n] {
    case "Host":
      r.Host = strings.Trim(v[n+1:], " ")
    default:
      fmt.Println("not handled header:", v)
    }
  }

  return r, nil
}

func pop(s []string) (string, []string) {
  return s[0], s[1:]
}
