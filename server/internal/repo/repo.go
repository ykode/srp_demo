package repo

import "github.com/ykode/srp_demo/server/internal/domain"

type IdentityRepository interface {
	Save(*domain.Identity) <-chan error
}

type SessionRepository interface {
	Save(*domain.Session) <-chan error
}
