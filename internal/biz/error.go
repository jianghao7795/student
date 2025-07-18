package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// ErrorInfo 错误信息模型
type ErrorInfo struct {
	ID               uint
	ErrorCode        int32
	ErrorType        string
	ErrorMessage     string
	ErrorDescription string
	Solution         string
	CreatedAt        *time.Time
	UpdatedAt        *time.Time

	// 格式化后的时间字符串
	CreatedAtStr string `gorm:"-" json:"created_at_str,omitempty"`
	UpdatedAtStr string `gorm:"-" json:"updated_at_str,omitempty"`
}

// FormatTimeFields 格式化时间字段
func (e *ErrorInfo) FormatTimeFields() {
	if e.CreatedAt != nil {
		e.CreatedAtStr = e.CreatedAt.Format(TimeFormat)
	}
	if e.UpdatedAt != nil {
		e.UpdatedAtStr = e.UpdatedAt.Format(TimeFormat)
	}
}

// ErrorForm 错误表单
type ErrorForm struct {
	ErrorCode        int32
	ErrorType        string
	ErrorMessage     string
	ErrorDescription string
	Solution         string
}

// CreateErrorMessage 创建错误消息
type CreateErrorMessage struct {
	ID      int32
	Message string
}

// 定义 Error 的操作接口
type ErrorRepo interface {
	GetErrorInfo(ctx context.Context, errorCode int32, errorType string) (*ErrorInfo, error)
	ListErrorCodes(ctx context.Context, errorType string, page, pageSize int32) ([]*ErrorInfo, int32, error)
	CreateCustomError(ctx context.Context, errorForm *ErrorForm) (*CreateErrorMessage, error)
}

type ErrorUsecase struct {
	repo ErrorRepo
	log  *log.Helper
}

// 初始化 ErrorUsecase
func NewErrorUsecase(repo ErrorRepo, logger log.Logger) *ErrorUsecase {
	return &ErrorUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// 获取错误信息
func (uc *ErrorUsecase) GetErrorInfo(ctx context.Context, errorCode int32, errorType string) (*ErrorInfo, error) {
	uc.log.Info("get error info", errorCode, errorType)
	return uc.repo.GetErrorInfo(ctx, errorCode, errorType)
}

// 获取错误码列表
func (uc *ErrorUsecase) ListErrorCodes(ctx context.Context, errorType string, page, pageSize int32) ([]*ErrorInfo, int32, error) {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	return uc.repo.ListErrorCodes(ctx, errorType, page, pageSize)
}

// 创建自定义错误
func (uc *ErrorUsecase) CreateCustomError(ctx context.Context, errorForm *ErrorForm) (*CreateErrorMessage, error) {
	uc.log.Info("create custom error", errorForm)
	return uc.repo.CreateCustomError(ctx, errorForm)
}

// 预定义的错误信息映射
var PredefinedErrors = map[int32]ErrorInfo{
	400: {
		ErrorCode:        400,
		ErrorType:        "CLIENT",
		ErrorMessage:     "Bad Request",
		ErrorDescription: "请求参数错误或格式不正确",
		Solution:         "请检查请求参数是否正确",
	},
	401: {
		ErrorCode:        401,
		ErrorType:        "AUTH",
		ErrorMessage:     "Unauthorized",
		ErrorDescription: "用户未授权或认证失败",
		Solution:         "请先登录或检查认证信息",
	},
	403: {
		ErrorCode:        403,
		ErrorType:        "AUTH",
		ErrorMessage:     "Forbidden",
		ErrorDescription: "用户没有权限访问该资源",
		Solution:         "请联系管理员获取相应权限",
	},
	404: {
		ErrorCode:        404,
		ErrorType:        "CLIENT",
		ErrorMessage:     "Not Found",
		ErrorDescription: "请求的资源不存在",
		Solution:         "请检查请求的URL是否正确",
	},
	500: {
		ErrorCode:        500,
		ErrorType:        "SERVER",
		ErrorMessage:     "Internal Server Error",
		ErrorDescription: "服务器内部错误",
		Solution:         "请稍后重试或联系技术支持",
	},
	1001: {
		ErrorCode:        1001,
		ErrorType:        "STUDENT",
		ErrorMessage:     "Student Not Found",
		ErrorDescription: "学生信息不存在",
		Solution:         "请检查学生ID是否正确",
	},
	1002: {
		ErrorCode:        1002,
		ErrorType:        "USER",
		ErrorMessage:     "User Not Found",
		ErrorDescription: "用户信息不存在",
		Solution:         "请检查用户ID是否正确",
	},
	1003: {
		ErrorCode:        1003,
		ErrorType:        "AUTH",
		ErrorMessage:     "Invalid Credentials",
		ErrorDescription: "用户名或密码错误",
		Solution:         "请检查用户名和密码是否正确",
	},
}

// GetPredefinedError 获取预定义错误信息
func GetPredefinedError(errorCode int32) (*ErrorInfo, bool) {
	if errorInfo, exists := PredefinedErrors[errorCode]; exists {
		return &errorInfo, true
	}
	return nil, false
}
