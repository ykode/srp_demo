package handlers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ykode/srp_demo/server/internal/domain"
	"math/big"
)

type SessionPayload struct {
	session *domain.Session
	salt    []byte
}

func (s SessionPayload) MarshalJSON() ([]byte, error) {

	id := s.session.ID()
	B := base64.StdEncoding.EncodeToString(s.session.B().Bytes())
	salt := base64.StdEncoding.EncodeToString(s.salt)

	j := struct {
		B         string    `json:"B"`
		SessionID uuid.UUID `json:"session_id"`
		Salt      string    `json:"salt"`
	}{
		B:         B,
		SessionID: id,
		Salt:      salt,
	}

	return json.Marshal(j)
}

type SessionAnswer struct {
	session *domain.Session
	m_s     *big.Int
}

func (s SessionAnswer) MarshalJSON() ([]byte, error) {
	id := s.session.ID()
	m_s := base64.StdEncoding.EncodeToString(s.m_s.Bytes())

	j := struct {
		SessionID uuid.UUID `json:"session_id"`
		M_s       string    `json:"M_s"`
	}{
		SessionID: id,
		M_s:       m_s,
	}

	return json.Marshal(j)
}
