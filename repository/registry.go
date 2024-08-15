package repository

import "github.com/NidzamuddinMuzakki/test-sharing-vision/common/util"

// @Notice: Register your repositories here

type IRegistry interface {
	GetPostsRepository() IPostsRepository
	GetUtilTx() *util.TransactionRunner
}

type Registry struct {
	postsRepository IPostsRepository
	masterUtilTx    *util.TransactionRunner
}

func NewRegistryRepository(
	masterUtilTx *util.TransactionRunner,
	postsRepository IPostsRepository,
) *Registry {
	return &Registry{
		masterUtilTx:    masterUtilTx,
		postsRepository: postsRepository,
	}
}

func (r Registry) GetUtilTx() *util.TransactionRunner {
	return r.masterUtilTx
}

func (r Registry) GetPostsRepository() IPostsRepository {
	return r.postsRepository
}
