package query

import (
	"github.com/google/uuid"
	"github.com/ykode/srp_demo/server/internal/domain"
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

type IdentityQuery interface {
	FindIdentityByUserName(username string) <-chan Result
}

type SessionQuery interface {
	FindSessionbyId(sessionId uuid.UUID) <-chan Result
}
