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
}

// StudentRepo is a Student repo.
type StudentRepo interface {
	Save(context.Context, *Student) (*Student, error)
	Update(context.Context, *Student) (*Student, error)
	FindByID(context.Context, int32) (*Student, error)
	ListAll(context.Context) ([]*Student, error)
	Delete(context.Context, int32) error
}

// StudentUsecase is a Student usecase.
type StudentUsecase struct {
	repo StudentRepo
	log  *log.Helper
}

// NewStudentUsecase new a Student usecase.
func NewStudentUsecase(repo StudentRepo, logger log.Logger) *StudentUsecase {
	return &StudentUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateStudent creates a Student, and returns the new Student.
func (uc *StudentUsecase) CreateStudent(ctx context.Context, s *Student) (*Student, error) {
	uc.log.WithContext(ctx).Infof("CreateStudent: %v", s)
	student, err := uc.repo.Save(ctx, s)
	if err != nil {
		return nil, err
	}
	student.FormatTimeFields()
	return student, nil
}

// GetStudent gets a Student by ID.
func (uc *StudentUsecase) GetStudent(ctx context.Context, id int32) (*Student, error) {
	uc.log.WithContext(ctx).Infof("GetStudent: %v", id)
	student, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	student.FormatTimeFields()
	return student, nil
}

// UpdateStudent updates a Student, and returns the updated Student.
func (uc *StudentUsecase) UpdateStudent(ctx context.Context, s *Student) (*Student, error) {
	uc.log.WithContext(ctx).Infof("UpdateStudent: %v", s)
	student, err := uc.repo.Update(ctx, s)
	if err != nil {
		return nil, err
	}
	student.FormatTimeFields()
	return student, nil
}

// DeleteStudent deletes a Student.
func (uc *StudentUsecase) DeleteStudent(ctx context.Context, id int32) error {
	uc.log.WithContext(ctx).Infof("DeleteStudent: %v", id)
	return uc.repo.Delete(ctx, id)
}

// ListStudent lists all Students.
func (uc *StudentUsecase) ListStudent(ctx context.Context) ([]*Student, error) {
	uc.log.WithContext(ctx).Infof("ListStudent")
	students, err := uc.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	FormatTimeFieldsBatch(students)
	return students, nil
}
