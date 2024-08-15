package handler

import (
	"github.com/NidzamuddinMuzakki/test-sharing-session/handler/health"
)

// @Notice: Register your http deliveries here

type IRegistry interface {
	GetHealth() health.IHealth
	GetPosts() IPosts
}

type Registry struct {
	health health.IHealth
	posts  IPosts
}

func NewRegistry(health health.IHealth, posts IPosts) *Registry {
	return &Registry{
		health: health,
		posts:  posts,
	}
}

func (r *Registry) GetHealth() health.IHealth {
	return r.health
}

func (r *Registry) GetPosts() IPosts {
	return r.posts
}
