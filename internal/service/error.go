package service

import (
	"context"

	pb "student/api/errors/v1"
	"student/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type ErrorService struct {
	pb.UnimplementedErrorServiceServer

	error *biz.ErrorUsecase
	log   *log.Helper
}

func NewErrorService(error *biz.ErrorUsecase, logger log.Logger) *ErrorService {
	return &ErrorService{
		error: error,
		log:   log.NewHelper(logger),
	}
}

func (s *ErrorService) GetErrorInfo(ctx context.Context, req *pb.GetErrorInfoRequest) (*pb.GetErrorInfoReply, error) {
	s.log.Info("get error info", req.ErrorCode, req.ErrorType)

	errorInfo, err := s.error.GetErrorInfo(ctx, req.ErrorCode, req.ErrorType)
	if err != nil {
		return nil, err
	}

	reply := &pb.GetErrorInfoReply{
		ErrorCode:        errorInfo.ErrorCode,
		ErrorType:        errorInfo.ErrorType,
		ErrorMessage:     errorInfo.ErrorMessage,
		ErrorDescription: errorInfo.ErrorDescription,
		Solution:         errorInfo.Solution,
		CreatedAt:        errorInfo.CreatedAtStr,
		UpdatedAt:        errorInfo.UpdatedAtStr,
	}

	return reply, nil
}

func (s *ErrorService) ListErrorCodes(ctx context.Context, req *pb.ListErrorCodesRequest) (*pb.ListErrorCodesReply, error) {
	s.log.Info("list error codes", req.ErrorType, req.Page, req.PageSize)

	var page, pageSize int32
	var err error

	if req.PageSize != 0 {
		pageSize = req.PageSize
	} else {
		pageSize = 10
	}

	if req.Page != 0 {
		page = req.Page
	} else {
		page = 1
	}

	errorInfos, total, err := s.error.ListErrorCodes(ctx, req.ErrorType, page, pageSize)
	if err != nil {
		return nil, err
	}

	var errors []*pb.ErrorInfo
	for _, errorInfo := range errorInfos {
		errors = append(errors, &pb.ErrorInfo{
			ErrorCode:        errorInfo.ErrorCode,
			ErrorType:        errorInfo.ErrorType,
			ErrorMessage:     errorInfo.ErrorMessage,
			ErrorDescription: errorInfo.ErrorDescription,
			Solution:         errorInfo.Solution,
			CreatedAt:        errorInfo.CreatedAtStr,
			UpdatedAt:        errorInfo.UpdatedAtStr,
		})
	}

	return &pb.ListErrorCodesReply{
		Errors:   errors,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *ErrorService) CreateCustomError(ctx context.Context, req *pb.CreateCustomErrorRequest) (*pb.CreateCustomErrorReply, error) {
	s.log.Info("create custom error", req.ErrorCode, req.ErrorType, req.ErrorMessage)

	errorForm := &biz.ErrorForm{
		ErrorCode:        req.ErrorCode,
		ErrorType:        req.ErrorType,
		ErrorMessage:     req.ErrorMessage,
		ErrorDescription: req.ErrorDescription,
		Solution:         req.Solution,
	}

	result, err := s.error.CreateCustomError(ctx, errorForm)
	if err != nil {
		return nil, err
	}

	reply := &pb.CreateCustomErrorReply{
		Message: result.Message,
		ErrorInfo: &pb.ErrorInfo{
			ErrorCode:        req.ErrorCode,
			ErrorType:        req.ErrorType,
			ErrorMessage:     req.ErrorMessage,
			ErrorDescription: req.ErrorDescription,
			Solution:         req.Solution,
		},
	}

	return reply, nil
}
