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
	Create(ctx context.Context, user *dto.UserCreate) (string, error)
	Update(ctx context.Context, user *dto.UserUpdate) error
	GetAll(ctx context.Context) ([]*domain.User, error)
	DeleteByUid(ctx context.Context, uid string) error
	GetByCredentials(ctx context.Context, credentials *dto.UserCredentials) (*domain.User, error)
}

type UserServiceImpl struct {
	userRepo        repo.UserRepo
	logger          logger.Logger
	passHashService PassHashService
}

func newUserService(repos *repo.Repositories, appLogger logger.Logger) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:        repos.UserRepo,
		logger:          appLogger,
		passHashService: &BCryptPassHashService{},
	}
}

func (us *UserServiceImpl) GetByUid(ctx context.Context, uid string) (*domain.User, error) {
	return us.userRepo.FindByUid(ctx, uid)
}

func (us *UserServiceImpl) GetAll(ctx context.Context) ([]*domain.User, error) {
	return us.userRepo.FindAll(ctx)
}

func (us *UserServiceImpl) GetByCredentials(ctx context.Context, credentials *dto.UserCredentials) (*domain.User, error) {
	user, err := us.GetByUid(ctx, credentials.Username)
	if err != nil {
		return nil, err
	}

	if !us.passHashService.CheckPasswordHash(credentials.Password, user.Password) {
		return nil, domain.ErrUserIncorrectPass
	}
	return user, nil
}

func (us *UserServiceImpl) Create(ctx context.Context, userDTO *dto.UserCreate) (string, error) {
	existingUser, _ := us.GetByUid(ctx, userDTO.Uid)
	if existingUser != nil {
		return "", domain.ErrUserAlreadyExists
	}
	password, err := us.passHashService.HashPassword(userDTO.Password)
	if err != nil {
		return "", err
	}

	user := &domain.User{}
	if err := copier.CopyWithOption(user, userDTO, copier.Option{}); err != nil {
		return "", err
	}
	user.Password = password
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
	return us.userRepo.DeleteByUid(ctx, uid)
}
