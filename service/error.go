package service

import "github.com/pkg/errors"

var (
	ErrNicknameEmpty = errors.New("empty nickname")
	ErrNameEmpty     = errors.New("empty name")
)
