package util

import "github.com/stretchr/testify/mock"

type MockPasswordHasher struct {
	mock.Mock
}

func (h *MockPasswordHasher) HashPassword(arg1 string) (string, error) {
	args := h.Called(arg1)
	return args.String(0), args.Error(1)
}

func (h *MockPasswordHasher) VerifyPassword(arg1 string, arg2 string) bool {
	args := h.Called(arg1, arg2)
	return args.Bool(0)
}
