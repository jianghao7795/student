package data

import (
	"context"
	"student/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type studentRepo struct {
	data *Data
	log  *log.Helper
}

func NewStudentRepo(data *Data, logger log.Logger) biz.StudentRepo {
	return &studentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *studentRepo) GetStudent(ctx context.Context, id int32) (*biz.Student, error) {
	// TODO: implement the logic of getting student by id
	var stu biz.Student
	err := r.data.gormDB.First(&stu, id).Error
	r.log.WithContext(ctx).Info("gormDB: GetStudent, id: %d, result: %v", id, stu)
	return &biz.Student{
		Name:   stu.Name,
		Status: stu.Status,
		Info:   stu.Info,
		MODEL: biz.MODEL{
			ID:        stu.ID,
			UpdatedAt: stu.UpdatedAt,
			CreatedAt: stu.CreatedAt,
		},
	}, err
}

func (r *studentRepo) CreateStudent(ctx context.Context, s *biz.StudentForm) (*biz.CreateStudentMessage, error) {
	// TODO: implement the logic of creating student
	var stu biz.Student
	stu.Name = s.Name
	stu.Info = s.Info
	stu.Status = s.Status
	stu.Age = s.Age
	err := r.data.gormDB.Create(&stu).Error
	r.log.WithContext(ctx).Info("gormDB: CreateStudent, student: %v", stu)
	return &biz.CreateStudentMessage{
		Message: "Create student success",
	}, err
}
