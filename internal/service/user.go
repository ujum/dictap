package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/ujum/dictap/internal/domain"
	"github.com/ujum/dictap/internal/dto"
	"github.com/ujum/dictap/internal/repo"
	"github.com/ujum/dictap/pkg/logger"
	"time"
)

type UserService interface {
	GetByUID(ctx context.Context, uid string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *dto.UserCreate) (string, error)
	Update(ctx context.Context, user *dto.UserUpdate) error
	GetAll(ctx context.Context) ([]*domain.User, error)
	DeleteByUid(ctx context.Context, uid string) error
	GetByCredentials(ctx context.Context, credentials *dto.UserCredentials) (*domain.User, error)
	ChangePassword(ctx context.Context, uid string, credentials *dto.ChangeUserPassword) error
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

func (us *UserServiceImpl) GetByUID(ctx context.Context, uid string) (*domain.User, error) {
	return us.userRepo.FindByUID(ctx, uid)
}

func (us *UserServiceImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return us.userRepo.FindByEmail(ctx, email)
}

func (us *UserServiceImpl) GetAll(ctx context.Context) ([]*domain.User, error) {
	return us.userRepo.FindAll(ctx)
}

func (us *UserServiceImpl) GetByCredentials(ctx context.Context, credentials *dto.UserCredentials) (*domain.User, error) {
	user, err := us.GetByEmail(ctx, credentials.Email)
	if err != nil {
		return nil, err
	}

	if !us.passHashService.CheckPasswordHash(credentials.Password, user.Password) {
		return nil, domain.ErrUserIncorrectPass
	}
	return user, nil
}

func (us *UserServiceImpl) Create(ctx context.Context, userDTO *dto.UserCreate) (string, error) {
	existingUser, _ := us.GetByEmail(ctx, userDTO.Email)
	if existingUser != nil {
		return "", domain.ErrUserAlreadyExists
	}
	user := &domain.User{}
	if err := copier.CopyWithOption(user, userDTO, copier.Option{}); err != nil {
		return "", err
	}
	if userDTO.Password != "" {
		hashedPass, err := us.passHashService.HashPassword(userDTO.Password)
		user.Password = hashedPass
		if err != nil {
			return "", err
		}
	}
	user.RegisteredAt = time.Now()
	user.UID = uuid.New().String()
	_, err := us.userRepo.Create(ctx, user)
	return user.UID, err
}

func (us *UserServiceImpl) Update(ctx context.Context, userDTO *dto.UserUpdate) error {
	user := &domain.User{}
	if err := copier.Copy(user, userDTO); err != nil {
		return err
	}
	return us.userRepo.Update(ctx, user)
}

func (us *UserServiceImpl) DeleteByUid(ctx context.Context, uid string) error {
	return us.userRepo.DeleteByUID(ctx, uid)
}

func (us *UserServiceImpl) ChangePassword(ctx context.Context, uid string, credentials *dto.ChangeUserPassword) error {
	user, err := us.GetByUID(ctx, uid)
	if err != nil {
		return err
	}
	if user.Password != "" && !us.passHashService.CheckPasswordHash(credentials.OldPassword, user.Password) {
		return domain.ErrUserIncorrectPass
	}
	hashedPass, err := us.passHashService.HashPassword(credentials.Password)
	user.Password = hashedPass
	if err != nil {
		return err
	}
	return us.userRepo.Update(ctx, user)
}
