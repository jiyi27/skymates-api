package service

import "skymates-api/internal/repository"

type Services struct {
	UserService UserService
}

func NewServices(
	userRepo repository.UserRepository,
	// 其他仓库...
) *Services {
	return &Services{
		UserService: NewUserService(userRepo),
		// 初始化其他服务...
	}
}
