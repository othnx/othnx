package usecases

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/primary"
	"aegis/internal/domain/ports/secondary"
	"aegis/internal/domain/services"
	"aegis/pkg/apperrors"
	"aegis/pkg/plugins/providers"
	"aegis/pkg/tokengen"
	"time"
)

type OAuthUseCases struct {
	Config                 entities.Config
	Provider               providers.OAuthProviderInterface
	UserRepository         secondary.UserRepository
	RefreshTokenRepository secondary.RefreshTokenRepository
	StateRepository        secondary.StateRepository
	UserService            *services.UserService
	TokenService           *services.TokenService
}

var _ primary.OAuthUseCasesInterface = (*OAuthUseCases)(nil)

func NewOAuthUseCases(
	c entities.Config,
	p providers.OAuthProviderInterface,
	userRepository secondary.UserRepository,
	refreshTokenRepository secondary.RefreshTokenRepository,
	stateRepository secondary.StateRepository,
) *OAuthUseCases {
	userService := services.NewUserService(userRepository, c)
	tokenService := services.NewTokenService(refreshTokenRepository, c)
	return &OAuthUseCases{
		Config:                 c,
		Provider:               p,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
		StateRepository:        stateRepository,
		UserService:            userService,
		TokenService:           tokenService,
	}
}

func (s OAuthUseCases) CheckAuthEnabled() bool {
	return s.Provider.IsEnabled()
}

func (s *OAuthUseCases) GetAuthURL(redirectUri string) (string, error) {
	state, err := tokengen.Generate("state_", 13)
	if err != nil {
		return "", err
	}
	if err := s.StateRepository.CreateState(entities.NewState(state)); err != nil {
		return "", err
	}
	redirectURL := s.Provider.GetOauthRedirectURL(state)
	return redirectURL, nil
}

func (s OAuthUseCases) ExchangeCode(code, state string) (*entities.TokenPair, error) {
	serverState, err := s.StateRepository.GetAndDeleteState(state)
	if err != nil {
		return nil, apperrors.ErrInvalidState
	}
	if serverState.IsExpired() {
		return nil, apperrors.ErrInvalidState
	}

	userInfos, err := s.Provider.ExchangeCodeForUserInfos(code, state)
	if err != nil {
		return nil, err
	}

	user, err := s.UserService.GetOrCreateUserIfAllowed(userInfos, s.Provider.GetName())
	if err != nil {
		return nil, err
	}

	if user.AuthMethod != s.Provider.GetName() {
		return nil, apperrors.ErrWrongAuthMethod
	}

	// todo device-id: pass one, since one session per device is allowed
	accessToken, atExpiresAt, newRefreshToken, rtExpiresAt, err := s.TokenService.GenerateTokensForUser(user, "device-id")
	if err != nil {
		return nil, err
	}

	result := &entities.TokenPair{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  time.Unix(atExpiresAt, 0),
		RefreshToken:          newRefreshToken,
		RefreshTokenExpiresAt: time.Unix(rtExpiresAt, 0),
	}

	return result, nil
}
