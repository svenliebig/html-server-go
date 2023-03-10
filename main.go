package main

import (
	"fmt"
	"log"
	"net"
  "github.com/svenliebig/html-server-go/http"
)

// HTTP Protocol: https://www.rfc-editor.org/rfc/rfc7231
// Content Type:  https://www.rfc-editor.org/rfc/rfc2046

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

func main() {
  s := http.Server{}

  id, err := s.AddRoute("GET", "/hello")
  
  if err != nil {
    log.Fatal(err)
  }

  c := make(chan http.Request, 1)

  s.SubscribeTo(id, c)

  go func () {
    for request := range c {
      request(func(req, res, n string) {
        fmt.Println(req, res, n)
      })
    }
  }()

  s.Listen()
}

func mainmelater() {
  listener, err := net.Listen("tcp", ":8080")

  if err != nil {
    log.Fatal("Error listening:", err.Error())
  }

  defer listener.Close()
  
  fmt.Println("Listening on :8080")

  for {
    conn, err := listener.Accept()

    if err != nil {
      fmt.Println("Error accepting: ", err.Error())
      continue
    }

    go handleConnection(conn)
  }
}
