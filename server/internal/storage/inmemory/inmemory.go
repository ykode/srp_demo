package inmemory

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ykode/srp_demo/server/internal/domain"
	"github.com/ykode/srp_demo/server/internal/query"
	"sync"
)

type InMemoryIdentityStorage struct {
	identities map[string]domain.Identity

	lock sync.RWMutex
}

func NewInMemoryIdentityStorage() *InMemoryIdentityStorage {
	return &InMemoryIdentityStorage{
		identities: make(map[string]domain.Identity),
	}
}

func (s *InMemoryIdentityStorage) FindIdentityByUserName(username string) <-chan query.Result {

	queryFunc := func(uname string) (*domain.Identity, error) {
		v, ok := s.identities[uname]
		if !ok {
			return nil, errors.New("Not Found")
		}

		return &v, nil
	}

	result := make(chan query.Result)

	go func() {
		s.lock.RLock()
		defer s.lock.RUnlock()

		v, err := queryFunc(username)

		result <- query.Result{Result: v, Err: err}

		close(result)

	}()

	return result
}

func (s *InMemoryIdentityStorage) SaveIdentity(id *domain.Identity) <-chan error {
	result := make(chan error)

	go func() {
		s.lock.Lock()
		defer s.lock.Unlock()

		s.identities[id.UserName()] = *id

		result <- nil

		close(result)
	}()

	return result
}

type InMemorySessionStorage struct {
	sessions map[uuid.UUID]domain.Session

	lock sync.RWMutex
}

func NewInMemorySessionStorage() *InMemorySessionStorage {
	return &InMemorySessionStorage{
		sessions: make(map[uuid.UUID]domain.Session),
	}
}

func (s *InMemorySessionStorage) FindSessionbyId(id uuid.UUID) <-chan query.Result {
	queryFunc := func(id uuid.UUID) (*domain.Session, error) {
		v, ok := s.sessions[id]
		if !ok {
			return nil, errors.New("Not Found")
		}

		return &v, nil
	}

	result := make(chan query.Result)

	go func() {
		s.lock.RLock()
		defer s.lock.RUnlock()

		v, err := queryFunc(id)

		result <- query.Result{Result: v, Err: err}

		close(result)
	}()

	return result
}

func (s *InMemorySessionStorage) SaveSession(sess *domain.Session) <-chan error {
	result := make(chan error)

	go func() {
		s.lock.Lock()
		defer s.lock.Unlock()

		s.sessions[sess.ID()] = *sess

		result <- nil

		close(result)
	}()

	return result
}
