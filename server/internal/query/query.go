package query

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ykode/srp_demo/server/internal/domain"
)

var (
	ErrNotFound = errors.New("Not Found")
)

type Result struct {
	Result interface{}
	Err    error
}

func (r *Result) IsError() bool {
	return r.Err != nil
}

func (r *Result) IdentityResult() (*domain.Identity, bool) {
	v, ok := r.Result.(*domain.Identity)

	return v, ok
}

func (r *Result) SessionResult() (*domain.Session, bool) {
	v, ok := r.Result.(*domain.Session)

	return v, ok
}

type QueryFunc = func(...interface{}) (interface{}, error)

func Query(qfn QueryFunc, params ...interface{}) <-chan Result {
	c := make(chan Result)

	go func() {
		r, err := qfn(params...)

		c <- Result{Result: r, Err: err}

		close(c)
	}()

	return c
}

type IdentityQuery interface {
	FindIdentityByUserName(username string) <-chan Result
}

type SessionQuery interface {
	FindSessionById(sessionId uuid.UUID) <-chan Result
}
