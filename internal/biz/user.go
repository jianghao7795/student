package biz

import (
	"context"
	"time"

	"student/internal/pkg/jwt"
	"student/internal/pkg/password"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

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

// UserForm 用户表单
type UserForm struct {
	Username string
	Email    string
	Phone    string
	Password string
	Status   int
	Age      int
	Avatar   string
}

// HashPassword 加密密码
func (u *UserForm) HashPassword() error {
	if u.Password != "" && !password.IsHashed(u.Password) {
		hashedPassword, err := password.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	return nil
}

// CheckPassword 验证密码
func (u *UserForm) CheckPassword(hashedPassword string) bool {
	return password.CheckPassword(u.Password, hashedPassword)
}

// CreateUserMessage 创建用户消息
type CreateUserMessage struct {
	ID      int32
	Message string
}

// UpdateUserMessage 更新用户消息
type UpdateUserMessage struct {
	Message string
}

// DeleteUserMessage 删除用户消息
type DeleteUserMessage struct {
	Message string
}

// 定义 User 的操作接口
type UserRepo interface {
	GetUser(ctx context.Context, id int32) (*User, error)
	CreateUser(ctx context.Context, u *UserForm) (*CreateUserMessage, error)
	UpdateUser(ctx context.Context, id int32, u *UserForm) (*UpdateUserMessage, error)
	DeleteUser(ctx context.Context, id int32) (*DeleteUserMessage, error)
	ListUsers(ctx context.Context, page int32, pageSize int32, username, email string) ([]*User, int32, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	RegisterUser(ctx context.Context, u *RegisterForm) (*RegisterMessage, error)
}

// LoginForm 登录表单
type LoginForm struct {
	Username string
	Password string
}

// LoginMessage 登录消息
type LoginMessage struct {
	User    *User
	Message string
	Success bool
	Token   string
}

type UserUsecase struct {
	repo    UserRepo
	rbacUC  *RBACUsecase
	log     *log.Helper
	jwtUtil *jwt.JWTUtil
}

// 初始化 UserUsecase
func NewUserUsecase(repo UserRepo, rbacUC *RBACUsecase, jwtUtil *jwt.JWTUtil, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo:    repo,
		rbacUC:  rbacUC,
		log:     log.NewHelper(logger),
		jwtUtil: jwtUtil,
	}
}

// 通过 id 获取用户信息
func (uc *UserUsecase) Get(ctx context.Context, id int32) (*User, error) {
	uc.log.Info("get user by id", id)
	return uc.repo.GetUser(ctx, id)
}

// 创建用户
func (uc *UserUsecase) Create(ctx context.Context, u *UserForm) (*CreateUserMessage, error) {
	uc.log.Info("create user", u)
	return uc.repo.CreateUser(ctx, u)
}

// 更新用户
func (uc *UserUsecase) Update(ctx context.Context, id int32, u *UserForm) (*UpdateUserMessage, error) {
	uc.log.Info("update user", id, u)
	return uc.repo.UpdateUser(ctx, id, u)
}

// 删除用户
func (uc *UserUsecase) Delete(ctx context.Context, id int32) (*DeleteUserMessage, error) {
	uc.log.Info("delete user", id)
	return uc.repo.DeleteUser(ctx, id)
}

// 获取用户列表
func (uc *UserUsecase) List(ctx context.Context, page int32, pageSize int32, username, email string) ([]*User, int32, error) {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	return uc.repo.ListUsers(ctx, page, pageSize, username, email)
}

// 通过用户名获取用户信息
func (uc *UserUsecase) GetByUsername(ctx context.Context, username string) (*User, error) {
	uc.log.Info("get user by username", username)
	return uc.repo.GetUserByUsername(ctx, username)
}

// 通过邮箱获取用户信息
func (uc *UserUsecase) GetByEmail(ctx context.Context, email string) (*User, error) {
	uc.log.Info("get user by email", email)
	return uc.repo.GetUserByEmail(ctx, email)
}

// 用户登录验证
func (uc *UserUsecase) Login(ctx context.Context, loginForm *LoginForm) (*LoginMessage, error) {
	uc.log.Info("user login", loginForm.Username)

	// 通过用户名获取用户
	user, err := uc.repo.GetUserByUsername(ctx, loginForm.Username)
	if err != nil {
		return &LoginMessage{
			Message: "用户名或密码错误",
			Success: false,
		}, nil
	}

	// 验证密码
	if !password.CheckPassword(loginForm.Password, user.Password) {
		return &LoginMessage{
			Message: "用户名或密码错误",
			Success: false,
		}, nil
	}

	// 检查用户状态
	if user.Status != 1 {
		return &LoginMessage{
			Message: "用户已被禁用",
			Success: false,
		}, nil
	}

	// 获取用户角色
	roles, err := uc.rbacUC.GetUserRoleNames(ctx, int32(user.ID))
	if err != nil {
		uc.log.Error("获取用户角色失败", err)
		return &LoginMessage{
			Message: "登录失败，请稍后重试",
			Success: false,
		}, nil
	}

	// 设置用户角色信息
	user.Roles = roles

	// 生成JWT token
	token, err := uc.jwtUtil.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		uc.log.Error("生成JWT token失败", err)
		return &LoginMessage{
			Message: "登录失败，请稍后重试",
			Success: false,
		}, nil
	}

	uc.log.Info("JWT token生成成功", "token_length", len(token))

	return &LoginMessage{
		User:    user,
		Message: "登录成功",
		Success: true,
		Token:   token,
	}, nil
}

// GetMeMessage 获取当前用户信息消息
type GetMeMessage struct {
	User    *User
	Message string
	Success bool
}

// RegisterForm 注册表单
type RegisterForm struct {
	Username string
	Email    string
	Phone    string
	Password string
	Age      int
	Avatar   string
}

// HashPassword 加密密码
func (r *RegisterForm) HashPassword() error {
	if r.Password != "" && !password.IsHashed(r.Password) {
		hashedPassword, err := password.HashPassword(r.Password)
		if err != nil {
			return err
		}
		r.Password = hashedPassword
	}
	return nil
}

// RegisterMessage 注册消息
type RegisterMessage struct {
	User    *User
	Message string
	Success bool
}

// 获取当前用户信息
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

	// 获取用户角色
	roles, err := uc.rbacUC.GetUserRoleNames(ctx, int32(user.ID))
	if err != nil {
		uc.log.Error("获取用户角色失败", err)
		return &GetMeMessage{
			Message: "获取用户信息失败",
			Success: false,
		}, nil
	}

	// 设置用户角色信息
	user.Roles = roles

	return &GetMeMessage{
		User:    user,
		Message: "获取用户信息成功",
		Success: true,
	}, nil
}

// 用户注册
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
