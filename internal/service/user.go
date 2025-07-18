package service

import (
	"context"
	"strconv"

	pb "student/api/user/v1"
	"student/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type UserService struct {
	pb.UnimplementedUserServer

	user *biz.UserUsecase
	log  *log.Helper
}

func NewUserService(user *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{
		user: user,
		log:  log.NewHelper(logger),
	}
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	user, err := s.user.Get(ctx, req.Id)
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

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	s.log.Info("create user", req.Username, req.Email, req.Phone, req.Status, req.Age, req.Avatar)
	user, err := s.user.Create(ctx, &biz.UserForm{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		Status:   int(req.Status),
		Age:      int(req.Age),
		Avatar:   req.Avatar,
	})

	if err != nil {
		return nil, err
	}
	s.log.Info("create user", user.Message)
	return &pb.CreateUserReply{
		Message: user.Message,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	s.log.Info("update user", req.Id, req.Username, req.Email, req.Phone, req.Status, req.Age, req.Avatar)
	user, err := s.user.Update(ctx, req.Id, &biz.UserForm{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		Status:   int(req.Status),
		Age:      int(req.Age),
		Avatar:   req.Avatar,
	})

	if err != nil {
		return nil, err
	}
	s.log.Info("update user", user.Message)
	return &pb.UpdateUserReply{
		Message: user.Message,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	s.log.Info("delete user", req.Id)
	_, err := s.user.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	s.log.Info("delete user success")
	return &pb.DeleteUserReply{
		Message: "delete user success",
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersReply, error) {
	s.log.Info("list users")
	var err error
	var pageSize int
	var total int32
	var users []*biz.User

	if req.PageSize != "" {
		pageSize, err = strconv.Atoi(req.PageSize)
		if err != nil {
			return nil, err
		}
	} else {
		pageSize = 10
	}

	var page int
	if req.Page != "" {
		page, _ = strconv.Atoi(req.Page)
	} else {
		page = 1
	}

	users, total, err = s.user.List(ctx, int32(page), int32(pageSize), req.Username, req.Email)
	if err != nil {
		return nil, err
	}

	var data []*pb.Users
	for _, user := range users {
		data = append(data, &pb.Users{
			Id:        int32(user.ID),
			Username:  user.Username,
			Email:     user.Email,
			Phone:     user.Phone,
			Status:    int32(user.Status),
			Age:       int32(user.Age),
			Avatar:    user.Avatar,
			CreatedAt: user.CreatedAtStr, // 使用格式化后的字符串，避免每次转换
			UpdatedAt: user.UpdatedAtStr, // 使用格式化后的字符串，避免每次转换
		})
	}

	return &pb.ListUsersReply{
		Data:  data,
		Total: int32(total),
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	s.log.Info("user login", req.Username)

	loginForm := &biz.LoginForm{
		Username: req.Username,
		Password: req.Password,
	}

	loginResult, err := s.user.Login(ctx, loginForm)
	if err != nil {
		return nil, err
	}

	s.log.Info("login result", "success", loginResult.Success, "token_length", len(loginResult.Token), "token", loginResult.Token)

	reply := &pb.LoginReply{
		Success: loginResult.Success,
		Message: loginResult.Message,
		Token:   loginResult.Token,
	}

	if loginResult.Success && loginResult.User != nil {
		reply.UserInfo = &pb.UserInfo{
			Id:        int32(loginResult.User.ID),
			Username:  loginResult.User.Username,
			Email:     loginResult.User.Email,
			Phone:     loginResult.User.Phone,
			Status:    int32(loginResult.User.Status),
			Age:       int32(loginResult.User.Age),
			Avatar:    loginResult.User.Avatar,
			CreatedAt: loginResult.User.CreatedAtStr,
			UpdatedAt: loginResult.User.UpdatedAtStr,
		}
	}

	s.log.Info("final reply", "token_length", len(reply.Token), "token", reply.Token)

	return reply, nil
}
