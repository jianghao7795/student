package data

import (
	"context"
	"errors"
	"time"

	"student/internal/pkg/jwt"
	"student/internal/pkg/password"
	"student/internal/user-service/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type userRepo struct {
	data    *Data
	log     *log.Helper
	jwtUtil *jwt.JWTUtil
}

func NewUserRepo(data *Data, logger log.Logger, jwtUtil *jwt.JWTUtil) biz.UserRepo {
	return &userRepo{
		data:    data,
		log:     log.NewHelper(logger),
		jwtUtil: jwtUtil,
	}
}

// GetUser 从 gormDB 中获取用户信息
func (r *userRepo) GetUser(ctx context.Context, id int32) (*biz.User, error) {
	var user biz.User
	err := r.data.gormDB.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: GetUser, id: %d, result: %v", id, user)

	// 格式化时间字段
	user.FormatTimeFields()

	return &user, nil
}

// ListUsers 获取用户列表
func (r *userRepo) ListUsers(ctx context.Context, req *biz.ListUsersRequest) (*biz.ListUsersReply, error) {
	var users []*biz.User
	var total int64

	query := r.data.gormDB.WithContext(ctx).Model(&biz.User{})

	// 添加筛选条件
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Email != "" {
		query = query.Where("email LIKE ?", "%"+req.Email+"%")
	}
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(int(offset)).Limit(int(req.PageSize)).Find(&users).Error; err != nil {
		return nil, err
	}

	return &biz.ListUsersReply{
		Total: int32(total),
		Users: users,
	}, nil
}

// CreateUser 创建用户
func (r *userRepo) CreateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	now := time.Now()
	user.CreatedAt = &now
	user.UpdatedAt = &now
	user.Status = 1 // 默认状态为启用

	if err := r.data.gormDB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	r.log.WithContext(ctx).Infof("gormDB: CreateUser, user: %v", user)
	return user, nil
}

// UpdateUser 更新用户
func (r *userRepo) UpdateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	now := time.Now()
	user.UpdatedAt = &now

	if err := r.data.gormDB.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}

	r.log.WithContext(ctx).Infof("gormDB: UpdateUser, user: %v", user)
	return user, nil
}

// DeleteUser 删除用户
func (r *userRepo) DeleteUser(ctx context.Context, id int32) error {
	if err := r.data.gormDB.WithContext(ctx).Delete(&biz.User{}, id).Error; err != nil {
		return err
	}

	r.log.WithContext(ctx).Infof("gormDB: DeleteUser, id: %d", id)
	return nil
}

// GetUserByUsername 根据用户名获取用户
func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*biz.User, error) {
	var user biz.User
	err := r.data.gormDB.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.WithContext(ctx).Infof("gormDB: GetUserByUsername, username: %s, result: record not found", username)
			return nil, err
		}
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: GetUserByUsername, username: %s, result: %v", username, user)
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	var user biz.User
	err := r.data.gormDB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.WithContext(ctx).Infof("gormDB: GetUserByEmail, email: %s, result: record not found", email)
			return nil, err
		}
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: GetUserByEmail, email: %s, result: %v", email, user)
	return &user, nil
}

// Login 用户登录
func (r *userRepo) Login(ctx context.Context, loginForm *biz.LoginForm) (*biz.LoginMessage, error) {
	r.log.WithContext(ctx).Infof("gormDB: Login, username: %s", loginForm.Username)

	// 根据用户名查找用户
	user, err := r.GetUserByUsername(ctx, loginForm.Username)
	if err != nil {
		return &biz.LoginMessage{
			Message: "用户名或密码错误",
			Success: false,
		}, nil
	}

	// 验证密码
	if !password.CheckPassword(loginForm.Password, user.Password) {
		return &biz.LoginMessage{
			Message: "用户名或密码错误",
			Success: false,
		}, nil
	}

	// 检查用户状态
	if user.Status != 1 {
		return &biz.LoginMessage{
			Message: "用户已被禁用",
			Success: false,
		}, nil
	}

	// 生成JWT令牌
	token, err := r.jwtUtil.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return &biz.LoginMessage{
			Message: "生成令牌失败",
			Success: false,
		}, err
	}

	// 格式化时间字段
	user.FormatTimeFields()

	return &biz.LoginMessage{
		User:    user,
		Token:   token,
		Message: "登录成功",
		Success: true,
	}, nil
}

// RegisterUser 用户注册
func (r *userRepo) RegisterUser(ctx context.Context, registerForm *biz.RegisterForm) (*biz.RegisterMessage, error) {
	r.log.WithContext(ctx).Infof("gormDB: RegisterUser, registerForm: %+v", registerForm)

	// 加密密码
	hashedPassword, err := password.HashPassword(registerForm.Password)
	if err != nil {
		return &biz.RegisterMessage{
			Message: "密码加密失败",
			Success: false,
		}, err
	}

	// 创建用户对象
	now := time.Now()
	user := &biz.User{
		Username:  registerForm.Username,
		Email:     registerForm.Email,
		Phone:     registerForm.Phone,
		Password:  hashedPassword,
		Status:    1, // 默认启用
		Age:       registerForm.Age,
		Avatar:    registerForm.Avatar,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// 保存到数据库
	if err := r.data.gormDB.WithContext(ctx).Create(user).Error; err != nil {
		return &biz.RegisterMessage{
			Message: "注册失败",
			Success: false,
		}, err
	}

	r.log.WithContext(ctx).Infof("gormDB: RegisterUser, user: %v", user)

	// 格式化时间字段
	user.FormatTimeFields()

	return &biz.RegisterMessage{
		User:    user,
		Message: "注册成功",
		Success: true,
	}, nil
}
