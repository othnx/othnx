package services

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/secondary"
	"aegis/pkg/apperrors"
	"aegis/pkg/plugins/providers"
)

type UserService struct {
	userRepository secondary.UserRepository
	config         entities.Config
}

func NewUserService(userRepository secondary.UserRepository, config entities.Config) *UserService {
	return &UserService{
		userRepository: userRepository,
		config:         config,
	}
}

// GetOrCreateUserIfAllowed handles the business logic for user creation and validation
// during OAuth authentication flows
func (s *UserService) GetOrCreateUserIfAllowed(userInfos *providers.UserInfos, authMethod string) (entities.User, error) {
	if userInfos.Email == "" {
		return entities.User{}, apperrors.ErrNoEmail
	}
	if userInfos.Name == "" {
		return entities.User{}, apperrors.ErrNoName
	}
	// Check if username already exists
	nameExists, err := s.userRepository.DoesNameExist(userInfos.Name)
	if err != nil {
		return entities.User{}, err
	}
	if nameExists {
		return entities.User{}, apperrors.ErrNameAlreadyExists
	}

	// Try to get existing user by email
	user, err := s.userRepository.GetUserByEmail(userInfos.Email)
	if err != nil && err.Error() != apperrors.ErrNoUser.Error() {
		return entities.User{}, err
	}

	// Create new user if doesn't exist
	if err != nil && err.Error() == apperrors.ErrNoUser.Error() {
		user, err = entities.NewUser(userInfos.Name, userInfos.Avatar, userInfos.Email, authMethod)
		if err != nil {
			return entities.User{}, err
		}
		roles := []entities.Role{entities.NewRole(user.ID, entities.RoleUser)}
		err = s.userRepository.CreateUser(user, roles)
		if err != nil {
			return entities.User{}, err
		}
		user, err = s.userRepository.GetUserByEmail(userInfos.Email)
		if err != nil {
			return entities.User{}, err
		}
	}

	// Validate user status
	if user.IsDeleted() {
		return entities.User{}, apperrors.ErrUserDeleted
	}
	if user.IsBlocked() {
		return entities.User{}, apperrors.ErrUserBlocked
	}
	if s.config.App.EarlyAdoptersOnly && !user.IsEarlyAdopter() {
		return entities.User{}, apperrors.ErrEarlyAdoptersOnly
	}

	return user, nil
}
