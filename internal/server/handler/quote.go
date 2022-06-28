package handler

import (
	"context"
	"errors"
	"log"

	"github.com/SergeySlonimsky/pow/internal/server"
	"github.com/SergeySlonimsky/pow/pkg/protocol"
)

//go:generate mockgen -source=./quote.go -destination=./mock/quote_mock.go

var ErrInvalidType = errors.New("cannot find handler for the type")

type quoteStorage interface {
	GetRandomQuote() string
}

type pow interface {
	Generate(ctx context.Context, resource string) (string, error)
	Verify(ctx context.Context, resource, data string) error
}

type QuoteHandler struct {
	quoteStorage quoteStorage
	pow          pow
}

func New(quoteStorage quoteStorage, pow pow) *QuoteHandler {
	return &QuoteHandler{
		quoteStorage: quoteStorage,
		pow:          pow,
	}
}

//nolint:ireturn // to implement handler interface
// Handle handles request from clients, and determines, which child handler should be run.
func (h *QuoteHandler) Handle(ctx context.Context, req server.Request) (server.Response, error) {
	switch req.GetMessage().GetType() { //nolint:exhaustive // protocol.TypeErr should return err
	case protocol.TypeResource:
		return h.handleResource(ctx, req)
	case protocol.TypeChallenge:
		return h.handleChallenge(ctx, req)
	default:
		return nil, ErrInvalidType
	}
}

func (h *QuoteHandler) handleResource(ctx context.Context, req server.Request) (protocol.Message, error) {
	log.Printf("called \"handleResource\": with %s", req.GetMessage().ToString())

	if err := h.pow.Verify(ctx, req.GetAddr(), req.GetMessage().GetBody()); err != nil {
		return protocol.Message{}, err
	}

	quote := h.quoteStorage.GetRandomQuote()

	return protocol.NewMessage(protocol.TypeResource, quote), nil
}

func (h *QuoteHandler) handleChallenge(ctx context.Context, req server.Request) (protocol.Message, error) {
	log.Printf("called \"handleChallenge\": with %s", req.GetMessage().ToString())

	challenge, err := h.pow.Generate(ctx, req.GetAddr())
	if err != nil {
		return protocol.Message{}, err
	}

	return protocol.NewMessage(protocol.TypeChallenge, challenge), nil
}
