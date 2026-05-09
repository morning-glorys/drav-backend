package apperror

import "errors"

var (
	ErrNotFound       = errors.New("data tidak ditemukan")
	ErrInvalidInput   = errors.New("input data tidak valid")
	ErrConflict       = errors.New("data sudah ada atau konflik")
	ErrInternalServer = errors.New("terjadi kesalahan internal server")
)
