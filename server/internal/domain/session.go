package domain

import (
	"crypto/sha256"
	"github.com/google/uuid"
	"golang.org/x/crypto/hkdf"
	"io"
	"math/big"
)

type SessionState string

const (
	SessionStateChallengeSent SessionState = "CHALLENGE_SENT"
	SessionStateCompleted     SessionState = "COMPLETED"
)

type Session struct {
	id        uuid.UUID
	salt      []byte
	masterKey []byte
	state     SessionState

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

func BuildSession(id uuid.UUID, salt, masterKey, v []byte, state SessionState) (*Session, error) {

	b, err := cryptrand(bitlen / 8)

	if err != nil {
		return nil, err
	}

	return &Session{
		id:        id,
		salt:      salt,
		masterKey: masterKey,
		state:     state,
		b:         b,
		v:         new(big.Int).SetBytes(v),
	}, nil
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
		id:    sessionId,
		b:     b,
		state: SessionStateChallengeSent,
	}, nil
}

func (s Session) ID() uuid.UUID {
	return s.id
}

func (s Session) State() SessionState {
	return s.state
}

func (s Session) B() *big.Int {
	return new(big.Int).Mod(new(big.Int).Add(new(big.Int).Mul(k, s.v), new(big.Int).Exp(g, s.b, N)), N)
}

func (s *Session) GenerateKey(A *big.Int) {
	B := s.B()
	u := calculateHashBigInt(A, B)
	s.masterKey = new(big.Int).Exp(new(big.Int).Mul(A, new(big.Int).Exp(s.v, u, N)), s.b, N).Bytes()
	s.state = SessionStateCompleted
}
