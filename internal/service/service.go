package service

import "skymates-api/internal/repository"

type Services struct {
	UserService UserService
	TermService TermService
}

func NewServices(
	userRepository repository.UserRepository,
	termRepository repository.TermRepository,
) *Services {
	return &Services{
		UserService: NewUserService(userRepository),
		TermService: NewTermService(termRepository),
	}
}
