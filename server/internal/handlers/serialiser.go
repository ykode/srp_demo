package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ykode/srp_demo/server/internal/domain"
	"math/big"
)

type SessionPayload domain.Session

func (s SessionPayload) MarshalJSON() ([]byte, error) {

	id := domain.Session(s).ID()
	B := domain.Session(s).B()

	j := struct {
		B         *big.Int  `json:"B"`
		SessionID uuid.UUID `json:"session_id"`
	}{
		B:         B,
		SessionID: id,
	}

	return json.Marshal(j)
}
