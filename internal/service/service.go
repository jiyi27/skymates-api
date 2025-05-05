package service

import "skymates-api/internal/repository"

type Services struct {
	UserService UserService
}

func NewServices(
	userRepository repository.UserRepository,
	// 其他仓库...
) *Services {
	return &Services{
		UserService: NewUserService(userRepository),
		// 初始化其他服务...
	}
}
