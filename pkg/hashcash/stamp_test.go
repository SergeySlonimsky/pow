package hashcash_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/SergeySlonimsky/pow/pkg/hashcash"
)

func TestStamp_Generate_Verify(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		zeroCount int
	}{
		{
			name:      "1 leading zero",
			zeroCount: 1,
		},
		{
			name:      "2 leading zeros",
			zeroCount: 2,
		},
		{
			name:      "3 leading zeros",
			zeroCount: 3,
		},
		{
			name:      "4 leading zeros",
			zeroCount: 4,
		},
		{
			name:      "5 leading zeros",
			zeroCount: 5,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stamp, err := hashcash.New("test", tt.zeroCount, time.Now().Unix())
			assert.NoError(t, err)

			assert.NoError(t, stamp.GenerateHash(9999999))
			assert.True(t, stamp.Verify())
		})
	}
}

func TestStamp_ToString(t *testing.T) {
	now := time.Now().Unix()
	stamp, err := hashcash.New("testResource", 4, now)

	assert.NoError(t, err)
	assert.Equal(t, "testResource", stamp.GetResource())
}
