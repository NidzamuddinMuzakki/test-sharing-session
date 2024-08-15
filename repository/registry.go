package repository

// @Notice: Register your repositories here

type IRegistry interface {
	GetPostsRepository() IPostsRepository
}

type Registry struct {
	postsRepository IPostsRepository
}

func NewRegistryRepository(
	postsRepository IPostsRepository,
) *Registry {
	return &Registry{
		postsRepository: postsRepository,
	}
}

func (r Registry) GetPostsRepository() IPostsRepository {
	return r.postsRepository
}
