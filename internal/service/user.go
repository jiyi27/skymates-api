package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	servererrors "skymates-api/errors"
	v1 "skymates-api/internal/dto/v1"
	"skymates-api/internal/model"
	"skymates-api/internal/repository"
	"skymates-api/pkg/auth"
)

// UserService 定义用户相关的业务逻辑接口
type UserService interface {
	Register(req v1.RegisterDto) (*model.User, error)
	Login(req v1.LoginDto) (*model.User, string, error)
}

// userService 实现 UserService 接口
type userService struct {
	userRepository repository.UserRepository
}

// NewUserService 创建 UserService 实例
func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

// Register 处理用户注册业务逻辑
// 成功时返回创建的用户，失败时返回错误
func (s *userService) Register(req v1.RegisterDto) (*model.User, error) {
	// 检查用户名是否存在
	exists, err := s.userRepository.CheckExists(repository.QueryByUsername, req.Username)
	if err != nil {
		log.Printf("UserService.Register: failed to check username exists: %v", err)
		return nil, errors.New("internal server error")
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否存在
	exists, err = s.userRepository.CheckExists(repository.QueryByEmail, req.Email)
	if err != nil {
		log.Printf("UserService.Register: failed to check email exists: %v", err)
		return nil, errors.New("internal server error")
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("UserService.Register: failed to hash password: %v", err)
		return nil, errors.New("internal server error")
	}

	// 创建用户
	user := &model.User{
		Username:       req.Username,
		HashedPassword: string(hashedPassword),
		Email:          req.Email,
	}

	// 保存用户到数据库
	if err := s.userRepository.Create(user); err != nil {
		log.Printf("UserService.Register: failed to create user: %v", err)
		return nil, errors.New("internal server error")
	}

	return user, nil
}

// Login 处理用户登录业务逻辑
// 成功时返回用户信息和JWT令牌，失败时返回错误
func (s *userService) Login(req v1.LoginDto) (*model.User, string, error) {
	// 查询用户
	user, err := s.userRepository.GetUserBy(repository.QueryByEmail, req.Email)
	if err != nil {
		var serverErr *servererrors.ServerError
		if errors.As(err, &serverErr) {
			if serverErr.Kind == servererrors.KindNotFound {
				return nil, "", errors.New("user not found")
			}
		}
		log.Printf("UserService.Login: failed to get user by email: %v", err)
		return nil, "", errors.New("internal server error")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// 生成JWT令牌
	jwtToken, err := auth.GenerateJwtToken(user)
	if err != nil {
		log.Printf("UserService.Login: failed to generate jwt token: %v", err)
		return nil, "", errors.New("internal server error")
	}

	return user, jwtToken, nil
}
