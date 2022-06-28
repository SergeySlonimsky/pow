package hashcash

import (
	"crypto/rand"
	"crypto/sha1" //nolint:gosec // it's ok for hashcash system
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const Version = 1

const defaultCounter = 0

const randBytesSize = 20

const zeroByte = 48

// Stamp is a struct for hashcash PoW algorithm to prevent DoS attacks.
// Contains all the necessary fields for hash generation and validation hashcash hash.
type Stamp struct {
	version   int
	zeroCount int
	date      int64
	resource  string
	rand      string
	counter   int
}

// New returns Stamp implementation of hashcash system with given resource and leading zero count.
// Returns error when can't generate random base64 string.
func New(resource string, zeroCount int, date int64) (*Stamp, error) {
	randBase64, err := generateRandBase64(randBytesSize)
	if err != nil {
		return nil, err
	}

	return &Stamp{
		version:   Version,
		zeroCount: zeroCount,
		date:      date,
		resource:  resource,
		rand:      randBase64,
		counter:   defaultCounter,
	}, nil
}

// FromString parses data string to hasshcash Stamp.
// String should be formatted as "version:zeroCount:date:resource:rand:counter".
func FromString(data string) (*Stamp, error) {
	parts := strings.Split(data, ":")
	if len(parts) != 6 {
		return nil, errors.New("invalid message format")
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil || version != Version {
		return nil, errors.New("invalid version")
	}

	zeroCount, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, errors.New("invalid zero count")
	}

	date, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, errors.New("invalid date")
	}

	counter, err := strconv.Atoi(parts[5])
	if err != nil {
		return nil, errors.New("invalid counter")
	}

	return &Stamp{
		version:   version,
		zeroCount: zeroCount,
		date:      date,
		resource:  parts[3],
		rand:      parts[4],
		counter:   counter,
	}, nil
}

// GenerateHash generates hash of the stamp, contains needed leading zero bytes.
// Returns nil if everything is ok and error if context deadline exceeded.
func (s *Stamp) GenerateHash(attempts int) error {
	for i := 1; i < attempts; i++ {
		s.counter = i
		hashStr := generateSha1(s.ToString())

		if checkHash(hashStr, s.zeroCount) {
			return nil
		}
	}

	return errors.New("max attempts exceeded")
}

// Verify verifies stamp
func (s *Stamp) Verify() bool {
	hash := generateSha1(s.ToString())

	return checkHash(hash, s.zeroCount)
}

// ToString returns string representation of the stamp data, separated with ":".
func (s *Stamp) ToString() string {
	return fmt.Sprintf("%d:%d:%d:%s:%s:%d", s.version, s.zeroCount, s.date, s.resource, s.rand, s.counter)
}

func (s *Stamp) GetRandValue() string {
	return s.rand
}

func (s *Stamp) GetResource() string {
	return s.resource
}

func generateRandBase64(randBytesSize int) (string, error) {
	bytes := make([]byte, randBytesSize)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func generateSha1(data string) string {
	hash := sha1.New() //nolint:gosec // it's ok for hashcash system
	hash.Write([]byte(data))

	return hex.EncodeToString(hash.Sum(nil))
}

func checkHash(hash string, zeroCount int) bool {
	// out of range protection
	if len(hash) < zeroCount {
		return false
	}

	leadingNumbers := hash[:zeroCount]

	for _, number := range leadingNumbers {
		if number != zeroByte {
			return false
		}
	}

	return true
}
