package data

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"student/internal/rbac-service/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type rbacRepo struct {
	data *Data
	log  *log.Helper
}

// NewRBACRepo .
func NewRBACRepo(data *Data, logger log.Logger) biz.RBACRepo {
	return &rbacRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// GetRoles 获取角色列表
func (r *rbacRepo) GetRoles(ctx context.Context) ([]*biz.Role, error) {
	var roles []*biz.Role
	err := r.data.gormDB.WithContext(ctx).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: GetRoles, count: %d", len(roles))
	return roles, nil
}

// CreateRole 创建角色
func (r *rbacRepo) CreateRole(ctx context.Context, role *biz.Role) (*biz.Role, error) {
	now := time.Now()
	role.CreatedAt = &now
	role.UpdatedAt = &now
	role.Status = 1 // 默认启用

	err := r.data.gormDB.WithContext(ctx).Create(role).Error
	if err != nil {
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: CreateRole, role: %v", role)
	return role, nil
}

// GetRole 获取角色
func (r *rbacRepo) GetRole(ctx context.Context, id int32) (*biz.Role, error) {
	var role biz.Role
	err := r.data.gormDB.WithContext(ctx).First(&role, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: GetRole, id: %d, result: %v", id, role)
	return &role, nil
}

// UpdateRole 更新角色
func (r *rbacRepo) UpdateRole(ctx context.Context, role *biz.Role) (*biz.Role, error) {
	now := time.Now()
	role.UpdatedAt = &now

	err := r.data.gormDB.WithContext(ctx).Save(role).Error
	if err != nil {
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: UpdateRole, role: %v", role)
	return role, nil
}

// DeleteRole 删除角色
func (r *rbacRepo) DeleteRole(ctx context.Context, id int32) error {
	err := r.data.gormDB.WithContext(ctx).Delete(&biz.Role{}, id).Error
	if err != nil {
		return err
	}
	r.log.WithContext(ctx).Infof("gormDB: DeleteRole, id: %d", id)
	return nil
}

// GetPermissions 获取权限列表
func (r *rbacRepo) GetPermissions(ctx context.Context) ([]*biz.Permission, error) {
	var permissions []*biz.Permission
	err := r.data.gormDB.WithContext(ctx).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: GetPermissions, count: %d", len(permissions))
	return permissions, nil
}

// CreatePermission 创建权限
func (r *rbacRepo) CreatePermission(ctx context.Context, permission *biz.Permission) (*biz.Permission, error) {
	now := time.Now()
	permission.CreatedAt = &now
	permission.UpdatedAt = &now
	permission.Status = 1 // 默认启用

	err := r.data.gormDB.WithContext(ctx).Create(permission).Error
	if err != nil {
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: CreatePermission, permission: %v", permission)
	return permission, nil
}

// GetPermission 获取权限
func (r *rbacRepo) GetPermission(ctx context.Context, id int32) (*biz.Permission, error) {
	var permission biz.Permission
	err := r.data.gormDB.WithContext(ctx).First(&permission, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: GetPermission, id: %d, result: %v", id, permission)
	return &permission, nil
}

// UpdatePermission 更新权限
func (r *rbacRepo) UpdatePermission(ctx context.Context, permission *biz.Permission) (*biz.Permission, error) {
	now := time.Now()
	permission.UpdatedAt = &now

	err := r.data.gormDB.WithContext(ctx).Save(permission).Error
	if err != nil {
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: UpdatePermission, permission: %v", permission)
	return permission, nil
}

// DeletePermission 删除权限
func (r *rbacRepo) DeletePermission(ctx context.Context, id int32) error {
	err := r.data.gormDB.WithContext(ctx).Delete(&biz.Permission{}, id).Error
	if err != nil {
		return err
	}
	r.log.WithContext(ctx).Infof("gormDB: DeletePermission, id: %d", id)
	return nil
}

// GetUserRoles 获取用户角色
func (r *rbacRepo) GetUserRoles(ctx context.Context, userID int32) ([]*biz.Role, error) {
	roles, err := r.data.enforcer.GetRolesForUser(strconv.Itoa(int(userID)))
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return []*biz.Role{}, nil
	}

	var roleModels []*biz.Role
	err = r.data.gormDB.WithContext(ctx).Where("name IN ?", roles).Find(&roleModels).Error
	if err != nil {
		return nil, err
	}

	r.log.WithContext(ctx).Infof("gormDB: GetUserRoles, userID: %d, roles: %v", userID, roles)
	return roleModels, nil
}

// GetUserRoleNames 获取用户角色名称列表
func (r *rbacRepo) GetUserRoleNames(ctx context.Context, userID int32) ([]string, error) {
	roles, err := r.data.enforcer.GetRolesForUser(strconv.Itoa(int(userID)))
	if err != nil {
		return nil, err
	}
	r.log.WithContext(ctx).Infof("Casbin: GetUserRoleNames, userID: %d, roles: %v", userID, roles)
	return roles, nil
}

// AssignRoleToUser 为用户分配角色
func (r *rbacRepo) AssignRoleToUser(ctx context.Context, userID int32, roleID int32) error {
	// 先获取角色名称
	role, err := r.GetRole(ctx, roleID)
	if err != nil {
		return err
	}

	// 为用户添加角色
	success, err := r.data.enforcer.AddRoleForUser(strconv.Itoa(int(userID)), role.Name)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("failed to assign role to user")
	}

	r.log.WithContext(ctx).Infof("Casbin: AssignRoleToUser, userID: %d, roleID: %d, roleName: %s", userID, roleID, role.Name)
	return nil
}

// RemoveRoleFromUser 移除用户角色
func (r *rbacRepo) RemoveRoleFromUser(ctx context.Context, userID int32, roleID int32) error {
	// 先获取角色名称
	role, err := r.GetRole(ctx, roleID)
	if err != nil {
		return err
	}

	// 移除用户角色
	success, err := r.data.enforcer.DeleteRoleForUser(strconv.Itoa(int(userID)), role.Name)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("failed to remove role from user")
	}

	r.log.WithContext(ctx).Infof("Casbin: RemoveRoleFromUser, userID: %d, roleID: %d, roleName: %s", userID, roleID, role.Name)
	return nil
}

// CheckPermission 检查权限
func (r *rbacRepo) CheckPermission(ctx context.Context, userID int32, resource string, action string) (bool, error) {
	allowed, err := r.data.enforcer.Enforce(strconv.Itoa(int(userID)), resource, action)
	if err != nil {
		return false, err
	}

	r.log.WithContext(ctx).Infof("Casbin: CheckPermission, userID: %d, resource: %s, action: %s, allowed: %v", userID, resource, action, allowed)
	return allowed, nil
}
