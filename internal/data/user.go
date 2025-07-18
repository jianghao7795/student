package data

import (
	"context"

	"student/internal/biz"

	errors "student/internal/data/errors"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// 实现 从 gormDB 中获取用户信息
func (r *userRepo) GetUser(ctx context.Context, id int32) (*biz.User, error) {
	var user biz.User
	err := r.data.gormDB.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Error404()
		}
		return nil, errors.Error400(err)
	}
	r.log.WithContext(ctx).Info("gormDB: GetUser, id: %d, result: %v", id, user)

	// 格式化时间字段
	user.FormatTimeFields()

	return &biz.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Phone:        user.Phone,
		Password:     user.Password,
		Status:       user.Status,
		Age:          user.Age,
		Avatar:       user.Avatar,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		CreatedAtStr: user.CreatedAtStr,
		UpdatedAtStr: user.UpdatedAtStr,
	}, err
}

// 实现 从 gormDB 中创建用户
func (r *userRepo) CreateUser(ctx context.Context, u *biz.UserForm) (*biz.CreateUserMessage, error) {
	// 加密密码
	if err := u.HashPassword(); err != nil {
		return nil, errors.Error400(err)
	}

	var user biz.User
	user.Username = u.Username
	user.Email = u.Email
	user.Phone = u.Phone
	user.Password = u.Password
	user.Status = u.Status
	user.Age = u.Age
	user.Avatar = u.Avatar

	err := r.data.gormDB.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, errors.Error400(err)
	}
	r.log.WithContext(ctx).Info("gormDB: CreateUser, user: %v", user)
	return &biz.CreateUserMessage{
		Message: "Create user success",
	}, err
}

// 实现 从 gormDB 中更新用户信息
func (r *userRepo) UpdateUser(ctx context.Context, id int32, u *biz.UserForm) (*biz.UpdateUserMessage, error) {
	var user biz.User
	err := r.data.gormDB.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, errors.Error404()
	}

	user.Username = u.Username
	user.Email = u.Email
	user.Phone = u.Phone
	if u.Password != "" {
		// 如果提供了新密码，则加密
		if err := u.HashPassword(); err != nil {
			return nil, errors.Error400(err)
		}
		user.Password = u.Password
	}
	user.Status = u.Status
	user.Age = u.Age
	user.Avatar = u.Avatar

	err = r.data.gormDB.WithContext(ctx).Save(&user).Error
	if err != nil {
		return nil, errors.Error400(err)
	}
	r.log.WithContext(ctx).Info("gormDB: UpdateUser, id: %d, user: %v", id, user)
	return &biz.UpdateUserMessage{
		Message: "Update user success",
	}, err
}

// 实现 从 gormDB 中删除用户
func (r *userRepo) DeleteUser(ctx context.Context, id int32) (*biz.DeleteUserMessage, error) {
	var user biz.User
	err := r.data.gormDB.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Error404()
		}
		return nil, errors.Error400(err)
	}
	err = r.data.gormDB.Delete(&user, id).Error
	r.log.WithContext(ctx).Info("gormDB: DeleteUser, id: %d", id)
	return &biz.DeleteUserMessage{
		Message: "Delete user success",
	}, err
}

// 实现 从 gormDB 中获取用户列表
func (r *userRepo) ListUsers(ctx context.Context, page int32, pageSize int32, username, email string) ([]*biz.User, int32, error) {
	var users []*biz.User
	var total int64
	var err error

	// 构建查询条件
	query := r.data.gormDB.WithContext(ctx).Model(&biz.User{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	// 获取总数
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, errors.Error400(err)
	}

	// 获取分页数据
	query = r.data.gormDB.WithContext(ctx).Model(&biz.User{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	err = query.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Order("id desc").Find(&users).Error
	if err != nil {
		return nil, 0, errors.Error400(err)
	}

	// 为每个用户记录格式化时间字段
	biz.FormatUserTimeFieldsBatch(users)

	return users, int32(total), err
}

// 实现 通过用户名获取用户信息
func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*biz.User, error) {
	var user biz.User
	err := r.data.gormDB.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Error404()
		}
		return nil, errors.Error400(err)
	}
	r.log.WithContext(ctx).Info("gormDB: GetUserByUsername, username: %s, result: %v", username, user)

	// 格式化时间字段
	user.FormatTimeFields()

	return &biz.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Phone:        user.Phone,
		Password:     user.Password,
		Status:       user.Status,
		Age:          user.Age,
		Avatar:       user.Avatar,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		CreatedAtStr: user.CreatedAtStr,
		UpdatedAtStr: user.UpdatedAtStr,
	}, err
}

// 实现 通过邮箱获取用户信息
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	var user biz.User
	err := r.data.gormDB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Error404()
		}
		return nil, errors.Error400(err)
	}
	r.log.WithContext(ctx).Info("gormDB: GetUserByEmail, email: %s, result: %v", email, user)

	// 格式化时间字段
	user.FormatTimeFields()

	return &biz.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Phone:        user.Phone,
		Password:     user.Password,
		Status:       user.Status,
		Age:          user.Age,
		Avatar:       user.Avatar,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		CreatedAtStr: user.CreatedAtStr,
		UpdatedAtStr: user.UpdatedAtStr,
	}, err
}
