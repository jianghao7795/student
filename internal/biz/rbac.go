package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
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

// UserRole 用户角色关联模型
type UserRole struct {
	ID        uint
	UserID    uint
	RoleID    uint
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`

	// 关联的角色信息
	Role *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// RolePermission 角色权限关联模型
type RolePermission struct {
	ID           uint
	RoleID       uint
	PermissionID uint
	CreatedAt    *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updated_at"`

	// 关联的权限信息
	Permission *Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

// RoleForm 角色表单
type RoleForm struct {
	Name        string
	Description string
	Status      int
}

// PermissionForm 权限表单
type PermissionForm struct {
	Name        string
	Resource    string
	Action      string
	Description string
	Status      int
}

// UserRoleForm 用户角色表单
type UserRoleForm struct {
	UserID uint
	RoleID uint
}

// RolePermissionForm 角色权限表单
type RolePermissionForm struct {
	RoleID       uint
	PermissionID uint
}

// 定义 RBAC 的操作接口
type RBACRepo interface {
	// 角色相关
	GetRole(ctx context.Context, id int32) (*Role, error)
	CreateRole(ctx context.Context, r *RoleForm) (*Role, error)
	UpdateRole(ctx context.Context, id int32, r *RoleForm) (*Role, error)
	DeleteRole(ctx context.Context, id int32) error
	ListRoles(ctx context.Context, page int32, pageSize int32, name string) ([]*Role, int32, error)
	GetRoleByName(ctx context.Context, name string) (*Role, error)

	// 权限相关
	GetPermission(ctx context.Context, id int32) (*Permission, error)
	CreatePermission(ctx context.Context, p *PermissionForm) (*Permission, error)
	UpdatePermission(ctx context.Context, id int32, p *PermissionForm) (*Permission, error)
	DeletePermission(ctx context.Context, id int32) error
	ListPermissions(ctx context.Context, page int32, pageSize int32, name, resource string) ([]*Permission, int32, error)
	GetPermissionByResourceAction(ctx context.Context, resource, action string) (*Permission, error)

	// 用户角色相关
	GetUserRoles(ctx context.Context, userID int32) ([]*UserRole, error)
	AssignUserRole(ctx context.Context, userID, roleID int32) error
	RemoveUserRole(ctx context.Context, userID, roleID int32) error
	GetUserRoleNames(ctx context.Context, userID int32) ([]string, error)

	// 角色权限相关
	GetRolePermissions(ctx context.Context, roleID int32) ([]*RolePermission, error)
	AssignRolePermission(ctx context.Context, roleID, permissionID int32) error
	RemoveRolePermission(ctx context.Context, roleID, permissionID int32) error
	GetRolePermissionNames(ctx context.Context, roleID int32) ([]string, error)

	// Casbin相关
	LoadPolicy(ctx context.Context) error
	SavePolicy(ctx context.Context) error
	AddPolicy(ctx context.Context, sub, obj, act string) error
	RemovePolicy(ctx context.Context, sub, obj, act string) error
	AddRoleForUser(ctx context.Context, user, role string) error
	RemoveRoleForUser(ctx context.Context, user, role string) error
	GetRolesForUser(ctx context.Context, user string) ([]string, error)
	GetPermissionsForUser(ctx context.Context, user string) ([][]string, error)
}

type RBACUsecase struct {
	repo RBACRepo
	log  *log.Helper
}

// 初始化 RBACUsecase
func NewRBACUsecase(repo RBACRepo, logger log.Logger) *RBACUsecase {
	return &RBACUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// 角色相关方法
func (uc *RBACUsecase) GetRole(ctx context.Context, id int32) (*Role, error) {
	uc.log.Info("get role by id", id)
	return uc.repo.GetRole(ctx, id)
}

func (uc *RBACUsecase) CreateRole(ctx context.Context, r *RoleForm) (*Role, error) {
	uc.log.Info("create role", r)
	return uc.repo.CreateRole(ctx, r)
}

func (uc *RBACUsecase) UpdateRole(ctx context.Context, id int32, r *RoleForm) (*Role, error) {
	uc.log.Info("update role", id, r)
	return uc.repo.UpdateRole(ctx, id, r)
}

func (uc *RBACUsecase) DeleteRole(ctx context.Context, id int32) error {
	uc.log.Info("delete role", id)
	return uc.repo.DeleteRole(ctx, id)
}

func (uc *RBACUsecase) ListRoles(ctx context.Context, page int32, pageSize int32, name string) ([]*Role, int32, error) {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	return uc.repo.ListRoles(ctx, page, pageSize, name)
}

// 权限相关方法
func (uc *RBACUsecase) GetPermission(ctx context.Context, id int32) (*Permission, error) {
	uc.log.Info("get permission by id", id)
	return uc.repo.GetPermission(ctx, id)
}

func (uc *RBACUsecase) CreatePermission(ctx context.Context, p *PermissionForm) (*Permission, error) {
	uc.log.Info("create permission", p)
	return uc.repo.CreatePermission(ctx, p)
}

func (uc *RBACUsecase) UpdatePermission(ctx context.Context, id int32, p *PermissionForm) (*Permission, error) {
	uc.log.Info("update permission", id, p)
	return uc.repo.UpdatePermission(ctx, id, p)
}

func (uc *RBACUsecase) DeletePermission(ctx context.Context, id int32) error {
	uc.log.Info("delete permission", id)
	return uc.repo.DeletePermission(ctx, id)
}

func (uc *RBACUsecase) ListPermissions(ctx context.Context, page int32, pageSize int32, name, resource string) ([]*Permission, int32, error) {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	return uc.repo.ListPermissions(ctx, page, pageSize, name, resource)
}

// 用户角色相关方法
func (uc *RBACUsecase) GetUserRoles(ctx context.Context, userID int32) ([]*UserRole, error) {
	uc.log.Info("get user roles", userID)
	return uc.repo.GetUserRoles(ctx, userID)
}

func (uc *RBACUsecase) AssignUserRole(ctx context.Context, userID, roleID int32) error {
	uc.log.Info("assign user role", userID, roleID)
	return uc.repo.AssignUserRole(ctx, userID, roleID)
}

func (uc *RBACUsecase) RemoveUserRole(ctx context.Context, userID, roleID int32) error {
	uc.log.Info("remove user role", userID, roleID)
	return uc.repo.RemoveUserRole(ctx, userID, roleID)
}

func (uc *RBACUsecase) GetUserRoleNames(ctx context.Context, userID int32) ([]string, error) {
	uc.log.Info("get user role names", userID)
	return uc.repo.GetUserRoleNames(ctx, userID)
}

// 角色权限相关方法
func (uc *RBACUsecase) GetRolePermissions(ctx context.Context, roleID int32) ([]*RolePermission, error) {
	uc.log.Info("get role permissions", roleID)
	return uc.repo.GetRolePermissions(ctx, roleID)
}

func (uc *RBACUsecase) AssignRolePermission(ctx context.Context, roleID, permissionID int32) error {
	uc.log.Info("assign role permission", roleID, permissionID)
	return uc.repo.AssignRolePermission(ctx, roleID, permissionID)
}

func (uc *RBACUsecase) RemoveRolePermission(ctx context.Context, roleID, permissionID int32) error {
	uc.log.Info("remove role permission", roleID, permissionID)
	return uc.repo.RemoveRolePermission(ctx, roleID, permissionID)
}

func (uc *RBACUsecase) GetRolePermissionNames(ctx context.Context, roleID int32) ([]string, error) {
	uc.log.Info("get role permission names", roleID)
	return uc.repo.GetRolePermissionNames(ctx, roleID)
}

// Casbin相关方法
func (uc *RBACUsecase) LoadPolicy(ctx context.Context) error {
	uc.log.Info("load casbin policy")
	return uc.repo.LoadPolicy(ctx)
}

func (uc *RBACUsecase) SavePolicy(ctx context.Context) error {
	uc.log.Info("save casbin policy")
	return uc.repo.SavePolicy(ctx)
}

func (uc *RBACUsecase) AddPolicy(ctx context.Context, sub, obj, act string) error {
	uc.log.Info("add casbin policy", sub, obj, act)
	return uc.repo.AddPolicy(ctx, sub, obj, act)
}

func (uc *RBACUsecase) RemovePolicy(ctx context.Context, sub, obj, act string) error {
	uc.log.Info("remove casbin policy", sub, obj, act)
	return uc.repo.RemovePolicy(ctx, sub, obj, act)
}

func (uc *RBACUsecase) AddRoleForUser(ctx context.Context, user, role string) error {
	uc.log.Info("add role for user", user, role)
	return uc.repo.AddRoleForUser(ctx, user, role)
}

func (uc *RBACUsecase) RemoveRoleForUser(ctx context.Context, user, role string) error {
	uc.log.Info("remove role for user", user, role)
	return uc.repo.RemoveRoleForUser(ctx, user, role)
}

func (uc *RBACUsecase) GetRolesForUser(ctx context.Context, user string) ([]string, error) {
	uc.log.Info("get roles for user", user)
	return uc.repo.GetRolesForUser(ctx, user)
}

func (uc *RBACUsecase) GetPermissionsForUser(ctx context.Context, user string) ([][]string, error) {
	uc.log.Info("get permissions for user", user)
	return uc.repo.GetPermissionsForUser(ctx, user)
}

// 权限检查方法
func (uc *RBACUsecase) CheckPermission(ctx context.Context, user string, obj string, act string) (bool, error) {
	uc.log.Info("check permission", user, obj, act)
	permissions, err := uc.repo.GetPermissionsForUser(ctx, user)
	if err != nil {
		return false, err
	}

	for _, permission := range permissions {
		if len(permission) >= 3 {
			// 检查资源匹配（支持通配符）
			if uc.matchResource(permission[1], obj) && (permission[2] == act || permission[2] == "*") {
				return true, nil
			}
		}
	}
	return false, nil
}

// 简单的资源匹配方法（支持通配符）
func (uc *RBACUsecase) matchResource(pattern, resource string) bool {
	// 这里可以实现更复杂的匹配逻辑
	// 目前简单实现：如果pattern以*结尾，则匹配前缀
	if len(pattern) > 1 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(resource) >= len(prefix) && resource[:len(prefix)] == prefix
	}
	return pattern == resource
}
