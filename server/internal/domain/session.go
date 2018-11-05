package domain

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
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
	masterKey []byte
	keys      [4]*big.Int
	state     SessionState

	b  *big.Int
	_B *big.Int
	_A *big.Int
	v  *big.Int
}

const keyinfo = "SRP Demo Key Information"
const bitlen = 1024

func hkdfFromKey(salt []byte, ikm []byte, iteration int) [][]byte {

	extractor := hkdf.Extract(sha256.New, ikm, salt)
	fmt.Printf("Extractor: %s\n", base64.StdEncoding.EncodeToString(extractor))

	hkdf := hkdf.New(sha256.New, ikm, salt, []byte(keyinfo))
	okm := make([]byte, 16)
	out := make([][]byte, iteration)

	for i := 0; i < iteration; i += 1 {
		io.ReadFull(hkdf, okm)
		fmt.Printf("K[%d]: 0x%x\n", i, new(big.Int).SetBytes(okm))
		out[i] = okm
	}

	return out
}

func BuildSession(id uuid.UUID, masterKey []byte, v, b, A *big.Int, state SessionState) (*Session, error) {

	B := new(big.Int).Mod(new(big.Int).Add(new(big.Int).Mul(k, v), new(big.Int).Exp(g, b, N)), N)

	u := calculateHashBigInt(A, B)

	var keys [4]*big.Int

	keysBytes := hkdfFromKey(u.Bytes(), masterKey, 4)

	for i := 0; i < 4; i += 1 {
		keys[i] = new(big.Int).SetBytes(keysBytes[i])
	}

	return &Session{
		id:        id,
		masterKey: masterKey,
		state:     state,
		b:         b,
		_B:        B,
		v:         v,
		_A:        A,
		keys:      keys,
	}, nil
}

func NewSession(v *big.Int) (*Session, error) {

	if v == nil {
		return nil, errors.New("v cannot be empty")
	}

	sessionId, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	b, err := cryptrand(bitlen / 8)

	if err != nil {
		return nil, err
	}

	B := new(big.Int).Mod(new(big.Int).Add(new(big.Int).Mul(k, v), new(big.Int).Exp(g, b, N)), N)

	return &Session{
		id:    sessionId,
		b:     b,
		state: SessionStateChallengeSent,
		v:     v,
		_B:    B,
	}, nil
}

func (s Session) ID() uuid.UUID {
	return s.id
}

func (s Session) State() SessionState {
	return s.state
}

func (s Session) SmallB() *big.Int {
	return s.b
}

func (s Session) B() *big.Int {
	return s._B
}

func (s Session) A() *big.Int {
	return s._A
}

func (s Session) V() *big.Int {
	return s.v
}

func (s Session) MasterKey() []byte {
	return s.masterKey
}

func (s *Session) GenerateKey(A *big.Int) error {

	if A.Cmp(new(big.Int)) == 0 {
		return errors.New("A cannot be Zero")
	}

	B := s._B
	u := calculateHashBigInt(A, B)
	S_s := new(big.Int).Exp(new(big.Int).Mul(A, new(big.Int).Exp(s.v, u, N)), s.b, N)
	fmt.Printf("S_s: 0x%x\n", S_s)
	s.masterKey = S_s.Bytes()

	keys := hkdfFromKey(u.Bytes(), s.masterKey, len(s.keys))

	for i := 0; i < len(s.keys); i += 1 {
		s.keys[i] = new(big.Int).SetBytes(keys[i])
	}

	s._A = A
	return nil
}

func (s *Session) VerifyClient(M1_c *big.Int) (*big.Int, bool) {

	M1_s := calculateHashBigInt(s.keys[0], new(big.Int).Exp(s._A, s._B, N))
	/*
		fmt.Printf("v: 0x%x\nA:%s\nB:%s\nS_s: 0x%x\nK1_s: 0x%x\nM1_s: %s",
			s.v,
			base64.StdEncoding.EncodeToString(s._A.Bytes()),
			base64.StdEncoding.EncodeToString(s._B.Bytes()),
			s.masterKey, s.keys[0], base64.StdEncoding.EncodeToString(M1_s.Bytes()))
	*/
	if eq := M1_s.Cmp(M1_c) == 0; !eq {
		return nil, false
	} else {
		s.state = SessionStateCompleted
		return calculateHashBigInt(s.keys[0], new(big.Int).Exp(s._A, M1_c, N)), eq
	}
}
