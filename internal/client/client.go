package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/SergeySlonimsky/pow/pkg/hashcash"
	"github.com/SergeySlonimsky/pow/pkg/protocol"
)

const defaultGenerationAttempts = 999999

func Run(_ context.Context, addr string) error { //nolint:gocyclo,cyclop // has to be refactored
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	defer conn.Close()

	if err := sendMessage(conn, protocol.NewMessage(protocol.TypeChallenge, "")); err != nil {
		return fmt.Errorf("send message: %s", err)
	}

	for {
		msg, err := readMessage(conn)
		if err != nil {
			return sendMessage(conn, createErrorMessage(err))
		}

		switch msg.GetType() {
		case protocol.TypeChallenge:
			log.Printf("challenge received: %s", msg.ToString())

			stamp, err := hashcash.FromString(msg.GetBody())
			if err != nil {
				return sendMessage(conn, createErrorMessage(err))
			}

			if err := stamp.GenerateHash(defaultGenerationAttempts); err != nil {
				return sendMessage(conn, createErrorMessage(err))
			}

			if err := sendMessage(conn, protocol.NewMessage(protocol.TypeResource, stamp.ToString())); err != nil {
				return fmt.Errorf("send message: %s", err)
			}
		case protocol.TypeResource:
			log.Printf("Quote received: %s", msg.GetBody())

			return nil
		case protocol.TypeErr:
			return errors.New(msg.GetBody())
		}
	}
}

func readMessage(r io.Reader) (protocol.Message, error) {
	return protocol.ParseFromReader(r)
}

func sendMessage(conn net.Conn, msg protocol.Message) error {
	text := fmt.Sprintf("%s\n", msg.ToString())

	if _, err := conn.Write([]byte(text)); err != nil {
		return fmt.Errorf("error writing data: %s", err.Error())
	}

	return nil
}

func createErrorMessage(err error) protocol.Message {
	return protocol.NewMessage(protocol.TypeErr, err.Error())
}
