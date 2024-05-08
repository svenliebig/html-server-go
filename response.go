package http

import (
	"fmt"
	"net"
)

type Response struct {
  conn net.Conn
}

func (r *Response) Write() {
  resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\t\n\t\nHello World\r\n"

  _, err := r.conn.Write([]byte(resp))

  if err != nil {
    fmt.Println("Error writing:", err.Error())
    return
  }

  r.conn.Close()
}
