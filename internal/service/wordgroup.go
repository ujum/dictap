package service

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/ujum/dictap/internal/domain"
	"github.com/ujum/dictap/internal/dto"
	"github.com/ujum/dictap/internal/repo"
	"github.com/ujum/dictap/pkg/logger"
)

type WordGroupService interface {
	Create(ctx context.Context, wordGroupDTO *dto.WordGroupCreate) (string, error)
	GetDefault(ctx context.Context, langID string, userUID string) (*domain.WordGroup, error)
	GetByIDAndUser(ctx context.Context, groupID string, userUID string) (*domain.WordGroup, error)
	GetAllByLangAndUser(ctx context.Context, langID string, userUID string) ([]*domain.WordGroup, error)
}

func newWordGroupService(repos *repo.Repositories, appLogger logger.Logger) *WordGroupServiceImpl {
	return &WordGroupServiceImpl{
		wordGroupRepo: repos.WordGroupRepo,
		logger:        appLogger,
	}
}

type WordGroupServiceImpl struct {
	logger        logger.Logger
	wordGroupRepo repo.WordGroupRepo
}

func (wgs *WordGroupServiceImpl) GetByIDAndUser(ctx context.Context, groupID string, userUID string) (*domain.WordGroup, error) {
	return wgs.wordGroupRepo.FindByIDAndUser(ctx, groupID, userUID)
}

func (wgs *WordGroupServiceImpl) GetDefault(ctx context.Context, langISO string, userUID string) (*domain.WordGroup, error) {
	return wgs.wordGroupRepo.FindByLangAndUser(ctx, langISO, userUID, true)
}

func (wgs *WordGroupServiceImpl) GetAllByLangAndUser(ctx context.Context, langID string, userUID string) ([]*domain.WordGroup, error) {
	return wgs.wordGroupRepo.FindAllByLangAndUser(ctx, langID, userUID)
}

func (wgs *WordGroupServiceImpl) Create(ctx context.Context, wordGroupDTO *dto.WordGroupCreate) (string, error) {
	wordGroup := &domain.WordGroup{}
	if err := copier.Copy(wordGroup, wordGroupDTO); err != nil {
		return "", err
	}
	return wgs.wordGroupRepo.Create(ctx, wordGroup)
}
