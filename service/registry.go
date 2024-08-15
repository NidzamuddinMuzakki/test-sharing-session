package service

import (
	"github.com/NidzamuddinMuzakki/test-sharing-vision/service/health"
)

// @Notice: Register your services here

type IRegistry interface {
	GetHealth() health.IHealth
	GetPostsService() IPostsService
}

type Registry struct {
	health       health.IHealth
	postsService IPostsService
}

func NewRegistry(health health.IHealth, postsService IPostsService) *Registry {
	return &Registry{
		health:       health,
		postsService: postsService,
	}
}

func (r *Registry) GetHealth() health.IHealth {
	return r.health
}

func (r *Registry) GetPostsService() IPostsService {
	return r.postsService
}
