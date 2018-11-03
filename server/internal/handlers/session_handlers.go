package handlers

import (
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/ykode/srp_demo/server/internal/domain"
	"github.com/ykode/srp_demo/server/internal/query"
	"github.com/ykode/srp_demo/server/internal/repo"
	"math/big"
	"net/http"
)

type SessionAction string

const (
	SessionActionStart  SessionAction = "start_session"
	SessionActionAnswer SessionAction = "answer"
)

type SessionHandler struct {
	repo    repo.SessionRepository
	query   query.SessionQuery
	idQuery query.IdentityQuery
}

func NewSessionHandler(sessionRepo repo.SessionRepository, sessionQuery query.SessionQuery,
	idQuery query.IdentityQuery) *SessionHandler {

	return &SessionHandler{
		repo:    sessionRepo,
		query:   sessionQuery,
		idQuery: idQuery,
	}
}

func (h *SessionHandler) Mount(g *echo.Group) {
	g.POST("", h.HandleSession)
	g.POST("/", h.HandleSession)
}

func extractBytesFromParam(c echo.Context, paramName string) ([]byte, error) {
	b64 := c.FormValue(paramName)

	if len(b64) == 0 {
		return nil, errors.New("Bad Value")
	}

	paramBytes, err := base64.StdEncoding.DecodeString(b64)

	if err != nil {
		return nil, err
	}

	return paramBytes, nil
}

func extractBigIntFromParam(c echo.Context, paramName string) (*big.Int, error) {

	paramBytes, err := extractBytesFromParam(c, paramName)

	if err != nil {
		return nil, err
	}

	return new(big.Int).SetBytes(paramBytes), nil
}

func (h *SessionHandler) StartSession(c echo.Context) error {

	userName := c.FormValue("user_name")

	if len(userName) == 0 {
		return c.String(http.StatusBadRequest, "UserName cannot be empty")
	}

	r := <-h.idQuery.FindIdentityByUserName(userName)

	if r.IsError() {
		return c.String(http.StatusBadRequest, r.Err.Error())
	}

	identity, ok := r.IdentityResult()

	if !ok {
		return c.String(http.StatusBadRequest, "Error Identity Type")
	}

	vb := identity.Verifier()

	v := new(big.Int).SetBytes(vb)

	A, err := extractBigIntFromParam(c, "A")

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	session, err := domain.NewSession(v)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = session.GenerateKey(A)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = <-h.repo.SaveSession(session)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, SessionPayload{session: session, salt: identity.Salt()})

}

func (h *SessionHandler) AnswerChallenge(c echo.Context) error {
	sessionIdStr := c.FormValue("session_id")
	mClientStr := c.FormValue("m_client")

	if len(sessionIdStr) == 0 {
		return c.String(http.StatusBadRequest, "Session Id cannot be empty")
	}

	if len(mClientStr) == 0 {
		return c.String(http.StatusBadRequest, "Client Verifier cannot be empty")
	}

	sessionId, err := uuid.Parse(sessionIdStr)

	if err != nil {
		return c.String(http.StatusBadRequest, "Session Error Parse")
	}

	r := <-h.query.FindSessionById(sessionId)

	if r.IsError() {
		return c.String(http.StatusBadRequest, r.Err.Error())
	}

	session, ok := r.SessionResult()

	if !ok {
		return c.String(http.StatusBadRequest, "Error System Type")
	}

	M_c, err := extractBigIntFromParam(c, "m_client")

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	M_s, ok := session.VerifyClient(M_c)

	if !ok {
		return c.String(http.StatusUnauthorized, "Client Verification Failed")
	}

	err = <-h.repo.SaveSession(session)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, SessionAnswer{session: session, m_s: M_s})
}

func (h *SessionHandler) HandleSession(c echo.Context) error {
	action := c.FormValue("action")

	switch SessionAction(action) {
	case SessionActionStart:
		return h.StartSession(c)
	case SessionActionAnswer:
		return h.AnswerChallenge(c)
	default:
		return c.String(http.StatusBadRequest, "Invalid action")
	}
}
