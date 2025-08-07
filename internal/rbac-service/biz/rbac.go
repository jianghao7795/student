package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

const (
	// TimeFormat 时间格式常量
	TimeFormat = "2006-01-02 15:04:05"
)

// Role 角色模型
type Role struct {
	ID          uint
	Name        string
	Description string
	Status      int
	CreatedAt   *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`

	// 格式化后的时间字符串
	CreatedAtStr string `gorm:"-" json:"created_at_str,omitempty"`
	UpdatedAtStr string `gorm:"-" json:"updated_at_str,omitempty"`
}

// FormatTimeFields 格式化时间字段
func (r *Role) FormatTimeFields() {
	if r.CreatedAt != nil {
		r.CreatedAtStr = r.CreatedAt.Format(TimeFormat)
	}
	if r.UpdatedAt != nil {
		r.UpdatedAtStr = r.UpdatedAt.Format(TimeFormat)
	}
}

// Permission 权限模型
type Permission struct {
	ID          uint
	Name        string
	Resource    string
	Action      string
	Description string
	Status      int
	CreatedAt   *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`

	// 格式化后的时间字符串
	CreatedAtStr string `gorm:"-" json:"created_at_str,omitempty"`
	UpdatedAtStr string `gorm:"-" json:"updated_at_str,omitempty"`
}

// FormatTimeFields 格式化时间字段
func (p *Permission) FormatTimeFields() {
	if p.CreatedAt != nil {
		p.CreatedAtStr = p.CreatedAt.Format(TimeFormat)
	}
	if p.UpdatedAt != nil {
		p.UpdatedAtStr = p.UpdatedAt.Format(TimeFormat)
	}
}

// RBACRepo RBAC仓储接口
type RBACRepo interface {
	// 角色管理
	GetRoles(ctx context.Context) ([]*Role, error)
	CreateRole(ctx context.Context, role *Role) (*Role, error)
	GetRole(ctx context.Context, id int32) (*Role, error)
	UpdateRole(ctx context.Context, role *Role) (*Role, error)
	DeleteRole(ctx context.Context, id int32) error

	// 权限管理
	GetPermissions(ctx context.Context) ([]*Permission, error)
	CreatePermission(ctx context.Context, permission *Permission) (*Permission, error)
	GetPermission(ctx context.Context, id int32) (*Permission, error)
	UpdatePermission(ctx context.Context, permission *Permission) (*Permission, error)
	DeletePermission(ctx context.Context, id int32) error

	// 用户角色关联
	GetUserRoles(ctx context.Context, userID int32) ([]*Role, error)
	GetUserRoleNames(ctx context.Context, userID int32) ([]string, error)
	AssignRoleToUser(ctx context.Context, userID int32, roleID int32) error
	RemoveRoleFromUser(ctx context.Context, userID int32, roleID int32) error

	// 权限检查
	CheckPermission(ctx context.Context, userID int32, resource string, action string) (bool, error)
}

// RBACUsecase RBAC用例
type RBACUsecase struct {
	repo RBACRepo
	log  *log.Helper
}

// NewRBACUsecase 创建RBAC用例
func NewRBACUsecase(repo RBACRepo, logger log.Logger) *RBACUsecase {
	return &RBACUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// GetRoles 获取角色列表
func (uc *RBACUsecase) GetRoles(ctx context.Context) ([]*Role, error) {
	roles, err := uc.repo.GetRoles(ctx)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		role.FormatTimeFields()
	}

	return roles, nil
}

// CreateRole 创建角色
func (uc *RBACUsecase) CreateRole(ctx context.Context, role *Role) (*Role, error) {
	createdRole, err := uc.repo.CreateRole(ctx, role)
	if err != nil {
		return nil, err
	}

	createdRole.FormatTimeFields()
	return createdRole, nil
}

// GetPermissions 获取权限列表
func (uc *RBACUsecase) GetPermissions(ctx context.Context) ([]*Permission, error) {
	permissions, err := uc.repo.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	for _, permission := range permissions {
		permission.FormatTimeFields()
	}

	return permissions, nil
}

// CreatePermission 创建权限
func (uc *RBACUsecase) CreatePermission(ctx context.Context, permission *Permission) (*Permission, error) {
	createdPermission, err := uc.repo.CreatePermission(ctx, permission)
	if err != nil {
		return nil, err
	}

	createdPermission.FormatTimeFields()
	return createdPermission, nil
}

// GetUserRoleNames 获取用户角色名称列表
func (uc *RBACUsecase) GetUserRoleNames(ctx context.Context, userID int32) ([]string, error) {
	return uc.repo.GetUserRoleNames(ctx, userID)
}

// CheckPermission 检查权限
func (uc *RBACUsecase) CheckPermission(ctx context.Context, userID int32, resource string, action string) (bool, error) {
	return uc.repo.CheckPermission(ctx, userID, resource, action)
}

// AssignRoleToUser 为用户分配角色
func (uc *RBACUsecase) AssignRoleToUser(ctx context.Context, userID int32, roleID int32) error {
	return uc.repo.AssignRoleToUser(ctx, userID, roleID)
}

// RemoveRoleFromUser 移除用户角色
func (uc *RBACUsecase) RemoveRoleFromUser(ctx context.Context, userID int32, roleID int32) error {
	return uc.repo.RemoveRoleFromUser(ctx, userID, roleID)
}
