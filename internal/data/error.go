package data

import (
	"context"

	"student/internal/biz"

	errors "student/internal/data/errors"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type errorRepo struct {
	data *Data
	log  *log.Helper
}

func NewErrorRepo(data *Data, logger log.Logger) biz.ErrorRepo {
	return &errorRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// 实现 从 gormDB 中获取错误信息
func (r *errorRepo) GetErrorInfo(ctx context.Context, errorCode int32, errorType string) (*biz.ErrorInfo, error) {
	// 首先尝试从预定义错误中获取
	if predefinedError, exists := biz.GetPredefinedError(errorCode); exists {
		// 如果指定了错误类型，检查是否匹配
		if errorType != "" && predefinedError.ErrorType != errorType {
			return nil, errors.Error404()
		}
		return predefinedError, nil
	}

	// 如果预定义错误中没有，尝试从数据库获取
	var errorInfo biz.ErrorInfo
	query := r.data.gormDB.WithContext(ctx).Model(&biz.ErrorInfo{}).Where("error_code = ?", errorCode)
	if errorType != "" {
		query = query.Where("error_type = ?", errorType)
	}

	err := query.First(&errorInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Error404()
		}
		return nil, errors.Error400(err)
	}

	r.log.WithContext(ctx).Info("gormDB: GetErrorInfo, errorCode: %d, errorType: %s, result: %v", errorCode, errorType, errorInfo)

	// 格式化时间字段
	errorInfo.FormatTimeFields()

	return &errorInfo, nil
}

// 实现 从 gormDB 中获取错误码列表
func (r *errorRepo) ListErrorCodes(ctx context.Context, errorType string, page, pageSize int32) ([]*biz.ErrorInfo, int32, error) {
	var errorInfos []*biz.ErrorInfo
	var total int64
	var err error

	// 构建查询条件
	query := r.data.gormDB.WithContext(ctx).Model(&biz.ErrorInfo{})
	if errorType != "" {
		query = query.Where("error_type = ?", errorType)
	}

	// 获取总数
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, errors.Error400(err)
	}

	// 获取分页数据
	query = r.data.gormDB.WithContext(ctx).Model(&biz.ErrorInfo{})
	if errorType != "" {
		query = query.Where("error_type = ?", errorType)
	}

	err = query.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Order("error_code asc").Find(&errorInfos).Error
	if err != nil {
		return nil, 0, errors.Error400(err)
	}

	// 为每个错误记录格式化时间字段
	for _, errorInfo := range errorInfos {
		errorInfo.FormatTimeFields()
	}

	return errorInfos, int32(total), err
}

// 实现 从 gormDB 中创建自定义错误
func (r *errorRepo) CreateCustomError(ctx context.Context, errorForm *biz.ErrorForm) (*biz.CreateErrorMessage, error) {
	var errorInfo biz.ErrorInfo
	errorInfo.ErrorCode = errorForm.ErrorCode
	errorInfo.ErrorType = errorForm.ErrorType
	errorInfo.ErrorMessage = errorForm.ErrorMessage
	errorInfo.ErrorDescription = errorForm.ErrorDescription
	errorInfo.Solution = errorForm.Solution

	err := r.data.gormDB.WithContext(ctx).Create(&errorInfo).Error
	if err != nil {
		return nil, errors.Error400(err)
	}

	r.log.WithContext(ctx).Info("gormDB: CreateCustomError, errorInfo: %v", errorInfo)
	return &biz.CreateErrorMessage{
		Message: "Create custom error success",
	}, err
}
