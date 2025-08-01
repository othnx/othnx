package middlewares

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/primary"
	"aegis/pkg/apperrors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OAuthMiddlewaresInterface interface {
	CheckAuthEnabled(next echo.HandlerFunc) echo.HandlerFunc
}

type OAuthMiddlewares struct {
	Config  entities.Config
	Service primary.OAuthUseCasesInterface
}

var _ OAuthMiddlewaresInterface = (*OAuthMiddlewares)(nil)

func NewOAuthMiddlewares(c entities.Config, s primary.OAuthUseCasesInterface) *OAuthMiddlewares {
	return &OAuthMiddlewares{
		Config:  c,
		Service: s,
	}
}

func (m OAuthMiddlewares) CheckAuthEnabled(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		enabled := m.Service.CheckAuthEnabled()
		if !enabled {
			return c.JSON(http.StatusForbidden, map[string]string{"error": apperrors.ErrAuthMethodNotEnabled.Error()})
		}
		return next(c)
	}
}
