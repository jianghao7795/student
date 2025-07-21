package biz

import (
	"student/internal/pkg/jwt"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewStudentUsecase,
	NewUserUsecase,
	NewRBACUsecase,
	NewErrorUsecase,
	jwt.NewJWTUtil,
)
