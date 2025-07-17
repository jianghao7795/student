package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

const (
	// TimeFormat 时间格式常量
	TimeFormat = "2006-01-02 15:04:05"
)

// Student is a Student model.
type Student struct {
	// MODEL
	ID        uint
	Name      string
	Info      string
	Status    int
	Age       int
	CreatedAt *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`

	// 格式化后的时间字符串，避免每次转换
	CreatedAtStr string `gorm:"-" json:"created_at_str,omitempty"` // gorm:"-" 表示不映射到数据库
	UpdatedAtStr string `gorm:"-" json:"updated_at_str,omitempty"` // gorm:"-" 表示不映射到数据库
}

// FormatTimeFields 格式化时间字段
func (s *Student) FormatTimeFields() {
	if s.CreatedAt != nil {
		s.CreatedAtStr = s.CreatedAt.Format(TimeFormat)
	}
	if s.UpdatedAt != nil {
		s.UpdatedAtStr = s.UpdatedAt.Format(TimeFormat)
	}
}

// FormatTimeFieldsBatch 批量格式化时间字段
func FormatTimeFieldsBatch(students []*Student) {
	for _, stu := range students {
		stu.FormatTimeFields()
	}
}

type StudentForm struct {
	Name   string
	Info   string
	Status int
	Age    int
	ID     uint
	// 其他字段...
}

type CreateStudentMessage struct {
	ID      int32
	Message string
}

type UpdateStudentMessage struct {
	Message string
}

type DeleteStudentMessage struct {
	Message string
}

// 定义 Student 的操作接口
type StudentRepo interface {
	GetStudent(ctx context.Context, id int32) (*Student, error)
	CreateStudent(ctx context.Context, s *StudentForm) (*CreateStudentMessage, error)
	UpdateStudent(ctx context.Context, id int32, s *StudentForm) (*UpdateStudentMessage, error)
	DeleteStudent(ctx context.Context, id int32) (*DeleteStudentMessage, error)
	ListStudents(ctx context.Context, page int32, pageSize int32, name string) ([]*Student, int32, error)
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

// update student
func (uc *StudentUsecase) Update(ctx context.Context, id int32, s *StudentForm) (*UpdateStudentMessage, error) {
	return uc.repo.UpdateStudent(ctx, id, s)
}

// delete student
func (uc *StudentUsecase) Delete(ctx context.Context, id int32) (*DeleteStudentMessage, error) {
	return uc.repo.DeleteStudent(ctx, id)
}

// get list student
func (uc *StudentUsecase) List(ctx context.Context, page int32, pageSize int32, name string) ([]*Student, int32, error) {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	return uc.repo.ListStudents(ctx, page, pageSize, name)
}
