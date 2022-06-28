package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/SergeySlonimsky/pow/pkg/protocol"
)

type Handler interface {
	Handle(ctx context.Context, req Request) (Response, error)
}

type Server struct {
	handler Handler
}

func New(handler Handler) *Server {
	return &Server{
		handler: handler,
	}
}

func (s *Server) Run(ctx context.Context, port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	defer func() {
		if err := listener.Close(); err != nil {
			log.Printf("error close tcp listener: %s", err.Error())
		}
	}()

	fmt.Println("listening", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accept connection: %w", err)
		}

		go s.handleConnection(ctx, conn)
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	log.Printf("new connection from %s", conn.RemoteAddr().String())

	for {
		message, err := protocol.ParseFromReader(conn)
		if err != nil {
			writeData(conn, createErrorMessage(err))

			return
		}

		if message.GetType() == protocol.TypeErr {
			if err = conn.Close(); err != nil {
				log.Printf("error close connection: %s", err.Error())

				return
			}
		}

		addr, err := cleanClientAddr(conn.RemoteAddr().String())
		if err != nil {
			writeData(conn, createErrorMessage(err))

			return
		}

		req := messageRequest{
			message: message,
			addr:    addr,
		}

		resp, err := s.handler.Handle(ctx, req)
		if err != nil {
			writeData(conn, createErrorMessage(err))
		} else {
			writeData(conn, resp)
		}
	}
}

func writeData(conn net.Conn, resp Response) {
	restText := fmt.Sprintf("%s\n", resp.ToString())

	if _, err := conn.Write([]byte(restText)); err != nil {
		log.Printf("error writing data: %s", err.Error())
	}
}

func cleanClientAddr(addr string) (string, error) {
	parts := strings.Split(addr, ":")
	if len(parts) > 0 {
		return parts[0], nil
	}

	return "", errors.New("invalid client address")
}

func createErrorMessage(err error) protocol.Message {
	return protocol.NewMessage(protocol.TypeErr, err.Error())
}
