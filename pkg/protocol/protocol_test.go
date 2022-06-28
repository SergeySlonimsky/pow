package protocol_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SergeySlonimsky/pow/pkg/protocol"
)

func TestParseFromReader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		str     string
		want    protocol.Message
		wantErr bool
	}{
		{
			name:    "without body",
			str:     "resource\n",
			want:    protocol.NewMessage(protocol.TypeResource, ""),
			wantErr: false,
		},
		{
			name:    "with body",
			str:     "challenge|test challenge\n",
			want:    protocol.NewMessage(protocol.TypeChallenge, "test challenge"),
			wantErr: false,
		},
		{
			name:    "invalid type",
			str:     "invalid|test challenge",
			want:    protocol.Message{},
			wantErr: true,
		},
		{
			name:    "without EOF",
			str:     "challenge|test challenge",
			want:    protocol.Message{},
			wantErr: true,
		},
		{
			name:    "ignore other EOFs",
			str:     "challenge|test challenge\nother date",
			want:    protocol.NewMessage(protocol.TypeChallenge, "test challenge"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := strings.NewReader(tt.str)

			got, err := protocol.ParseFromReader(r)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.GetType(), got.GetType())
				assert.Equal(t, tt.want.GetBody(), got.GetBody())
			}
		})
	}
}

func TestMessage_ToString(t *testing.T) {
	msg := protocol.NewMessage(protocol.TypeResource, "testBody")

	assert.Equal(t, fmt.Sprintf("%s|%s", protocol.TypeResource, "testBody"), msg.ToString())
}
