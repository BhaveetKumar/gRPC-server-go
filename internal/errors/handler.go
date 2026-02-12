package errors

import (
	"github.com/BhaveetKumar/gRPC-server-go/internal/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToStatus(err error, log *logger.Logger) error {
	if err == nil {
		return nil
	}

	switch err {
	case ErrPostNotFound:
		log.Error("post not found")
		return status.Error(codes.NotFound, err.Error())
	case ErrInvalidInput:
		log.Error("invalid input")
		return status.Error(codes.InvalidArgument, err.Error())
	case ErrDuplicatePost:
		log.Error("duplicate post")
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		log.Error("internal error")
		return status.Error(codes.Internal, ErrInternal.Error())
	}
}
