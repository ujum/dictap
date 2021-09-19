package service

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/ujum/dictap/internal/domain"
	"github.com/ujum/dictap/internal/dto"
	"github.com/ujum/dictap/internal/repo"
	"github.com/ujum/dictap/pkg/logger"
)

type WordService interface {
	Create(ctx context.Context, wordDTO *dto.WordCreate) (string, error)
	GetByGroup(ctx context.Context, groupID string) ([]*domain.Word, error)
	GetByName(ctx context.Context, name string) (*domain.Word, error)
	AddToGroup(ctx context.Context, name string, groupID string) error
	RemoveFromGroup(ctx context.Context, name string, groupID string) error
	MoveToGroup(ctx context.Context, name string, fromGroupID, toGroupID string) error
}

func newWordService(repos *repo.Repositories, appLogger logger.Logger) *WordServiceImpl {
	return &WordServiceImpl{
		wordRepo: repos.WordRepo,
		logger:   appLogger,
	}
}

type WordServiceImpl struct {
	logger   logger.Logger
	wordRepo repo.WordRepo
}

func (w *WordServiceImpl) Create(ctx context.Context, wordDTO *dto.WordCreate) (string, error) {
	var wordId string
	word, err := w.GetByName(ctx, wordDTO.Name)
	if err != nil && word == nil {

		word = &domain.Word{}
		if err := copier.Copy(word, wordDTO); err != nil {
			return "", err
		}
		wordId, err = w.wordRepo.Create(ctx, word)
		if err != nil {
			return "", err
		}
	} else {
		wordId = word.ID
	}
	return wordId, err
}

func (w *WordServiceImpl) GetByGroup(ctx context.Context, groupID string) ([]*domain.Word, error) {
	return w.wordRepo.FindByGroup(ctx, groupID)
}

func (w *WordServiceImpl) GetByName(ctx context.Context, name string) (*domain.Word, error) {
	return w.wordRepo.FindByName(ctx, name)
}

func (w *WordServiceImpl) AddToGroup(ctx context.Context, name string, groupID string) error {
	word, err := w.wordRepo.FindByNameAndGroup(ctx, name, groupID)
	if err == nil && word != nil {
		return nil
	}
	return w.wordRepo.AddToGroup(ctx, name, groupID)
}

func (w *WordServiceImpl) RemoveFromGroup(ctx context.Context, name string, groupID string) error {
	return w.wordRepo.RemoveFromGroup(ctx, name, groupID)
}

func (w *WordServiceImpl) MoveToGroup(ctx context.Context, name string, fromGroupID, toGroupID string) error {
	err := w.AddToGroup(ctx, name, toGroupID)
	if err != nil {
		return err
	}
	err = w.RemoveFromGroup(ctx, name, fromGroupID)
	if err != nil {
		return err
	}
	return nil
}
