package domain

import "github.com/google/uuid"

type Session struct {
	id           uuid.UUID
	masterKey    []byte
	keyIteration int
}

func NewSession() (*Session, error) {

	sessionId, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	return &Session{
		id: sessionId,
	}, nil
}

func (s *Session) ID() uuid.UUID {
	return s.id
}

func (s *Session) KeyIteration() int {
	return s.keyIteration
}

func (s *Session) Initialised() bool {
	return len(s.masterKey) != 0
}
