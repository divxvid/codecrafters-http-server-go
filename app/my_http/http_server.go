package myhttp

import (
	"fmt"
	"net"
)

type Server struct {
	Host string
	Port int

	Listener net.Listener
	Router   *Router
}

func NewServer(router *Router) *Server {
	return &Server{
		Router: router,
	}
}

func (s *Server) Start(host string, port int) error {
	fmt.Printf("Starting to listen over %s:%d\n", host, port)
	s.Host = host
	s.Port = port
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	s.Listener = listener
	go s.AcceptLoop()
	return nil
}

func (s *Server) AcceptLoop() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Printf("Error while accepting a connection: %v\n", err)
		}
		go s.ProcessConnection(conn)
	}
}

func (s *Server) ProcessConnection(conn net.Conn) {
	defer conn.Close()

	request, err := FromReader(conn)
	if err != nil {
		fmt.Println("Encountered an error", err)
		return
	}

	response := s.Router.HandleRequest(request)
	response.WriteToConn(conn)
}
