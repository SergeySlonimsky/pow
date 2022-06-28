package pow

import (
	"context"
	"errors"
	"time"

	"github.com/SergeySlonimsky/pow/pkg/hashcash"
)

//go:generate mockgen -source=./pow.go -destination=./mock/pow_mock.go

const defaultZeroCount = 4

const defaultStampTTL = time.Minute * 2

type cache interface {
	Add(ctx context.Context, key, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type PoW struct {
	cache cache
}

func New(cache cache) *PoW {
	return &PoW{
		cache: cache,
	}
}

// Generate generates Proof of Work string for a client with given resource.
func (p *PoW) Generate(ctx context.Context, resource string) (string, error) {
	stamp, err := hashcash.New(resource, defaultZeroCount, time.Now().Unix())
	if err != nil {
		return "", err
	}

	if err := p.cache.Add(ctx, resource, stamp.GetRandValue(), defaultStampTTL); err != nil {
		return "", err
	}

	return stamp.ToString(), nil
}

// Verify verifies Proof of Work string from a client by given resource and PoW algorithm.
func (p *PoW) Verify(ctx context.Context, resource, data string) error {
	stamp, err := hashcash.FromString(data)
	if err != nil {
		return err
	}

	randValue, err := p.cache.Get(ctx, resource)
	if err != nil {
		return err
	}

	if randValue != stamp.GetRandValue() {
		return errors.New("invalid rand value or challenge timeout exceeded")
	}

	if !stamp.Verify() {
		return errors.New("invalid challenge stamp")
	}

	return nil
}
