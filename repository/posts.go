package repository

import (
	"context"

	commonDs "github.com/NidzamuddinMuzakki/test-sharing-session/go-lib-common/data_source"
	common "github.com/NidzamuddinMuzakki/test-sharing-session/go-lib-common/registry"
	"github.com/NidzamuddinMuzakki/test-sharing-session/model"
	"github.com/jmoiron/sqlx"
)

type IPostsRepository interface {
	CreatePosts(ctx context.Context, payload model.PostsModel) error
}

type posts struct {
	common common.IRegistry
	master *sqlx.DB
	slave  *sqlx.DB
}

func NewPostsRepository(common common.IRegistry, master *sqlx.DB, slave *sqlx.DB) IPostsRepository {
	return &posts{
		common: common,
		master: master,
		slave:  slave,
	}
}

func (r posts) CreatePosts(ctx context.Context, data model.PostsModel) error {
	// const logCtx = "repository.MasterTnC.Create"

	err := commonDs.Exec(ctx, r.master, commonDs.NewStatement(
		nil,
		"insert into posts (title,content,category,status, created_date, updated_date) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		data.Title,
		data.Content,

		data.Category,
		data.Status,
		data.CreatedDate,
		data.UpdatedDate,
	))
	if err != nil {
		return err
	}

	return nil
}
