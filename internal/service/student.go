package service

import (
	"context"
	"strconv"

	pb "student/api/student/v1"
	"student/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
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
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	s.log.Info("get student", stu.CreatedAt, stu.UpdatedAt)
	student := pb.GetStudentReply{
		Name:      stu.Name,
		Info:      stu.Info,
		Status:    int32(stu.Status),
		Age:       int32(stu.Age),
		CreatedAt: stu.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: stu.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	// err = gconv.Struct(stu, &student)
	return &student, err
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

func (s *StudentService) UpdateStudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.UpdateStudentReply, error) {
	s.log.Info("update student", req.Id, req.Name, req.Age, req.Status, req.Info)
	stu, err := s.student.Update(ctx, req.Id, &biz.StudentForm{
		Name:   req.Name,
		Info:   req.Info,
		Status: int(req.Status),
		Age:    int(req.Age),
	})

	if err != nil {
		return nil, err
	}
	s.log.Info("update student", stu.Message)
	return &pb.UpdateStudentReply{
		Message: stu.Message,
	}, nil
}

func (s *StudentService) DeleteStudent(ctx context.Context, req *pb.DeleteStudentRequest) (*pb.DeleteStudentReply, error) {
	s.log.Info("delete student", req.Id)
	_, err := s.student.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	s.log.Info("delete student success")
	return &pb.DeleteStudentReply{
		Message: "delete student success",
	}, nil
}

func (s *StudentService) ListStudents(ctx context.Context, req *pb.ListStudentsRequest) (*pb.ListStudentsReply, error) {
	s.log.Info("list student")
	var err error
	var pageSize int
	var total int32
	var students []*biz.Student
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
		// if err != nil {
		// return nil, err
		// }
	} else {
		page = 1
	}
	students, total, err = s.student.List(ctx, int32(page), int32(pageSize), req.Name)
	if err != nil {
		return nil, err
	}
	var data []*pb.Students
	for _, stu := range students {
		data = append(data, &pb.Students{
			Id:        int32(stu.ID),
			Name:      stu.Name,
			Info:      stu.Info,
			Status:    int32(stu.Status),
			Age:       int32(stu.Age),
			CreatedAt: stu.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: stu.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &pb.ListStudentsReply{
		Data:  data,
		Total: int32(total),
	}, nil
}
