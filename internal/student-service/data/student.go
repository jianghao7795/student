package data

import (
	"context"
	"time"

	"student/internal/student-service/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type studentRepo struct {
	data *Data
	log  *log.Helper
}

// NewStudentRepo .
func NewStudentRepo(data *Data, logger log.Logger) biz.StudentRepo {
	return &studentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *studentRepo) Save(ctx context.Context, g *biz.Student) (*biz.Student, error) {
	now := time.Now()
	g.CreatedAt = &now
	g.UpdatedAt = &now

	err := r.data.gormDB.WithContext(ctx).Create(g).Error
	if err != nil {
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: Save, student: %v", g)
	return g, nil
}

func (r *studentRepo) Update(ctx context.Context, g *biz.Student) (*biz.Student, error) {
	now := time.Now()
	g.UpdatedAt = &now

	err := r.data.gormDB.WithContext(ctx).Save(g).Error
	if err != nil {
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: Update, student: %v", g)
	return g, nil
}

func (r *studentRepo) FindByID(ctx context.Context, id int32) (*biz.Student, error) {
	var student biz.Student
	err := r.data.gormDB.WithContext(ctx).First(&student, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: FindByID, id: %d, result: %v", id, student)
	return &student, nil
}

func (r *studentRepo) ListAll(ctx context.Context) ([]*biz.Student, error) {
	var students []*biz.Student
	err := r.data.gormDB.WithContext(ctx).Find(&students).Error
	if err != nil {
		return nil, err
	}
	r.log.WithContext(ctx).Infof("gormDB: ListAll, count: %d", len(students))
	return students, nil
}

func (r *studentRepo) Delete(ctx context.Context, id int32) error {
	err := r.data.gormDB.WithContext(ctx).Delete(&biz.Student{}, id).Error
	if err != nil {
		return err
	}
	r.log.WithContext(ctx).Infof("gormDB: Delete, id: %d", id)
	return nil
}
