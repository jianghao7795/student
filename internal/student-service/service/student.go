package service

import (
	"context"
	"time"

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

// HealthCheck 健康检查
func (s *StudentService) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckReply, error) {
	return &pb.HealthCheckReply{
		Status:    "OK",
		Message:   "Student service is healthy",
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
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

func (s *StudentService) ListStudents(ctx context.Context, req *pb.ListStudentsRequest) (*pb.ListStudentsReply, error) {
	// 获取学生列表
	students, err := s.uc.ListStudent(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为proto消息
	var studentsProto []*pb.Students
	for _, student := range students {
		studentsProto = append(studentsProto, &pb.Students{
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
		Data:  studentsProto,
		Total: int32(len(studentsProto)), // 这里应该返回总数，暂时返回当前页数量
	}, nil
}
