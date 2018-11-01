package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo"
	"github.com/ykode/srp_demo/server/internal/domain"
	"github.com/ykode/srp_demo/server/internal/query"
	"github.com/ykode/srp_demo/server/internal/repo"
	"net/http"
)

type IdentityHandler struct {
	repo  repo.IdentityRepository
	query query.IdentityQuery
}

func NewIdentityHandler(idRepo repo.IdentityRepository, idQuery query.IdentityQuery) *IdentityHandler {
	return &IdentityHandler{
		repo:  idRepo,
		query: idQuery,
	}
}

func (h *IdentityHandler) Mount(g *echo.Group) {
	g.POST("", h.RegisterIdentity)
	g.POST("/", h.RegisterIdentity)
}

func (h *IdentityHandler) RegisterIdentity(c echo.Context) error {
	userName := c.FormValue("user_name")
	salt_base64 := c.FormValue("salt")
	v_base64 := c.FormValue("v")

	salt, err := base64.StdEncoding.DecodeString(salt_base64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid salt")
	}

	v, err := base64.StdEncoding.DecodeString(v_base64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid verifier value")
	}

	id, err := domain.NewIdentity(userName, salt, v)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.repo.SaveIdentity(id)

	fmt.Printf("Storage: \n%+v\n", h.repo)

	return c.String(http.StatusCreated, "Created")
}
