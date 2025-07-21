package data

import (
	"context"
	"fmt"
	"strconv"

	"student/internal/biz"
	"student/internal/data/errors"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type rbacRepo struct {
	data     *Data
	log      *log.Helper
	enforcer *casbin.Enforcer
}

func NewRBACRepo(data *Data, logger log.Logger, modelPath string) biz.RBACRepo {
	// 创建Casbin适配器
	adapter, err := gormadapter.NewAdapterByDB(data.gormDB)
	if err != nil {
		log.NewHelper(logger).Error("failed to create casbin adapter", err)
		return nil
	}

	// 创建Casbin执行器
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		log.NewHelper(logger).Error("failed to create casbin enforcer", err)
		return nil
	}

	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		log.NewHelper(logger).Error("failed to load casbin policy", err)
	}

	return &rbacRepo{
		data:     data,
		log:      log.NewHelper(logger),
		enforcer: enforcer,
	}
}

// 角色相关方法实现
func (r *rbacRepo) GetRole(ctx context.Context, id int32) (*biz.Role, error) {
	var role biz.Role
	err := r.data.gormDB.WithContext(ctx).First(&role, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Error404()
		}
		return nil, errors.Error400(err)
	}

	role.FormatTimeFields()
	return &role, nil
}

func (r *rbacRepo) CreateRole(ctx context.Context, roleForm *biz.RoleForm) (*biz.Role, error) {
	role := biz.Role{
		Name:        roleForm.Name,
		Description: roleForm.Description,
		Status:      roleForm.Status,
	}

	err := r.data.gormDB.WithContext(ctx).Create(&role).Error
	if err != nil {
		return nil, errors.Error400(err)
	}

	role.FormatTimeFields()
	return &role, nil
}

func (r *rbacRepo) UpdateRole(ctx context.Context, id int32, roleForm *biz.RoleForm) (*biz.Role, error) {
	var role biz.Role
	err := r.data.gormDB.WithContext(ctx).First(&role, id).Error
	if err != nil {
		return nil, errors.Error404()
	}

	role.Name = roleForm.Name
	role.Description = roleForm.Description
	role.Status = roleForm.Status

	err = r.data.gormDB.WithContext(ctx).Save(&role).Error
	if err != nil {
		return nil, errors.Error400(err)
	}

	role.FormatTimeFields()
	return &role, nil
}

func (r *rbacRepo) DeleteRole(ctx context.Context, id int32) error {
	var role biz.Role
	err := r.data.gormDB.WithContext(ctx).First(&role, id).Error
	if err != nil {
		return errors.Error404()
	}

	return r.data.gormDB.WithContext(ctx).Delete(&role).Error
}

func (r *rbacRepo) ListRoles(ctx context.Context, page int32, pageSize int32, name string) ([]*biz.Role, int32, error) {
	var roles []*biz.Role
	var total int64

	query := r.data.gormDB.WithContext(ctx).Model(&biz.Role{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, errors.Error400(err)
	}

	query = r.data.gormDB.WithContext(ctx).Model(&biz.Role{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err = query.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Order("id desc").Find(&roles).Error
	if err != nil {
		return nil, 0, errors.Error400(err)
	}

	for _, role := range roles {
		role.FormatTimeFields()
	}

	return roles, int32(total), nil
}

func (r *rbacRepo) GetRoleByName(ctx context.Context, name string) (*biz.Role, error) {
	var role biz.Role
	err := r.data.gormDB.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Error404()
		}
		return nil, errors.Error400(err)
	}

	role.FormatTimeFields()
	return &role, nil
}

// 权限相关方法实现
func (r *rbacRepo) GetPermission(ctx context.Context, id int32) (*biz.Permission, error) {
	var permission biz.Permission
	err := r.data.gormDB.WithContext(ctx).First(&permission, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Error404()
		}
		return nil, errors.Error400(err)
	}

	permission.FormatTimeFields()
	return &permission, nil
}

func (r *rbacRepo) CreatePermission(ctx context.Context, permissionForm *biz.PermissionForm) (*biz.Permission, error) {
	permission := biz.Permission{
		Name:        permissionForm.Name,
		Resource:    permissionForm.Resource,
		Action:      permissionForm.Action,
		Description: permissionForm.Description,
		Status:      permissionForm.Status,
	}

	err := r.data.gormDB.WithContext(ctx).Create(&permission).Error
	if err != nil {
		return nil, errors.Error400(err)
	}

	permission.FormatTimeFields()
	return &permission, nil
}

func (r *rbacRepo) UpdatePermission(ctx context.Context, id int32, permissionForm *biz.PermissionForm) (*biz.Permission, error) {
	var permission biz.Permission
	err := r.data.gormDB.WithContext(ctx).First(&permission, id).Error
	if err != nil {
		return nil, errors.Error404()
	}

	permission.Name = permissionForm.Name
	permission.Resource = permissionForm.Resource
	permission.Action = permissionForm.Action
	permission.Description = permissionForm.Description
	permission.Status = permissionForm.Status

	err = r.data.gormDB.WithContext(ctx).Save(&permission).Error
	if err != nil {
		return nil, errors.Error400(err)
	}

	permission.FormatTimeFields()
	return &permission, nil
}

func (r *rbacRepo) DeletePermission(ctx context.Context, id int32) error {
	var permission biz.Permission
	err := r.data.gormDB.WithContext(ctx).First(&permission, id).Error
	if err != nil {
		return errors.Error404()
	}

	return r.data.gormDB.WithContext(ctx).Delete(&permission).Error
}

func (r *rbacRepo) ListPermissions(ctx context.Context, page int32, pageSize int32, name, resource string) ([]*biz.Permission, int32, error) {
	var permissions []*biz.Permission
	var total int64

	query := r.data.gormDB.WithContext(ctx).Model(&biz.Permission{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if resource != "" {
		query = query.Where("resource LIKE ?", "%"+resource+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, errors.Error400(err)
	}

	query = r.data.gormDB.WithContext(ctx).Model(&biz.Permission{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if resource != "" {
		query = query.Where("resource LIKE ?", "%"+resource+"%")
	}

	err = query.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Order("id desc").Find(&permissions).Error
	if err != nil {
		return nil, 0, errors.Error400(err)
	}

	for _, permission := range permissions {
		permission.FormatTimeFields()
	}

	return permissions, int32(total), nil
}

func (r *rbacRepo) GetPermissionByResourceAction(ctx context.Context, resource, action string) (*biz.Permission, error) {
	var permission biz.Permission
	err := r.data.gormDB.WithContext(ctx).Where("resource = ? AND action = ?", resource, action).First(&permission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Error404()
		}
		return nil, errors.Error400(err)
	}

	permission.FormatTimeFields()
	return &permission, nil
}

// 用户角色相关方法实现
func (r *rbacRepo) GetUserRoles(ctx context.Context, userID int32) ([]*biz.UserRole, error) {
	var userRoles []*biz.UserRole
	err := r.data.gormDB.WithContext(ctx).Preload("Role").Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, errors.Error400(err)
	}
	return userRoles, nil
}

func (r *rbacRepo) AssignUserRole(ctx context.Context, userID, roleID int32) error {
	// 检查用户和角色是否存在
	var user biz.User
	err := r.data.gormDB.WithContext(ctx).First(&user, userID).Error
	if err != nil {
		return errors.Error404()
	}

	var role biz.Role
	err = r.data.gormDB.WithContext(ctx).First(&role, roleID).Error
	if err != nil {
		return errors.Error404()
	}

	// 创建用户角色关联
	userRole := biz.UserRole{
		UserID: uint(userID),
		RoleID: uint(roleID),
	}

	err = r.data.gormDB.WithContext(ctx).Create(&userRole).Error
	if err != nil {
		return errors.Error400(err)
	}

	// 同步到Casbin
	_, err = r.enforcer.AddRoleForUser(strconv.Itoa(int(userID)), role.Name)
	return err
}

func (r *rbacRepo) RemoveUserRole(ctx context.Context, userID, roleID int32) error {
	// 获取角色名称
	var role biz.Role
	err := r.data.gormDB.WithContext(ctx).First(&role, roleID).Error
	if err != nil {
		return errors.Error404()
	}

	// 删除用户角色关联
	err = r.data.gormDB.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&biz.UserRole{}).Error
	if err != nil {
		return errors.Error400(err)
	}

	// 同步到Casbin
	_, err = r.enforcer.RemoveFilteredGroupingPolicy(0, strconv.Itoa(int(userID)), role.Name)
	return err
}

func (r *rbacRepo) GetUserRoleNames(ctx context.Context, userID int32) ([]string, error) {
	return r.enforcer.GetRolesForUser(strconv.Itoa(int(userID)))
}

// 角色权限相关方法实现
func (r *rbacRepo) GetRolePermissions(ctx context.Context, roleID int32) ([]*biz.RolePermission, error) {
	var rolePermissions []*biz.RolePermission
	err := r.data.gormDB.WithContext(ctx).Preload("Permission").Where("role_id = ?", roleID).Find(&rolePermissions).Error
	if err != nil {
		return nil, errors.Error400(err)
	}
	return rolePermissions, nil
}

func (r *rbacRepo) AssignRolePermission(ctx context.Context, roleID, permissionID int32) error {
	// 检查角色和权限是否存在
	var role biz.Role
	err := r.data.gormDB.WithContext(ctx).First(&role, roleID).Error
	if err != nil {
		return errors.Error404()
	}

	var permission biz.Permission
	err = r.data.gormDB.WithContext(ctx).First(&permission, permissionID).Error
	if err != nil {
		return errors.Error404()
	}

	// 创建角色权限关联
	rolePermission := biz.RolePermission{
		RoleID:       uint(roleID),
		PermissionID: uint(permissionID),
	}

	err = r.data.gormDB.WithContext(ctx).Create(&rolePermission).Error
	if err != nil {
		return errors.Error400(err)
	}

	// 同步到Casbin
	_, err = r.enforcer.AddPolicy(role.Name, permission.Resource, permission.Action)
	return err
}

func (r *rbacRepo) RemoveRolePermission(ctx context.Context, roleID, permissionID int32) error {
	// 获取角色和权限信息
	var role biz.Role
	err := r.data.gormDB.WithContext(ctx).First(&role, roleID).Error
	if err != nil {
		return errors.Error404()
	}

	var permission biz.Permission
	err = r.data.gormDB.WithContext(ctx).First(&permission, permissionID).Error
	if err != nil {
		return errors.Error404()
	}

	// 删除角色权限关联
	err = r.data.gormDB.WithContext(ctx).Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&biz.RolePermission{}).Error
	if err != nil {
		return errors.Error400(err)
	}

	// 同步到Casbin
	_, err = r.enforcer.RemovePolicy(role.Name, permission.Resource, permission.Action)
	return err
}

func (r *rbacRepo) GetRolePermissionNames(ctx context.Context, roleID int32) ([]string, error) {
	var role biz.Role
	err := r.data.gormDB.WithContext(ctx).First(&role, roleID).Error
	if err != nil {
		return nil, errors.Error404()
	}

	policies, _ := r.enforcer.GetFilteredPolicy(0, role.Name)
	var permissions []string
	for _, policy := range policies {
		if len(policy) >= 3 {
			permissions = append(permissions, fmt.Sprintf("%s:%s", policy[1], policy[2]))
		}
	}
	return permissions, nil
}

// Casbin相关方法实现
func (r *rbacRepo) LoadPolicy(ctx context.Context) error {
	return r.enforcer.LoadPolicy()
}

func (r *rbacRepo) SavePolicy(ctx context.Context) error {
	return r.enforcer.SavePolicy()
}

func (r *rbacRepo) AddPolicy(ctx context.Context, sub, obj, act string) error {
	_, err := r.enforcer.AddPolicy(sub, obj, act)
	return err
}

func (r *rbacRepo) RemovePolicy(ctx context.Context, sub, obj, act string) error {
	_, err := r.enforcer.RemovePolicy(sub, obj, act)
	return err
}

func (r *rbacRepo) AddRoleForUser(ctx context.Context, user, role string) error {
	_, err := r.enforcer.AddRoleForUser(user, role)
	return err
}

func (r *rbacRepo) RemoveRoleForUser(ctx context.Context, user, role string) error {
	_, err := r.enforcer.RemoveFilteredGroupingPolicy(0, user, role)
	return err
}

func (r *rbacRepo) GetRolesForUser(ctx context.Context, user string) ([]string, error) {
	return r.enforcer.GetRolesForUser(user)
}

func (r *rbacRepo) GetPermissionsForUser(ctx context.Context, user string) ([][]string, error) {
	return r.enforcer.GetPermissionsForUser(user)
}
