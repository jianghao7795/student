package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Student is a Student model.
type Student struct {
	MODEL
	Name   string
	Info   string
	Status int
	Age    int
}

type StudentForm struct {
	Name   string
	Info   string
	Status int
	Age    int
	MODEL
}

type CreateStudentMessage struct {
	Message string
}

// 定义 Student 的操作接口
type StudentRepo interface {
	GetStudent(ctx context.Context, id int32) (*Student, error)
	CreateStudent(ctx context.Context, s *StudentForm) (*CreateStudentMessage, error)
}

type StudentUsecase struct {
	repo StudentRepo
	log  *log.Helper
}

// 初始化 StudentUsecase
func NewStudentUsecase(repo StudentRepo, logger log.Logger) *StudentUsecase {
	return &StudentUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// 通过 id 获取 student 信息
func (uc *StudentUsecase) Get(ctx context.Context, id int32) (*Student, error) {
	uc.log.Info("get student by id", id)
	return uc.repo.GetStudent(ctx, id)
}

// create student
func (uc *StudentUsecase) Create(ctx context.Context, s *StudentForm) (*CreateStudentMessage, error) {
	uc.log.Info("create student", s)
	return uc.repo.CreateStudent(ctx, s)
}
