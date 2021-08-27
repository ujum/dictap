package service

import (
	"github.com/ujum/dictap/internal/repo"
	"github.com/ujum/dictap/pkg/logger"
)

type UserService interface {
	GetByUid(uid string) string
}

func (us *UserServiceImpl) GetByUid(uid string) string {
	return us.userRepo.GetByUid(uid)
}

type UserServiceImpl struct {
	userRepo repo.UserRepo
	logger   logger.Logger
}
