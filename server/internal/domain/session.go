package domain

import (
	"crypto/sha256"
	"github.com/google/uuid"
	"golang.org/x/crypto/hkdf"
	"io"
	"math/big"
)

type Session struct {
	id                  uuid.UUID
	salt                []byte
	masterKey           []byte
	currentKeyIteration int
	currentKey          []byte

	b *big.Int
	v *big.Int
}

const keyinfo = "SRP Demo Key Information"
const bitlen = 1024

func hkdfFromKey(salt []byte, ikm []byte, iteration int) []byte {
	hkdf := hkdf.New(sha256.New, ikm, salt, []byte(keyinfo))
	okm := make([]byte, 16)

	for i := 0; i < iteration; i += 1 {
		io.ReadFull(hkdf, okm)
	}

	return okm
}

func NewSession(salt []byte, v *big.Int) (*Session, error) {

	sessionId, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}
	b, err := cryptrand(bitlen / 8)

	if err != nil {
		return nil, err
	}

	return &Session{
		id: sessionId,
		b:  b,
	}, nil
}

func (s *Session) ID() uuid.UUID {
	return s.id
}

func (s *Session) Initialised() bool {
	return len(s.masterKey) != 0
}

func (s *Session) GenerateNewKey() []byte {
	s.currentKeyIteration += 1

	s.currentKey = hkdfFromKey(s.salt, s.masterKey, s.currentKeyIteration)

	return s.currentKey
}

func (s *Session) CurrentIteration() int {
	return s.currentKeyIteration
}

func (s *Session) CurrentKey() []byte {
	return s.currentKey
}

func (s *Session) GenerateKey(A *big.Int) error {
	B := new(big.Int).Mod(new(big.Int).Add(new(big.Int).Mul(k, s.v), new(big.Int).Exp(g, s.b, N)), N)
	u := calculateHashBigInt(A, B)
	s.masterKey = new(big.Int).Exp(new(big.Int).Mul(A, new(big.Int).Exp(s.v, u, N)), s.b, N).Bytes()
	return nil
}
