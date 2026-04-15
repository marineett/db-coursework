package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
)

type IAuthService interface {
	Authorize(auth_data types.AuthData) (types.AuthVerdict, error)
	CheckLogin(username string) (bool, error)
}

type AuthService struct {
	AuthRepository data_base.IAuthRepository
}

func CreateAuthService(authRepository data_base.IAuthRepository) *AuthService {
	return &AuthService{
		AuthRepository: authRepository,
	}
}

func (s *AuthService) Authorize(auth_data types.AuthData) (types.AuthVerdict, error) {
	return s.AuthRepository.Authorize(auth_data)
}

func (s *AuthService) CheckLogin(username string) (bool, error) {
	return s.AuthRepository.CheckLogin(username)
}
