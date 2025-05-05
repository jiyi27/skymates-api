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
	"strconv"
)

// UserService 定义用户相关的业务逻辑接口
type UserService interface {
	Register(registerDto v1.RegisterDto) (*model.User, error)
	Login(loginDto v1.LoginDto) (*model.User, string, error)
	GetUserById(id int64) (*model.User, error)
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
// user_service.go
func (s *userService) Register(registerDto v1.RegisterDto) (*model.User, error) {
	exists, err := s.userRepository.CheckExists(repository.QueryByUsername, registerDto.Username)
	if err != nil {
		log.Printf("UserService.Register: failed to check username exists: %v", err)
		return nil, servererrors.NewInternalError("检查用户名是否存在失败", err)
	}
	if exists {
		return nil, servererrors.NewAlreadyExistsError("用户名已存在", nil)
	}

	exists, err = s.userRepository.CheckExists(repository.QueryByEmail, registerDto.Email)
	if err != nil {
		log.Printf("UserService.Register: failed to check email exists: %v", err)
		return nil, servererrors.NewInternalError("检查邮箱是否存在失败", err)
	}
	if exists {
		return nil, servererrors.NewAlreadyExistsError("邮箱已存在", nil)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerDto.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("UserService.Register: failed to hash password: %v", err)
		return nil, servererrors.NewInternalError("密码加密失败", err)
	}

	user := &model.User{
		Username: registerDto.Username,
		Password: string(hashedPassword),
		Email:    registerDto.Email,
	}
	if err := s.userRepository.Create(user); err != nil {
		log.Printf("UserService.Register: failed to create user: %v", err)
		return nil, servererrors.NewInternalError("创建用户失败", err)
	}

	return user, nil
}

// Login 处理用户登录业务逻辑
// 成功时返回用户信息和JWT令牌，失败时返回错误
func (s *userService) Login(loginDto v1.LoginDto) (*model.User, string, error) {
	// 1. 查询用户
	user, err := s.userRepository.GetUserBy(repository.QueryByEmail, loginDto.Email)
	if err != nil {
		// 如果是未找到，映射成 NotFoundError
		var se *servererrors.ServerError
		if errors.As(err, &se) && se.Kind == servererrors.KindNotFound {
			return nil, "", servererrors.NewNotFoundError("用户不存在", nil)
		}
		// 其他视为内部错误
		log.Printf("UserService.Login: failed to get user by email: %v", err)
		return nil, "", servererrors.NewInternalError("获取用户失败", err)
	}

	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password)); err != nil {
		// 密码不匹配当作 Unauthorized
		return nil, "", servererrors.NewUnauthorizedError("凭证无效", nil)
	}

	// 3. 生成 JWT
	jwtToken, err := auth.GenerateJwtToken(user)
	if err != nil {
		log.Printf("UserService.Login: failed to generate jwt token: %v", err)
		return nil, "", servererrors.NewInternalError("生成令牌失败", err)
	}

	return user, jwtToken, nil
}

func (s *userService) GetUserById(id int64) (*model.User, error) {
	user, err := s.userRepository.GetUserBy(repository.QueryByID, strconv.FormatInt(id, 10))
	if err != nil {
		var se *servererrors.ServerError
		if errors.As(err, &se) && se.Kind == servererrors.KindNotFound {
			return nil, servererrors.NewNotFoundError("用户不存在", nil)
		}
		log.Printf("UserService.GetUserById: failed to get user by id: %v", err)
		return nil, servererrors.NewInternalError("获取用户失败", err)
	}

	return user, nil
}
