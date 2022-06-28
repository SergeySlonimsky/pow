package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Type string

const (
	TypeChallenge Type = "challenge" // TypeChallenge is sent by server, when pass challenge, and by client, when client wants to be challenged.
	TypeResource  Type = "resource"  // TypeResource is sent by server when send resource, and by client, when client passed the challege.
	TypeErr       Type = "error"     // TypeResource is sent by server and client when smth went wrong.
)

// Delimiter between type and body.
const delimiter = "|"

// Message is a struct for communication between client and server.
// Contains message type for determine, how this message should be understood, and body with message payload.
type Message struct {
	messageType Type
	body        string
}

// NewMessage creates Message from message type and body
func NewMessage(messageType Type, body string) Message {
	return Message{
		messageType: messageType,
		body:        body,
	}
}

// ParseFromReader parses string until \n from reader and tries to create Message from it.
func ParseFromReader(rd io.Reader) (Message, error) {
	reader := bufio.NewReader(rd)

	request, err := reader.ReadString('\n')
	if err != nil {
		return Message{}, err
	}

	request = strings.TrimSuffix(request, "\n")

	parts := strings.Split(request, delimiter)

	if len(parts) == 0 {
		return Message{}, errors.New("invalid request format")
	}

	messageType := convertMessageType(parts[0])

	body := ""
	if len(parts) > 1 {
		body = parts[1]
	}

	return Message{
		messageType: messageType,
		body:        body,
	}, nil
}

func (m Message) GetType() Type {
	return m.messageType
}

func (m Message) GetBody() string {
	return m.body
}

// ToString returns string view of message separated by delimiter. When body is empty, just returns message type.
func (m Message) ToString() string {
	if m.body == "" {
		return string(m.messageType)
	}

	return fmt.Sprintf("%s|%s", m.messageType, m.body)
}

func convertMessageType(incomeType string) Type {
	switch incomeType {
	case "resource":
		return TypeResource
	case "challenge":
		return TypeChallenge
	default:
		return TypeErr
	}
}
