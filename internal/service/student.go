package service

import (
	"context"
	"fmt"

	pb "student/api/student/v1"
	"student/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type StudentService struct {
	pb.UnimplementedStudentServer

	student *biz.StudentUsecase
	log     *log.Helper
}

func NewStudentService(student *biz.StudentUsecase, logger log.Logger) *StudentService {
	return &StudentService{
		student: student,
		log:     log.NewHelper(logger),
	}
}

func (s *StudentService) GetStudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.GetStudentReply, error) {
	stu, err := s.student.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	s.log.Info("get student", stu.CreatedAt, stu.UpdatedAt)
	fmt.Println("bbs to bbs", stu.CreatedAt, stu.UpdatedAt)
	return &pb.GetStudentReply{
		Id:     int32(stu.ID),
		Name:   stu.Name,
		Status: int32(stu.Status),
	}, nil
}

func (s *StudentService) CreateStudent(ctx context.Context, req *pb.CreateStudentRequest) (*pb.CreateStudentReply, error) {
	s.log.Info("create student", req.Name, req.Age, req.Status, req.Info)
	stu, err := s.student.Create(ctx, &biz.StudentForm{
		Name:   req.Name,
		Info:   req.Info,
		Status: int(req.Status),
		Age:    int(req.Age),
	})

	if err != nil {
		return nil, err
	}
	s.log.Info("create student", stu.Message)
	return &pb.CreateStudentReply{
		Message: stu.Message,
	}, nil
}
