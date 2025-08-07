package service

import (
	"context"
	"strconv"

	pb "student/api/user/v1"
	"student/internal/pkg/jwt"
	"student/internal/user-service/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type UserService struct {
	pb.UnimplementedUserServer

	uc      *biz.UserUsecase
	log     *log.Helper
	jwtUtil *jwt.JWTUtil
}

func NewUserService(uc *biz.UserUsecase, logger log.Logger, jwtUtil *jwt.JWTUtil) *UserService {
	return &UserService{
		uc:      uc,
		log:     log.NewHelper(logger),
		jwtUtil: jwtUtil,
	}
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	user, err := s.uc.GetUser(ctx, req.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	s.log.Info("get user", user.CreatedAt, user.UpdatedAt)
	userReply := pb.GetUserReply{
		Id:        int32(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Status:    int32(user.Status),
		Age:       int32(user.Age),
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAtStr, // 使用格式化后的字符串，避免每次转换
		UpdatedAt: user.UpdatedAtStr, // 使用格式化后的字符串，避免每次转换
	}
	return &userReply, err
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersReply, error) {
	// 转换字符串参数为数字
	page := int32(1)
	pageSize := int32(10)

	if req.Page != "" {
		if p, err := strconv.ParseInt(req.Page, 10, 32); err == nil {
			page = int32(p)
		}
	}

	if req.PageSize != "" {
		if ps, err := strconv.ParseInt(req.PageSize, 10, 32); err == nil {
			pageSize = int32(ps)
		}
	}

	bizReq := &biz.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
		Username: req.Username,
		Email:    req.Email,
		Status:   0, // 默认值，因为proto中没有status字段
	}

	reply, err := s.uc.ListUsers(ctx, bizReq)
	if err != nil {
		return nil, err
	}

	var users []*pb.Users
	for _, user := range reply.Users {
		users = append(users, &pb.Users{
			Id:        int32(user.ID),
			Username:  user.Username,
			Email:     user.Email,
			Phone:     user.Phone,
			Status:    int32(user.Status),
			Age:       int32(user.Age),
			Avatar:    user.Avatar,
			CreatedAt: user.CreatedAtStr,
			UpdatedAt: user.UpdatedAtStr,
		})
	}

	return &pb.ListUsersReply{
		Total: reply.Total,
		Data:  users,
	}, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	user := &biz.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		Status:   int(req.Status),
		Age:      int(req.Age),
		Avatar:   req.Avatar,
	}

	_, err := s.uc.CreateUser(ctx, user)
	if err != nil {
		return &pb.CreateUserReply{
			Message: "创建用户失败",
		}, err
	}

	return &pb.CreateUserReply{
		Message: "创建用户成功",
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	user := &biz.User{
		ID:       uint(req.Id),
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		Status:   int(req.Status),
		Age:      int(req.Age),
		Avatar:   req.Avatar,
	}

	_, err := s.uc.UpdateUser(ctx, user)
	if err != nil {
		return &pb.UpdateUserReply{
			Message: "更新用户失败",
		}, err
	}

	return &pb.UpdateUserReply{
		Message: "更新用户成功",
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	err := s.uc.DeleteUser(ctx, req.Id)
	if err != nil {
		return &pb.DeleteUserReply{
			Message: "删除失败",
		}, err
	}

	return &pb.DeleteUserReply{
		Message: "删除成功",
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	s.log.Info("user login", req.Username, req.Password)

	loginForm := &biz.LoginForm{
		Username: req.Username,
		Password: req.Password,
	}

	result, err := s.uc.Login(ctx, loginForm)
	if err != nil {
		return &pb.LoginReply{
			Success: false,
			Message: "登录失败，请稍后重试",
		}, err
	}

	var userInfo *pb.UserInfo
	if result.User != nil {
		userInfo = &pb.UserInfo{
			Id:        int32(result.User.ID),
			Username:  result.User.Username,
			Email:     result.User.Email,
			Phone:     result.User.Phone,
			Status:    int32(result.User.Status),
			Age:       int32(result.User.Age),
			Avatar:    result.User.Avatar,
			CreatedAt: result.User.CreatedAtStr,
			UpdatedAt: result.User.UpdatedAtStr,
		}
	}

	return &pb.LoginReply{
		Success:  result.Success,
		Message:  result.Message,
		Token:    result.Token,
		UserInfo: userInfo,
	}, nil
}

func (s *UserService) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.GetMeReply, error) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Value("user_id").(int64)
	if !exists {
		// 尝试从请求中获取
		userIDStr, exists := ctx.Value("user_id").(string)
		if !exists {
			return &pb.GetMeReply{
				Success: false,
				Message: "未找到用户ID",
			}, nil
		}
		var err error
		userID, err = strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return &pb.GetMeReply{
				Success: false,
				Message: "用户ID格式错误",
			}, nil
		}
	}

	result, err := s.uc.GetMe(ctx, uint(userID))
	if err != nil {
		return &pb.GetMeReply{
			Success: false,
			Message: "获取用户信息失败",
		}, err
	}

	var userInfo *pb.UserInfo
	if result.User != nil {
		userInfo = &pb.UserInfo{
			Id:        int32(result.User.ID),
			Username:  result.User.Username,
			Email:     result.User.Email,
			Phone:     result.User.Phone,
			Status:    int32(result.User.Status),
			Age:       int32(result.User.Age),
			Avatar:    result.User.Avatar,
			CreatedAt: result.User.CreatedAtStr,
			UpdatedAt: result.User.UpdatedAtStr,
		}
	}

	return &pb.GetMeReply{
		Success:  result.Success,
		Message:  result.Message,
		UserInfo: userInfo,
	}, nil
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	s.log.Info("user register", req.Username, req.Email, req.Age)

	registerForm := &biz.RegisterForm{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Age:      int(req.Age),
		Avatar:   req.Avatar,
	}

	result, err := s.uc.Register(ctx, registerForm)
	if err != nil {
		return &pb.RegisterReply{
			Success: false,
			Message: "注册失败，请稍后重试",
		}, err
	}

	var userInfo *pb.UserInfo
	if result.User != nil {
		userInfo = &pb.UserInfo{
			Id:        int32(result.User.ID),
			Username:  result.User.Username,
			Email:     result.User.Email,
			Phone:     result.User.Phone,
			Status:    int32(result.User.Status),
			Age:       int32(result.User.Age),
			Avatar:    result.User.Avatar,
			CreatedAt: result.User.CreatedAtStr,
			UpdatedAt: result.User.UpdatedAtStr,
		}
	}

	return &pb.RegisterReply{
		Success:  result.Success,
		Message:  result.Message,
		UserInfo: userInfo,
	}, nil
}
