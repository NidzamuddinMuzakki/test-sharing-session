package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	constants "github.com/NidzamuddinMuzakki/test-sharing-vision/common/constant"
	"github.com/NidzamuddinMuzakki/test-sharing-vision/common/util"
	"github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/logger"
	"github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/registry"
	"github.com/NidzamuddinMuzakki/test-sharing-vision/model"
	"github.com/NidzamuddinMuzakki/test-sharing-vision/repository"
	"github.com/jmoiron/sqlx"
)

type IPostsService interface {
	CreatePosts(ctx context.Context, payload model.RequestPostModel) (id int, err error)
	GetPostsList(ctx context.Context, payload model.RequestGetListPostModel) ([]model.ResponsePostModel, uint64, error)
	GetPostsDetail(ctx context.Context, id int) (*model.ResponsePostModel, error)
	GetListLogPosts(ctx context.Context, id int) ([]model.ResponseLogPostsModel, error)
	DeletePosts(ctx context.Context, id int) error
	UpdatePosts(ctx context.Context, payload model.RequestUpdatePostModel) error
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
func (s postsService) GetListLogPosts(ctx context.Context, id int) ([]model.ResponseLogPostsModel, error) {

	ccc, err := s.repoRegistry.GetPostsRepository().GetListLogPosts(ctx, id)
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, err
	}
	responses := []model.ResponseLogPostsModel{}

	for _, dataR := range ccc {
		response := model.ResponseLogPostsModel{
			Id:        dataR.Id,
			ArticleId: dataR.ArticleId,

			CategoryStatus: dataR.CategoryStatus,
			CreatedDate:    dataR.CreatedDate,
			UpdatedDate:    dataR.UpdatedDate,
		}
		if dataR.DataBefore != "" {
			var dataString interface{}
			err := json.Unmarshal([]byte(dataR.DataBefore), &dataString)
			if err != nil {
				logger.Error(ctx, err.Error(), err)

			}
			response.DataBefore = dataString
		}

		if dataR.DataAfter != "" {
			var dataString interface{}
			err := json.Unmarshal([]byte(dataR.DataAfter), &dataString)
			if err != nil {
				logger.Error(ctx, err.Error(), err)

			}
			response.DataAfter = dataString
		}
		responses = append(responses, response)

	}

	return responses, nil
}
func (s postsService) GetPostsDetail(ctx context.Context, id int) (*model.ResponsePostModel, error) {

	ccc, err := s.repoRegistry.GetPostsRepository().GetDetailPosts(ctx, id)
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, err
	}
	var dataString string
	json.Unmarshal([]byte(ccc.Content), &dataString)
	response := model.ResponsePostModel{
		Id:       ccc.Id,
		Title:    ccc.Title,
		Content:  dataString,
		Category: ccc.Category,
		Status:   ccc.Status,
	}
	return &response, nil
}
func (s postsService) GetPostsList(ctx context.Context, payload model.RequestGetListPostModel) ([]model.ResponsePostModel, uint64, error) {

	responses := []model.ResponsePostModel{}

	if payload.Search != "" {
		payload.Search = strings.ToLower(payload.Search)
	}

	resp, totalCounts, err := s.repoRegistry.GetPostsRepository().GetListPosts(ctx, payload)

	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, 0, err
	}

	for _, dataR := range resp {
		var dataString string
		json.Unmarshal([]byte(dataR.Content), &dataString)
		response := model.ResponsePostModel{
			Id:       dataR.Id,
			Title:    dataR.Title,
			Content:  dataString,
			Category: dataR.Category,
			Status:   dataR.Status,
		}
		responses = append(responses, response)
	}

	return responses, totalCounts, nil

}
func (s postsService) UpdatePosts(ctx context.Context, payload model.RequestUpdatePostModel) error {
	now := time.Now()
	dataGet, err := s.repoRegistry.GetPostsRepository().GetDetailPosts(ctx, payload.Id)
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return err
	}
	dataContent, _ := json.Marshal(payload.Content)

	data := model.PostsModel{
		Id:          dataGet.Id,
		Title:       payload.Title,
		Content:     string(dataContent),
		Category:    payload.Category,
		Status:      payload.Status,
		CreatedDate: &now,
		UpdatedDate: &now,
	}
	DataAfter, errr := json.Marshal(data)
	DataBefore, errr := json.Marshal(dataGet)
	fmt.Println(errr, "ni")
	dataLog := model.LogPostsModel{

		DataBefore:     string(DataBefore),
		DataAfter:      string(DataAfter),
		CategoryStatus: "updated-data-" + payload.Status,
		CreatedDate:    &now,
		UpdatedDate:    &now,
	}
	doFunc := util.TxFunc(func(tx *sqlx.Tx) (id int, err error) {
		err = s.repoRegistry.GetPostsRepository().UpdatePosts(ctx, tx, data)
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			// if strings.Contains(err.Error(), "Duplicate entry") {
			// 	err = errors.New("Duplicate entry")
			// } else if strings.Contains(err.Error(), "foreign key constraint fails") {
			// 	err = errors.New("App ID not Found in Master App")
			// }

			return 0, err
		}

		dataLog.ArticleId = dataGet.Id

		err = s.repoRegistry.GetPostsRepository().CreateLogPosts(ctx, tx, dataLog)
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return 0, err
		}
		return 0, nil
	})

	_, err = s.repoRegistry.GetUtilTx().WithTx(ctx, doFunc, nil)
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return err
	}

	return nil
}
func (s postsService) DeletePosts(ctx context.Context, id int) error {
	now := time.Now()
	dataGet, err := s.repoRegistry.GetPostsRepository().GetDetailPosts(ctx, id)
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return err
	}

	data := model.PostsModel{
		Id:          dataGet.Id,
		Title:       dataGet.Title,
		Content:     dataGet.Content,
		Category:    dataGet.Category,
		Status:      constants.StatusThrash,
		CreatedDate: &now,
		UpdatedDate: &now,
	}
	DataAfter, _ := json.Marshal(data)
	DataBefore, _ := json.Marshal(dataGet)
	dataLog := model.LogPostsModel{

		DataBefore:     string(DataBefore),
		DataAfter:      string(DataAfter),
		CategoryStatus: "updated-status-" + constants.StatusThrash,
		CreatedDate:    &now,
		UpdatedDate:    &now,
	}
	doFunc := util.TxFunc(func(tx *sqlx.Tx) (id int, err error) {
		err = s.repoRegistry.GetPostsRepository().UpdatePosts(ctx, tx, data)
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			// if strings.Contains(err.Error(), "Duplicate entry") {
			// 	err = errors.New("Duplicate entry")
			// } else if strings.Contains(err.Error(), "foreign key constraint fails") {
			// 	err = errors.New("App ID not Found in Master App")
			// }

			return 0, err
		}

		dataLog.ArticleId = dataGet.Id

		err = s.repoRegistry.GetPostsRepository().CreateLogPosts(ctx, tx, dataLog)
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return 0, err
		}
		return 0, nil
	})

	_, err = s.repoRegistry.GetUtilTx().WithTx(ctx, doFunc, nil)
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return err
	}

	return nil
}
func (s postsService) CreatePosts(ctx context.Context, payload model.RequestPostModel) (id int, err error) {
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
	DataAfter, _ := json.Marshal(data)
	dataLog := model.LogPostsModel{

		DataBefore:     "",
		DataAfter:      string(DataAfter),
		CategoryStatus: "created-" + payload.Status,
		CreatedDate:    &now,
		UpdatedDate:    &now,
	}
	doFunc := util.TxFunc(func(tx *sqlx.Tx) (id int, err error) {
		ids, err := s.repoRegistry.GetPostsRepository().CreatePosts(ctx, tx, data)
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			// if strings.Contains(err.Error(), "Duplicate entry") {
			// 	err = errors.New("Duplicate entry")
			// } else if strings.Contains(err.Error(), "foreign key constraint fails") {
			// 	err = errors.New("App ID not Found in Master App")
			// }

			return 0, err
		}

		dataLog.ArticleId = int(ids)

		err = s.repoRegistry.GetPostsRepository().CreateLogPosts(ctx, tx, dataLog)
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return 0, err
		}
		return int(ids), nil
	})

	idsss, err := s.repoRegistry.GetUtilTx().WithTx(ctx, doFunc, nil)
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return 0, err
	}

	return idsss, nil
}
