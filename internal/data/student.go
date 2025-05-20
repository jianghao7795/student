package data

import (
	"context"
	"fmt"

	"student/internal/biz"

	errors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
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
		Name:      stu.Name,
		Status:    stu.Status,
		Info:      stu.Info,
		ID:        stu.ID,
		Age:       stu.Age,
		CreatedAt: stu.CreatedAt,
		UpdatedAt: stu.UpdatedAt,
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

func (r *studentRepo) UpdateStudent(ctx context.Context, id int32, s *biz.StudentForm) (*biz.UpdateStudentMessage, error) {
	// TODO: implement the logic of updating student
	var stu biz.Student
	err := r.data.gormDB.First(&stu, id).Error
	if err != nil {
		return nil, err
	}
	stu.Name = s.Name
	stu.Info = s.Info
	stu.Status = s.Status
	stu.Age = s.Age
	err = r.data.gormDB.Save(&stu).Error
	r.log.WithContext(ctx).Info("gormDB: UpdateStudent, id: %d, student: %v", id, stu)
	return &biz.UpdateStudentMessage{
		Message: "Update student success",
	}, err
}

func (r *studentRepo) DeleteStudent(ctx context.Context, id int32) (*biz.DeleteStudentMessage, error) {
	// TODO: implement the logic of deleting student
	var stu biz.Student
	if _, err := r.GetStudent(ctx, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(404, "NOT_FOUND", fmt.Sprintf("student not found with id: %d", id)).WithMetadata(map[string]string{
				"data": "未找到",
			})
		}
		return nil, err
	}
	err := r.data.gormDB.Delete(&stu, id).Error
	r.log.WithContext(ctx).Info("gormDB: DeleteStudent, id: %d", id)
	return &biz.DeleteStudentMessage{
		Message: "Delete student success",
	}, err
}

func (r *studentRepo) ListStudents(ctx context.Context, page int32, pageSize int32, name string) ([]*biz.Student, int32, error) {
	var stus []*biz.Student
	var total int64
	var err error
	if name == "" {
		err = r.data.gormDB.Model(&biz.Student{}).Count(&total).Error
	} else {
		err = r.data.gormDB.Model(&biz.Student{}).Where("name LIKE ?", "%"+name+"%").Count(&total).Error
	}
	if err != nil {
		return nil, 0, err
	}
	// logger.Println(page, pageSize)
	if name == "" {
		err = r.data.gormDB.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Order("id desc").Find(&stus).Error
	} else {
		err = r.data.gormDB.Where("name LIKE ?", "%"+name+"%").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Order("id desc").Find(&stus).Error
	}

	return stus, int32(total), err
}
