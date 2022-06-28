package hashcash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checkHash(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		hash      string
		zeroCount int
		want      bool
	}{
		{
			name:      "valid hash with 4 leading zero",
			hash:      "0000aek01f",
			zeroCount: 4,
			want:      true,
		},
		{
			name:      "invalid hash with 4 ending zero",
			hash:      "aek01f0000",
			zeroCount: 3,
			want:      false,
		},
		{
			name:      "invalid leading zero count",
			hash:      "0aef00k0f",
			zeroCount: 2,
			want:      false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, checkHash(tt.hash, tt.zeroCount))
		})
	}
}
