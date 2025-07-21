package service

import (
	"context"

	v1 "student/api/rbac/v1"
	"student/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type RBACService struct {
	v1.UnimplementedRBACServiceServer

	rbacUC *biz.RBACUsecase
	log    *log.Helper
}

func NewRBACService(rbacUC *biz.RBACUsecase, logger log.Logger) *RBACService {
	return &RBACService{
		rbacUC: rbacUC,
		log:    log.NewHelper(logger),
	}
}

// 角色相关服务方法
func (s *RBACService) GetRole(ctx context.Context, req *v1.GetRoleRequest) (*v1.GetRoleResponse, error) {
	role, err := s.rbacUC.GetRole(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 转换为protobuf消息
	roleProto := &v1.Role{
		Id:          uint32(role.ID),
		Name:        role.Name,
		Description: role.Description,
		Status:      int32(role.Status),
		CreatedAt:   role.CreatedAtStr,
		UpdatedAt:   role.UpdatedAtStr,
	}

	return &v1.GetRoleResponse{
		Role: roleProto,
	}, nil
}

func (s *RBACService) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*v1.CreateRoleResponse, error) {
	roleForm := &biz.RoleForm{
		Name:        req.Name,
		Description: req.Description,
		Status:      int(req.Status),
	}

	role, err := s.rbacUC.CreateRole(ctx, roleForm)
	if err != nil {
		return nil, err
	}

	// 转换为protobuf消息
	roleProto := &v1.Role{
		Id:          uint32(role.ID),
		Name:        role.Name,
		Description: role.Description,
		Status:      int32(role.Status),
		CreatedAt:   role.CreatedAtStr,
		UpdatedAt:   role.UpdatedAtStr,
	}

	return &v1.CreateRoleResponse{
		Role: roleProto,
	}, nil
}

func (s *RBACService) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.UpdateRoleResponse, error) {
	roleForm := &biz.RoleForm{
		Name:        req.Name,
		Description: req.Description,
		Status:      int(req.Status),
	}

	role, err := s.rbacUC.UpdateRole(ctx, req.Id, roleForm)
	if err != nil {
		return nil, err
	}

	// 转换为protobuf消息
	roleProto := &v1.Role{
		Id:          uint32(role.ID),
		Name:        role.Name,
		Description: role.Description,
		Status:      int32(role.Status),
		CreatedAt:   role.CreatedAtStr,
		UpdatedAt:   role.UpdatedAtStr,
	}

	return &v1.UpdateRoleResponse{
		Role: roleProto,
	}, nil
}

func (s *RBACService) DeleteRole(ctx context.Context, req *v1.DeleteRoleRequest) (*v1.DeleteRoleResponse, error) {
	err := s.rbacUC.DeleteRole(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteRoleResponse{
		Message: "角色删除成功",
	}, nil
}

func (s *RBACService) ListRoles(ctx context.Context, req *v1.ListRolesRequest) (*v1.ListRolesResponse, error) {
	roles, total, err := s.rbacUC.ListRoles(ctx, req.Page, req.PageSize, req.Name)
	if err != nil {
		return nil, err
	}

	// 转换为protobuf消息
	var rolesProto []*v1.Role
	for _, role := range roles {
		roleProto := &v1.Role{
			Id:          uint32(role.ID),
			Name:        role.Name,
			Description: role.Description,
			Status:      int32(role.Status),
			CreatedAt:   role.CreatedAtStr,
			UpdatedAt:   role.UpdatedAtStr,
		}
		rolesProto = append(rolesProto, roleProto)
	}

	return &v1.ListRolesResponse{
		Roles: rolesProto,
		Total: total,
	}, nil
}

// 权限相关服务方法
func (s *RBACService) GetPermission(ctx context.Context, req *v1.GetPermissionRequest) (*v1.GetPermissionResponse, error) {
	permission, err := s.rbacUC.GetPermission(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 转换为protobuf消息
	permissionProto := &v1.Permission{
		Id:          uint32(permission.ID),
		Name:        permission.Name,
		Resource:    permission.Resource,
		Action:      permission.Action,
		Description: permission.Description,
		Status:      int32(permission.Status),
		CreatedAt:   permission.CreatedAtStr,
		UpdatedAt:   permission.UpdatedAtStr,
	}

	return &v1.GetPermissionResponse{
		Permission: permissionProto,
	}, nil
}

func (s *RBACService) CreatePermission(ctx context.Context, req *v1.CreatePermissionRequest) (*v1.CreatePermissionResponse, error) {
	permissionForm := &biz.PermissionForm{
		Name:        req.Name,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
		Status:      int(req.Status),
	}

	permission, err := s.rbacUC.CreatePermission(ctx, permissionForm)
	if err != nil {
		return nil, err
	}

	// 转换为protobuf消息
	permissionProto := &v1.Permission{
		Id:          uint32(permission.ID),
		Name:        permission.Name,
		Resource:    permission.Resource,
		Action:      permission.Action,
		Description: permission.Description,
		Status:      int32(permission.Status),
		CreatedAt:   permission.CreatedAtStr,
		UpdatedAt:   permission.UpdatedAtStr,
	}

	return &v1.CreatePermissionResponse{
		Permission: permissionProto,
	}, nil
}

func (s *RBACService) UpdatePermission(ctx context.Context, req *v1.UpdatePermissionRequest) (*v1.UpdatePermissionResponse, error) {
	permissionForm := &biz.PermissionForm{
		Name:        req.Name,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
		Status:      int(req.Status),
	}

	permission, err := s.rbacUC.UpdatePermission(ctx, req.Id, permissionForm)
	if err != nil {
		return nil, err
	}

	// 转换为protobuf消息
	permissionProto := &v1.Permission{
		Id:          uint32(permission.ID),
		Name:        permission.Name,
		Resource:    permission.Resource,
		Action:      permission.Action,
		Description: permission.Description,
		Status:      int32(permission.Status),
		CreatedAt:   permission.CreatedAtStr,
		UpdatedAt:   permission.UpdatedAtStr,
	}

	return &v1.UpdatePermissionResponse{
		Permission: permissionProto,
	}, nil
}

func (s *RBACService) DeletePermission(ctx context.Context, req *v1.DeletePermissionRequest) (*v1.DeletePermissionResponse, error) {
	err := s.rbacUC.DeletePermission(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.DeletePermissionResponse{
		Message: "权限删除成功",
	}, nil
}

func (s *RBACService) ListPermissions(ctx context.Context, req *v1.ListPermissionsRequest) (*v1.ListPermissionsResponse, error) {
	permissions, total, err := s.rbacUC.ListPermissions(ctx, req.Page, req.PageSize, req.Name, req.Resource)
	if err != nil {
		return nil, err
	}

	// 转换为protobuf消息
	var permissionsProto []*v1.Permission
	for _, permission := range permissions {
		permissionProto := &v1.Permission{
			Id:          uint32(permission.ID),
			Name:        permission.Name,
			Resource:    permission.Resource,
			Action:      permission.Action,
			Description: permission.Description,
			Status:      int32(permission.Status),
			CreatedAt:   permission.CreatedAtStr,
			UpdatedAt:   permission.UpdatedAtStr,
		}
		permissionsProto = append(permissionsProto, permissionProto)
	}

	return &v1.ListPermissionsResponse{
		Permissions: permissionsProto,
		Total:       total,
	}, nil
}

// 用户角色相关服务方法
func (s *RBACService) GetUserRoles(ctx context.Context, req *v1.GetUserRolesRequest) (*v1.GetUserRolesResponse, error) {
	userRoles, err := s.rbacUC.GetUserRoles(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	// 转换为protobuf消息
	var userRolesProto []*v1.UserRole
	for _, userRole := range userRoles {
		userRoleProto := &v1.UserRole{
			Id:        uint32(userRole.ID),
			UserId:    uint32(userRole.UserID),
			RoleId:    uint32(userRole.RoleID),
			CreatedAt: userRole.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: userRole.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if userRole.Role != nil {
			userRoleProto.Role = &v1.Role{
				Id:          uint32(userRole.Role.ID),
				Name:        userRole.Role.Name,
				Description: userRole.Role.Description,
				Status:      int32(userRole.Role.Status),
				CreatedAt:   userRole.Role.CreatedAtStr,
				UpdatedAt:   userRole.Role.UpdatedAtStr,
			}
		}
		userRolesProto = append(userRolesProto, userRoleProto)
	}

	return &v1.GetUserRolesResponse{
		UserRoles: userRolesProto,
	}, nil
}

func (s *RBACService) AssignUserRole(ctx context.Context, req *v1.AssignUserRoleRequest) (*v1.AssignUserRoleResponse, error) {
	err := s.rbacUC.AssignUserRole(ctx, req.UserId, req.RoleId)
	if err != nil {
		return nil, err
	}

	return &v1.AssignUserRoleResponse{
		Message: "用户角色分配成功",
	}, nil
}

func (s *RBACService) RemoveUserRole(ctx context.Context, req *v1.RemoveUserRoleRequest) (*v1.RemoveUserRoleResponse, error) {
	err := s.rbacUC.RemoveUserRole(ctx, req.UserId, req.RoleId)
	if err != nil {
		return nil, err
	}

	return &v1.RemoveUserRoleResponse{
		Message: "用户角色移除成功",
	}, nil
}

// 角色权限相关服务方法
func (s *RBACService) GetRolePermissions(ctx context.Context, req *v1.GetRolePermissionsRequest) (*v1.GetRolePermissionsResponse, error) {
	rolePermissions, err := s.rbacUC.GetRolePermissions(ctx, req.RoleId)
	if err != nil {
		return nil, err
	}

	// 转换为protobuf消息
	var rolePermissionsProto []*v1.RolePermission
	for _, rolePermission := range rolePermissions {
		rolePermissionProto := &v1.RolePermission{
			Id:           uint32(rolePermission.ID),
			RoleId:       uint32(rolePermission.RoleID),
			PermissionId: uint32(rolePermission.PermissionID),
			CreatedAt:    rolePermission.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    rolePermission.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if rolePermission.Permission != nil {
			rolePermissionProto.Permission = &v1.Permission{
				Id:          uint32(rolePermission.Permission.ID),
				Name:        rolePermission.Permission.Name,
				Resource:    rolePermission.Permission.Resource,
				Action:      rolePermission.Permission.Action,
				Description: rolePermission.Permission.Description,
				Status:      int32(rolePermission.Permission.Status),
				CreatedAt:   rolePermission.Permission.CreatedAtStr,
				UpdatedAt:   rolePermission.Permission.UpdatedAtStr,
			}
		}
		rolePermissionsProto = append(rolePermissionsProto, rolePermissionProto)
	}

	return &v1.GetRolePermissionsResponse{
		RolePermissions: rolePermissionsProto,
	}, nil
}

func (s *RBACService) AssignRolePermission(ctx context.Context, req *v1.AssignRolePermissionRequest) (*v1.AssignRolePermissionResponse, error) {
	err := s.rbacUC.AssignRolePermission(ctx, req.RoleId, req.PermissionId)
	if err != nil {
		return nil, err
	}

	return &v1.AssignRolePermissionResponse{
		Message: "角色权限分配成功",
	}, nil
}

func (s *RBACService) RemoveRolePermission(ctx context.Context, req *v1.RemoveRolePermissionRequest) (*v1.RemoveRolePermissionResponse, error) {
	err := s.rbacUC.RemoveRolePermission(ctx, req.RoleId, req.PermissionId)
	if err != nil {
		return nil, err
	}

	return &v1.RemoveRolePermissionResponse{
		Message: "角色权限移除成功",
	}, nil
}

// 权限检查服务方法
func (s *RBACService) CheckPermission(ctx context.Context, req *v1.CheckPermissionRequest) (*v1.CheckPermissionResponse, error) {
	hasPermission, err := s.rbacUC.CheckPermission(ctx, req.User, req.Resource, req.Action)
	if err != nil {
		return nil, err
	}

	return &v1.CheckPermissionResponse{
		HasPermission: hasPermission,
	}, nil
}
