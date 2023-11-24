package ingressgateway

import (
	"net"

	"github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
	"github.com/rs/zerolog/log"
)

type Server struct {
	header      *header.Header
	destination string
	closed      bool
	listener    net.Listener
}

func New(envID, destination string) *Server {
	h := header.NewV1()
	h.Set(propagation.EnvIDHeader, envID)
	return &Server{
		header:      h,
		destination: destination,
	}
}

func (s *Server) Start(addr string) {
	s.closed = false
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal().Err(err).Str("listen address", addr).Msg("failed to listen")
	}
	defer ln.Close()
	s.listener = ln

	log.Info().Msg("starting server")
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed {
				return
			}
			log.Error().Err(err).Str("listen address", addr).Msg("failed to accept")
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) Close() {
	log.Info().Msg("shutdown")
	s.closed = true
	s.listener.Close()
}

func (s *Server) handle(clientConn net.Conn) {
	defer clientConn.Close()

	serverConn, err := net.Dial("tcp", s.destination)
	if err != nil {
		log.Error().
			Err(err).
			Stringer("client local address", clientConn.LocalAddr()).
			Stringer("client remote address", clientConn.RemoteAddr()).
			Str("server address", s.destination).
			Msg("failed to dial server")
		return
	}
	defer serverConn.Close()

	_, err = s.header.WriteTo(serverConn)
	if err != nil {
		log.Error().
			Err(err).
			Stringer("client local address", clientConn.LocalAddr()).
			Stringer("client remote address", clientConn.RemoteAddr()).
			Str("server address", s.destination).
			Msg("failed to write header")
		return
	}

	err = proxy(clientConn, serverConn)
	if err != nil {
		log.Error().
			Err(err).
			Stringer("client local address", clientConn.LocalAddr()).
			Stringer("client remote address", clientConn.RemoteAddr()).
			Stringer("server local address", serverConn.LocalAddr()).
			Stringer("server remote address", serverConn.RemoteAddr()).
			Msg("failed to proxy")
	}
}
