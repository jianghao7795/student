package service

import (
	"context"

	pb "student/api/student/v1"
	"student/internal/student-service/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type StudentService struct {
	pb.UnimplementedStudentServer

	uc  *biz.StudentUsecase
	log *log.Helper
}

func NewStudentService(uc *biz.StudentUsecase, logger log.Logger) *StudentService {
	return &StudentService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *StudentService) CreateStudent(ctx context.Context, req *pb.CreateStudentRequest) (*pb.CreateStudentReply, error) {
	student := &biz.Student{
		Name:   req.Name,
		Info:   req.Info,
		Status: int(req.Status),
		Age:    int(req.Age),
	}

	_, err := s.uc.CreateStudent(ctx, student)
	if err != nil {
		return nil, err
	}

	return &pb.CreateStudentReply{
		Message: "创建学生成功",
	}, nil
}

func (s *StudentService) UpdateStudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.UpdateStudentReply, error) {
	student := &biz.Student{
		ID:     uint(req.Id),
		Name:   req.Name,
		Info:   req.Info,
		Status: int(req.Status),
		Age:    int(req.Age),
	}

	_, err := s.uc.UpdateStudent(ctx, student)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateStudentReply{
		Message: "更新学生成功",
	}, nil
}

func (s *StudentService) DeleteStudent(ctx context.Context, req *pb.DeleteStudentRequest) (*pb.DeleteStudentReply, error) {
	err := s.uc.DeleteStudent(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteStudentReply{
		Message: "删除成功",
	}, nil
}

func (s *StudentService) GetStudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.GetStudentReply, error) {
	student, err := s.uc.GetStudent(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetStudentReply{
		Id:        int32(student.ID),
		Name:      student.Name,
		Info:      student.Info,
		Status:    int32(student.Status),
		Age:       int32(student.Age),
		CreatedAt: student.CreatedAtStr,
		UpdatedAt: student.UpdatedAtStr,
	}, nil
}

func (s *StudentService) ListStudent(ctx context.Context, req *pb.ListStudentsRequest) (*pb.ListStudentsReply, error) {
	students, err := s.uc.ListStudent(ctx)
	if err != nil {
		return nil, err
	}

	var pbStudents []*pb.Students
	for _, student := range students {
		pbStudents = append(pbStudents, &pb.Students{
			Id:        int32(student.ID),
			Name:      student.Name,
			Info:      student.Info,
			Status:    int32(student.Status),
			Age:       int32(student.Age),
			CreatedAt: student.CreatedAtStr,
			UpdatedAt: student.UpdatedAtStr,
		})
	}

	return &pb.ListStudentsReply{
		Data:  pbStudents,
		Total: int32(len(pbStudents)),
	}, nil
}
