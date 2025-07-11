package primary

import "aegis/internal/domain/entities"

type UseCasesForHandlers interface {
	GetSession(accessToken string) (entities.Session, error)
	Logout(refreshToken string) (*entities.TokenPair, error)
}

type UseCasesForMiddlewares interface {
	CheckAndRefreshToken(accessToken, refreshToken string, forceRefresh bool) (*entities.TokenPair, error)
}

type UseCasesInterface interface {
	UseCasesForHandlers
	UseCasesForMiddlewares
}
