package http

import (
	"fmt"
	"log"
	"math/rand"
	"net"
)

type Server struct {
}

type Request func(func (req, res, n string))

const max = uint16(1 << 15)

var subscribers = map[uint16]chan Request{}
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

func (s *Server) SubscribeTo(n uint16, c chan Request) {
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

    subscribers[routes["/hello"]] <- func(h func (req, res, n string)) {
      h("hello", "mister", "world")
    }

    if err != nil {
      fmt.Println("Error accepting: ", err.Error())
      continue
    }

    go handleConnection(conn)
  }
}

func handleConnection(conn net.Conn) {
  buf := make([]byte, 1024)
  n, err := conn.Read(buf)

  if err != nil {
    fmt.Println("Error reading:", err.Error())
    return
  }

  req := string(buf[:n])
  fmt.Println("Received:", req)

  // Send the response back to the client
  resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\t\n\t\nHello World\r\n"
  _, err = conn.Write([]byte(resp))

  if err != nil {
    fmt.Println("Error writing:", err.Error())
    return
  }

  conn.Close()
}
