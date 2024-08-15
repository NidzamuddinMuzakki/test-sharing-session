package repository

import (
	"context"
	"errors"
	"fmt"

	commonDs "github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/data_source"
	"github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/logger"
	common "github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/registry"
	"github.com/NidzamuddinMuzakki/test-sharing-vision/model"
	"github.com/jmoiron/sqlx"
)

type IPostsRepository interface {
	GetListPosts(ctx context.Context, payload model.RequestGetListPostModel) ([]model.PostsModel, uint64, error)
	GetListLogPosts(ctx context.Context, id int) ([]model.LogPostsModel, error)
	GetDetailPosts(ctx context.Context, id int) (model.PostsModel, error)
	CreatePosts(ctx context.Context, tx *sqlx.Tx, payload model.PostsModel) (id int64, err error)
	UpdatePosts(ctx context.Context, tx *sqlx.Tx, payload model.PostsModel) error
	CreateLogPosts(ctx context.Context, tx *sqlx.Tx, data model.LogPostsModel) error
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
func (r posts) GetDetailPosts(ctx context.Context, id int) (model.PostsModel, error) {

	var (
		data model.PostsModel
	)

	selectQuery := `select 
    	id,
		title, 
		content, 
		category, 
		status,
		created_date,
		updated_date
	from posts
	where id = ? `

	err := commonDs.Exec(ctx, r.master, commonDs.NewStatement(&data, selectQuery, id))
	if err != nil {
		return data, err
	}

	return data, nil

}

func (r posts) GetListLogPosts(ctx context.Context, id int) ([]model.LogPostsModel, error) {

	var (
		list []model.LogPostsModel
		args []any
	)

	selectQuery := "SELECT " +
		"id, " +
		"article_id, " +
		"data_before, " +
		"data_after, " +
		"category_status, " +
		"created_date, " +
		"updated_date " +
		"FROM log_posts " +
		"where 1=1 and  article_id=? order by id desc "
	args = append(args, id)
	err := commonDs.Exec(ctx, r.slave, commonDs.NewStatement(&list, selectQuery, args...))
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, err
	}

	return list, nil

}

func (r posts) GetListPosts(ctx context.Context, payload model.RequestGetListPostModel) ([]model.PostsModel, uint64, error) {

	var (
		list              []model.PostsModel
		totalTransactions uint64
		filters           []string
		args              []any
	)

	countQuery := "select count(id) from posts where 1=1 "
	selectQuery := "SELECT " +
		"id, " +
		"title, " +
		"content, " +
		"category, " +
		"status, " +
		"created_date, " +
		"updated_date " +
		"FROM posts " +
		"where 1=1 "

	filters = append(filters, "and status = ? ")
	args = append(args, payload.Status)

	if payload.Search != "" {
		filters = append(filters, "and  LOWER(title) like ?  ")
		args = append(args, "%"+payload.Search+"%")
	}

	for _, f := range filters {
		countQuery = fmt.Sprintf("%s %s", countQuery, f)
		selectQuery = fmt.Sprintf("%s %s", selectQuery, f)
	}

	offset := (payload.Limit * payload.Offset) - payload.Limit
	selectQuery = fmt.Sprintf("%s LIMIT %d OFFSET %d", selectQuery, payload.Limit, offset)

	err := commonDs.Exec(ctx, r.slave, commonDs.NewStatement(&list, selectQuery, args...))
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, 0, err
	}

	err = commonDs.Exec(ctx, r.slave, commonDs.NewStatement(&totalTransactions, countQuery, args...))
	if err != nil {

		logger.Error(ctx, err.Error(), err)
		return nil, 0, err
	}

	return list, totalTransactions, nil

}

func (r posts) UpdatePosts(ctx context.Context, tx *sqlx.Tx, data model.PostsModel) error {

	insertQuery := "update posts set title=?,content=?,category=?,status=?, updated_date=? where id=?"
	stmtx, err := tx.PreparexContext(ctx, insertQuery)
	if err != nil {
		return err
	}
	res, err := stmtx.ExecContext(
		ctx,
		data.Title,
		data.Content,
		data.Category,
		data.Status,
		data.UpdatedDate,
		data.Id,
	)
	intss, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if intss == 0 {
		err = errors.New("nothing update")
		return err
	}

	return nil
}
func (r posts) CreatePosts(ctx context.Context, tx *sqlx.Tx, data model.PostsModel) (id int64, err error) {

	insertQuery := "insert into posts (title,content,category,status, created_date, updated_date) VALUES(?,?,?,?,?,?)"
	stmtx, err := tx.PreparexContext(ctx, insertQuery)
	if err != nil {
		return 0, err
	}
	res, err := stmtx.ExecContext(
		ctx,
		data.Title,
		data.Content,
		data.Category,
		data.Status,
		data.CreatedDate,
		data.UpdatedDate,
	)
	intss, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if intss == 0 {
		err = errors.New("nothing update")
		return 0, err
	}

	return intss, nil
}

func (r posts) CreateLogPosts(ctx context.Context, tx *sqlx.Tx, data model.LogPostsModel) error {
	fmt.Println(data.DataAfter, data.DataBefore)
	insertQuery := "insert into log_posts (article_id,data_before,data_after,category_status, created_date, updated_date) VALUES(?,?,?,?,?,?)"
	stmtx, err := tx.PreparexContext(ctx, insertQuery)
	if err != nil {
		return err
	}
	res, err := stmtx.ExecContext(
		ctx,
		data.ArticleId,
		data.DataBefore,
		data.DataAfter,
		data.CategoryStatus,
		data.CreatedDate,
		data.UpdatedDate,
	)
	intss, err := res.RowsAffected()
	fmt.Println(err, intss)
	if err != nil {
		return err
	}
	if intss == 0 {
		err = errors.New("nothing update")
		return err
	}

	return nil
}
