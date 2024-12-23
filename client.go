package main

import (
	"log"
	"net"
)

type client struct {
	Name string
	conn net.Conn
	srv  *Server
}

func (s *Server) createClient(conn net.Conn, name string) *client {

	c := &client{
		Name: name,
		srv:  s,
		conn: conn,
	}

	log.Printf("new client created --> %s - remote addr = %s", c.Name, c.conn.RemoteAddr())
	log.Printf("received on %s", s.Addr)
	return c

}
