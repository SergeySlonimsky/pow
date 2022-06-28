package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/SergeySlonimsky/pow/internal/server"
	"github.com/SergeySlonimsky/pow/internal/server/handler"
	mockHandler "github.com/SergeySlonimsky/pow/internal/server/handler/mock"
	"github.com/SergeySlonimsky/pow/pkg/protocol"
)

type testRequest struct {
	message protocol.Message
	addr    string
}

func (t testRequest) GetMessage() protocol.Message {
	return t.message
}

func (t testRequest) GetAddr() string {
	return t.addr
}

func TestQuoteHandler_Handle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)

	t.Run("handleResource", func(t *testing.T) {
		tests := []struct {
			name     string
			mockFunc func(ctrl *gomock.Controller) (server.Request, server.Response, *mockHandler.Mockpow, *mockHandler.MockquoteStorage)
			wantErr  bool
		}{
			{
				name: "valid pow and quote",
				mockFunc: func(ctrl *gomock.Controller) (server.Request, server.Response, *mockHandler.Mockpow, *mockHandler.MockquoteStorage) {
					req := testRequest{
						message: protocol.NewMessage(protocol.TypeResource, "challenge"),
						addr:    "192.168.1.1",
					}
					resp := protocol.NewMessage(protocol.TypeResource, "test quote")

					pow := mockHandler.NewMockpow(ctrl)
					storage := mockHandler.NewMockquoteStorage(ctrl)

					pow.EXPECT().Verify(ctx, "192.168.1.1", "challenge").Return(nil)
					storage.EXPECT().GetRandomQuote().Return("test quote")

					return req, resp, pow, storage
				},
				wantErr: false,
			},
			{
				name: "invalid pow",
				mockFunc: func(ctrl *gomock.Controller) (server.Request, server.Response, *mockHandler.Mockpow, *mockHandler.MockquoteStorage) {
					req := testRequest{
						message: protocol.NewMessage(protocol.TypeResource, ""),
						addr:    "192.168.1.1",
					}
					resp := protocol.NewMessage(protocol.TypeResource, "resource")

					pow := mockHandler.NewMockpow(ctrl)
					storage := mockHandler.NewMockquoteStorage(ctrl)

					pow.EXPECT().Verify(ctx, "192.168.1.1", "").Return(errors.New("invalid pow"))

					return req, resp, pow, storage
				},
				wantErr: true,
			},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				req, resp, pow, storage := tt.mockFunc(ctrl)

				h := handler.New(storage, pow)

				got, err := h.Handle(ctx, req)

				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, resp.ToString(), got.ToString())
				}
			})
		}
	})

	t.Run("handleChallenge", func(t *testing.T) {
		tests := []struct {
			name     string
			mockFunc func(ctrl *gomock.Controller) (server.Request, server.Response, *mockHandler.Mockpow, *mockHandler.MockquoteStorage)
			wantErr  bool
		}{
			{
				name: "valid challenge",
				mockFunc: func(ctrl *gomock.Controller) (server.Request, server.Response, *mockHandler.Mockpow, *mockHandler.MockquoteStorage) {
					req := testRequest{
						message: protocol.NewMessage(protocol.TypeChallenge, ""),
						addr:    "192.168.1.1",
					}
					resp := protocol.NewMessage(protocol.TypeChallenge, "test challenge")

					pow := mockHandler.NewMockpow(ctrl)
					storage := mockHandler.NewMockquoteStorage(ctrl)

					pow.EXPECT().Generate(ctx, "192.168.1.1").Return("test challenge", nil)

					return req, resp, pow, storage
				},
				wantErr: false,
			},
			{
				name: "invalid challenge",
				mockFunc: func(ctrl *gomock.Controller) (server.Request, server.Response, *mockHandler.Mockpow, *mockHandler.MockquoteStorage) {
					req := testRequest{
						message: protocol.NewMessage(protocol.TypeChallenge, ""),
						addr:    "192.168.1.1",
					}
					resp := protocol.NewMessage(protocol.TypeChallenge, "test challenge")

					pow := mockHandler.NewMockpow(ctrl)
					storage := mockHandler.NewMockquoteStorage(ctrl)

					pow.EXPECT().Generate(ctx, "192.168.1.1").Return("", errors.New("challenge error"))

					return req, resp, pow, storage
				},
				wantErr: true,
			},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				req, resp, pow, storage := tt.mockFunc(ctrl)

				h := handler.New(storage, pow)

				got, err := h.Handle(ctx, req)

				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, resp.ToString(), got.ToString())
				}
			})
		}
	})
}
