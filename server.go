package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {

	//Basic server name for readability
	Name string

	//Network Address for our server
	Addr    string
	TCPAddr *net.TCPAddr

	//Our listener
	listenerConfig net.ListenConfig
	listener       net.Listener

	//Concurrency, Context & sync
	srvWG     *sync.WaitGroup
	srvCtx    context.Context
	srvCancel context.CancelFunc
	mu        sync.Mutex
}

func NewServer(serverName, host, port string, lc net.ListenConfig) *Server {

	addr := net.JoinHostPort(host, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	createdAt := time.Now()

	name := fmt.Sprintf("%s_%s", serverName, createdAt.Format("20060102150405"))

	ctx, cancel := context.WithCancel(context.Background())

	s := &Server{
		Name:      name,
		Addr:      addr,
		TCPAddr:   tcpAddr,
		srvCtx:    ctx,
		srvCancel: cancel,

		listenerConfig: lc,
	}
	return s
}

func (s *Server) StartServer() {

	log.Printf("Starting server: %s", s.Name)

	s.AcceptLoop("client-test")

}

func (s *Server) AcceptLoop(name string) {

	s.mu.Lock()

	ctx, cancel := context.WithCancel(s.srvCtx)
	defer cancel()

	log.Printf("Starting client accept loop -- %s\n", name)

	l, err := s.listenerConfig.Listen(ctx, s.TCPAddr.Network(), s.TCPAddr.String())
	if err != nil {
		log.Printf("Error creating listener: %s\n", err)
	}

	log.Printf("Listening on %s\n", l.Addr())

	// Add listener to the server
	s.listener = l

	// Can begin go-routine for accepting connections
	go s.acceptConnection(l, "client-test",
		func(conn net.Conn) {
			s.createClient(conn, "normal-client")
		},
		func(err error) bool {
			select {
			case <-ctx.Done():
				//log.Println("accept loop context canceled -- exiting loop")
				return true
			default:
				//log.Printf("accept loop context error -- %s\n", err)
				return false
			}
		})

	s.mu.Unlock()
}

func (s *Server) acceptConnection(l net.Listener, name string, createConnFunc func(conn net.Conn), customErr func(err error) bool) {

	for {
		conn, err := l.Accept()
		if err != nil {
			if customErr != nil && customErr(err) {
				log.Println("custom error called")
				return // we break here to come out of the loop - if we can't reconnect during a reconnect strategy then we break
			}
			continue
		}
		go createConnFunc(conn)
	}
}
