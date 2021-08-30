package service

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/ujum/dictap/internal/domain"
	"github.com/ujum/dictap/internal/dto"
	"github.com/ujum/dictap/internal/repo"
	"github.com/ujum/dictap/pkg/logger"
	"time"
)

type UserService interface {
	GetByUid(ctx context.Context, uid string) (*domain.User, error)
	Create(ctx context.Context, user *dto.UserCreate) error
	Update(ctx context.Context, user *dto.UserUpdate) error
	GetAll(ctx context.Context) ([]*domain.User, error)
	DeleteByUid(ctx context.Context, uid string) error
}

type UserServiceImpl struct {
	userRepo repo.UserRepo
	logger   logger.Logger
}

func (us *UserServiceImpl) GetByUid(ctx context.Context, uid string) (*domain.User, error) {
	return us.userRepo.FindByUid(ctx, uid)
}

func (us *UserServiceImpl) GetAll(ctx context.Context) ([]*domain.User, error) {
	return us.userRepo.FindAll(ctx)
}

func (us *UserServiceImpl) Create(ctx context.Context, userDTO *dto.UserCreate) error {
	existingUser, _ := us.GetByUid(ctx, userDTO.Uid)
	if existingUser != nil {
		return domain.ErrUserAlreadyExists
	}
	user := &domain.User{}
	if err := copier.CopyWithOption(user, userDTO, copier.Option{}); err != nil {
		return err
	}
	user.RegisteredAt = time.Now()

	return us.userRepo.Create(ctx, user)
}

func (us *UserServiceImpl) Update(ctx context.Context, userDTO *dto.UserUpdate) error {
	user := &domain.User{}
	if err := copier.Copy(user, userDTO); err != nil {
		return err
	}
	return us.userRepo.Update(ctx, user)
}

func (us *UserServiceImpl) DeleteByUid(ctx context.Context, uid string) error {
	//user, _ := us.GetByUid(ctx, uid)
	//if user == nil {
	//	return domain.ErrUserNotFound
	//}
	return us.userRepo.DeleteByUid(ctx, uid)
}
