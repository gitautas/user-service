package models

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Status struct {
	Code    int32
	Message string
}

func NewStatus(code int32, message string) *Status {
	return &Status{
		Code:    code,
		Message: message,
	}
}

func (s *Status) RPC() error {
	switch s.Code {
	case http.StatusBadRequest:
		s.Code = int32(codes.InvalidArgument)
	case http.StatusInternalServerError:
		s.Code = int32(codes.Internal)
	case http.StatusNotFound:
		s.Code = int32(codes.NotFound)
	}

	return status.Error(codes.Code(s.Code), s.Message)
}
