package pow_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/SergeySlonimsky/pow/internal/server/pow"
	mockPow "github.com/SergeySlonimsky/pow/internal/server/pow/mock"
)

const (
	ipAddr    = "192.168.1.1"
	stampData = "1:4:1656370862:172.21.0.4:FrZUho0yFjtWiiMonJTt55OFQ9k=:0"
)

func TestPoW_Generate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	tests := []struct {
		name     string
		mockFunc func(ctrl *gomock.Controller, ctx context.Context) *mockPow.Mockcache
		wantErr  bool
	}{
		{
			name: "valid",
			mockFunc: func(ctrl *gomock.Controller, ctx context.Context) *mockPow.Mockcache {
				cache := mockPow.NewMockcache(ctrl)
				cache.EXPECT().Add(ctx, "192.168.1.1", gomock.Any(), time.Minute*2).Return(nil)

				return cache
			},
			wantErr: false,
		},
		{
			name: "cache error",
			mockFunc: func(ctrl *gomock.Controller, ctx context.Context) *mockPow.Mockcache {
				cache := mockPow.NewMockcache(ctrl)
				cache.EXPECT().Add(ctx, "192.168.1.1", gomock.Any(), time.Minute*2).Return(errors.New("cache error"))

				return cache
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache := tt.mockFunc(ctrl, ctx)

			pw := pow.New(cache)
			result, err := pw.Generate(ctx, ipAddr)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, ipAddr, result)
			}
		})
	}
}

func TestPoW_Verify(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	tests := []struct {
		name     string
		mockFunc func(ctrl *gomock.Controller, ctx context.Context) *mockPow.Mockcache
		wantErr  bool
	}{
		{
			name: "valid",
			mockFunc: func(ctrl *gomock.Controller, ctx context.Context) *mockPow.Mockcache {
				cache := mockPow.NewMockcache(ctrl)
				cache.EXPECT().Get(ctx, "192.168.1.1").Return("FrZUho0yFjtWiiMonJTt55OFQ9k=", nil)

				return cache
			},
			wantErr: false,
		},
		{
			name: "cache error",
			mockFunc: func(ctrl *gomock.Controller, ctx context.Context) *mockPow.Mockcache {
				cache := mockPow.NewMockcache(ctrl)
				cache.EXPECT().Get(ctx, "192.168.1.1").Return("", errors.New("cache error"))

				return cache
			},
			wantErr: true,
		},
		{
			name: "rand number verify error",
			mockFunc: func(ctrl *gomock.Controller, ctx context.Context) *mockPow.Mockcache {
				cache := mockPow.NewMockcache(ctrl)
				cache.EXPECT().Get(ctx, "192.168.1.1").Return("invalid rand string", nil)

				return cache
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache := tt.mockFunc(ctrl, ctx)

			pw := pow.New(cache)
			err := pw.Verify(ctx, ipAddr, stampData)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
