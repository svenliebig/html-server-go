package http

import (
	"fmt"
	"log"
	"testing"
)

func TestServer(t *testing.T) {
	t.Run("should add a route and be able to call it with a handler", func(t *testing.T) {
		s := Server{}

		id, err := s.AddRoute("GET", "/hello")

		if err != nil {
			log.Fatal(err)
		}

		c := make(chan Handler, 1)

		s.SubscribeTo(id, c)

		go func() {
			for request := range c {
				request(func(req *Request, res *Response, n string) {
					fmt.Println("Request:", req.Method, req.Path, req.Host)
					res.Write()
				})
			}
		}()

		s.Listen()
	})
}
