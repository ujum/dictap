package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/ujum/dictap/internal/domain"
	"github.com/ujum/dictap/internal/dto"
	derr "github.com/ujum/dictap/internal/error"
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
	FlagUserAsOAuth(ctx context.Context, user *domain.User) error
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
		return nil, derr.ErrUserIncorrectPass
	}
	return user, nil
}

func (us *UserServiceImpl) Create(ctx context.Context, userDTO *dto.UserCreate) (string, error) {
	existingUser, _ := us.GetByEmail(ctx, userDTO.Email)
	if existingUser != nil {
		return "", derr.ErrAlreadyExists
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
	us.setActiveLangBinding(user)
	_, err := us.userRepo.Create(ctx, user)
	return user.UID, err
}

func (us *UserServiceImpl) setActiveLangBinding(user *domain.User) {
	if len(user.LangBinding) > 0 {
		var hasActive bool
		for i := 0; i < len(user.LangBinding); i++ {
			if hasActive {
				user.LangBinding[i].Active = false
			} else {
				hasActive = user.LangBinding[i].Active
			}
		}
		if !hasActive {
			user.LangBinding[0].Active = true
		}
	}
}

func (us *UserServiceImpl) Update(ctx context.Context, userDTO *dto.UserUpdate) error {
	user := &domain.User{}
	if err := copier.Copy(user, userDTO); err != nil {
		return err
	}
	return us.userRepo.Update(ctx, user)
}

func (us *UserServiceImpl) DeleteByUid(ctx context.Context, uid string) error {
	err := us.userRepo.DeleteByUID(ctx, uid)
	if err != derr.ErrNotFound {
		return err
	}
	return nil
}

func (us *UserServiceImpl) ChangePassword(ctx context.Context, uid string, credentials *dto.ChangeUserPassword) error {
	user, err := us.GetByUID(ctx, uid)
	if err != nil {
		return err
	}
	if user.Password != "" && !us.passHashService.CheckPasswordHash(credentials.OldPassword, user.Password) {
		return derr.ErrUserIncorrectPass
	}
	hashedPass, err := us.passHashService.HashPassword(credentials.Password)
	user.Password = hashedPass
	if err != nil {
		return err
	}
	return us.userRepo.Update(ctx, user)
}

func (us *UserServiceImpl) FlagUserAsOAuth(ctx context.Context, user *domain.User) error {
	user.OAuthUser = true
	return us.userRepo.Update(ctx, user)
}
