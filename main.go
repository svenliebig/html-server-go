package main

import (
	"fmt"
	"log"
  "github.com/svenliebig/html-server-go/http"
)

func main() {
  s := http.Server{}

  id, err := s.AddRoute("GET", "/hello")
  
  if err != nil {
    log.Fatal(err)
  }

  c := make(chan http.Handler, 1)

  s.SubscribeTo(id, c)

  go func () {
    for request := range c {
      request(func(req *http.Request, res *http.Response, n string) {
        fmt.Println("Request:", req.Method, req.Path, req.Host)
        res.Write()
      })
    }
  }()

  s.Listen()
}

