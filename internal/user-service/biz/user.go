package biz

import (
	"context"
	"time"

	"student/internal/pkg/password"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// TimeFormat 时间格式
const TimeFormat = "2006-01-02 15:04:05"

// User 用户模型
type User struct {
	ID        uint
	Username  string
	Email     string
	Phone     string
	Password  string
	Status    int
	Age       int
	Avatar    string
	CreatedAt *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`

	// 格式化后的时间字符串，避免每次转换
	CreatedAtStr string `gorm:"-" json:"created_at_str,omitempty"` // gorm:"-" 表示不映射到数据库
	UpdatedAtStr string `gorm:"-" json:"updated_at_str,omitempty"` // gorm:"-" 表示不映射到数据库

	// 角色信息
	Roles []string `gorm:"-" json:"roles,omitempty"`
}

// FormatTimeFields 格式化时间字段
func (u *User) FormatTimeFields() {
	if u.CreatedAt != nil {
		u.CreatedAtStr = u.CreatedAt.Format(TimeFormat)
	}
	if u.UpdatedAt != nil {
		u.UpdatedAtStr = u.UpdatedAt.Format(TimeFormat)
	}
}

// FormatUserTimeFieldsBatch 批量格式化用户时间字段
func FormatUserTimeFieldsBatch(users []*User) {
	for _, user := range users {
		user.FormatTimeFields()
	}
}

// UserRepo 用户仓储接口
type UserRepo interface {
	GetUser(ctx context.Context, id int32) (*User, error)
	ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersReply, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, id int32) error
	Login(ctx context.Context, loginForm *LoginForm) (*LoginMessage, error)
	RegisterUser(ctx context.Context, registerForm *RegisterForm) (*RegisterMessage, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

// UserUsecase 用户用例
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase 创建用户用例
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// ListUsersRequest 用户列表请求
type ListUsersRequest struct {
	Page     int32
	PageSize int32
	Username string
	Email    string
	Status   int32
}

// ListUsersReply 用户列表响应
type ListUsersReply struct {
	Total int32   `json:"total"`
	Users []*User `json:"users"`
}

// LoginForm 登录表单
type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginMessage 登录结果
type LoginMessage struct {
	User    *User  `json:"user"`
	Token   string `json:"token"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// RegisterForm 注册表单
type RegisterForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Age      int    `json:"age"`
	Avatar   string `json:"avatar"`
}

// GetMeMessage 获取当前用户信息结果
type GetMeMessage struct {
	User    *User  `json:"user"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// RegisterMessage 注册结果
type RegisterMessage struct {
	User    *User
	Message string
	Success bool
}

// GetUser 获取用户
func (uc *UserUsecase) GetUser(ctx context.Context, id int32) (*User, error) {
	user, err := uc.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FormatTimeFields()
	return user, nil
}

// ListUsers 获取用户列表
func (uc *UserUsecase) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersReply, error) {
	reply, err := uc.repo.ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	FormatUserTimeFieldsBatch(reply.Users)
	return reply, nil
}

// CreateUser 创建用户
func (uc *UserUsecase) CreateUser(ctx context.Context, user *User) (*User, error) {
	// 加密密码
	hashedPassword, err := password.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	createdUser, err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	createdUser.FormatTimeFields()
	return createdUser, nil
}

// UpdateUser 更新用户
func (uc *UserUsecase) UpdateUser(ctx context.Context, user *User) (*User, error) {
	// 如果有新密码，则加密
	if user.Password != "" {
		hashedPassword, err := password.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	updatedUser, err := uc.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	updatedUser.FormatTimeFields()
	return updatedUser, nil
}

// DeleteUser 删除用户
func (uc *UserUsecase) DeleteUser(ctx context.Context, id int32) error {
	return uc.repo.DeleteUser(ctx, id)
}

// Login 用户登录
func (uc *UserUsecase) Login(ctx context.Context, loginForm *LoginForm) (*LoginMessage, error) {
	uc.log.Info("user login", loginForm.Username)

	// 调用数据层进行登录
	result, err := uc.repo.Login(ctx, loginForm)
	if err != nil {
		return &LoginMessage{
			Message: "登录失败，请稍后重试",
			Success: false,
		}, err
	}

	return result, nil
}

// GetMe 获取当前用户信息
func (uc *UserUsecase) GetMe(ctx context.Context, userID uint) (*GetMeMessage, error) {
	uc.log.Info("get current user", userID)

	// 通过用户ID获取用户信息
	user, err := uc.repo.GetUser(ctx, int32(userID))
	if err != nil {
		return &GetMeMessage{
			Message: "用户不存在",
			Success: false,
		}, nil
	}

	return &GetMeMessage{
		User:    user,
		Message: "获取用户信息成功",
		Success: true,
	}, nil
}

// Register 用户注册
func (uc *UserUsecase) Register(ctx context.Context, registerForm *RegisterForm) (*RegisterMessage, error) {
	uc.log.Info("user register", registerForm.Username)

	// 检查用户名是否已存在
	existingUser, err := uc.repo.GetUserByUsername(ctx, registerForm.Username)
	if err == nil && existingUser != nil {
		return &RegisterMessage{
			Message: "用户名已存在",
			Success: false,
		}, nil
	}

	// 检查邮箱是否已存在
	if registerForm.Email != "" {
		existingUser, err = uc.repo.GetUserByEmail(ctx, registerForm.Email)
		if err == nil && existingUser != nil {
			return &RegisterMessage{
				Message: "邮箱已存在",
				Success: false,
			}, nil
		}
	}

	// 调用数据层进行注册
	result, err := uc.repo.RegisterUser(ctx, registerForm)
	if err != nil {
		return &RegisterMessage{
			Message: "注册失败，请稍后重试",
			Success: false,
		}, err
	}

	return result, nil
}
