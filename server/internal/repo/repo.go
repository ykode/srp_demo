package repo

import "github.com/ykode/srp_demo/server/internal/domain"

type IdentityRepository interface {
	SaveIdentity(*domain.Identity) <-chan error
}

type SessionRepository interface {
	SaveSession(*domain.Session) <-chan error
}
