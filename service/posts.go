package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/NidzamuddinMuzakki/test-sharing-session/go-lib-common/logger"
	"github.com/NidzamuddinMuzakki/test-sharing-session/go-lib-common/registry"
	"github.com/NidzamuddinMuzakki/test-sharing-session/model"
	"github.com/NidzamuddinMuzakki/test-sharing-session/repository"
)

type IPostsService interface {
	CreatePosts(ctx context.Context, payload model.RequestPostModel) error
}

type postsService struct {
	common       registry.IRegistry
	repoRegistry repository.IRegistry
}

func NewPostsService(common registry.IRegistry, repoRegistry repository.IRegistry) IPostsService {
	return &postsService{
		common:       common,
		repoRegistry: repoRegistry,
	}
}

func (s postsService) CreatePosts(ctx context.Context, payload model.RequestPostModel) error {
	now := time.Now()
	dataContent, _ := json.Marshal(payload.Content)

	data := model.PostsModel{
		Title:       payload.Title,
		Content:     string(dataContent),
		Category:    payload.Category,
		Status:      payload.Status,
		CreatedDate: &now,
		UpdatedDate: &now,
	}

	err := s.repoRegistry.GetPostsRepository().CreatePosts(ctx, data)
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		// if strings.Contains(err.Error(), "Duplicate entry") {
		// 	err = errors.New("Duplicate entry")
		// } else if strings.Contains(err.Error(), "foreign key constraint fails") {
		// 	err = errors.New("App ID not Found in Master App")
		// }

		return err
	}

	return nil
}
