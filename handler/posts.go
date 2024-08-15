package handler

import (
	"net/http"

	"github.com/NidzamuddinMuzakki/test-sharing-vision/common/response"
	"github.com/NidzamuddinMuzakki/test-sharing-vision/common/util"
	"github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/logger"
	common "github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/registry"
	commonModel "github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/response/model"
	"github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/validator"
	"github.com/NidzamuddinMuzakki/test-sharing-vision/model"
	"github.com/NidzamuddinMuzakki/test-sharing-vision/service"
	"github.com/gin-gonic/gin"
)

type IPosts interface {
	CreatePosts(c *gin.Context)
	GetListPosts(c *gin.Context)
	GetListLogPosts(c *gin.Context)
	GetDetailPosts(c *gin.Context)
	DeletePosts(c *gin.Context)
	UpdatePosts(c *gin.Context)
}
type DataResponse struct {
	Before interface{}
	Now    interface{}
}
type posts struct {
	common          common.IRegistry
	serviceRegistry service.IRegistry
}

func NewPosts(common common.IRegistry, serviceRegistry service.IRegistry) IPosts {
	return &posts{
		common:          common,
		serviceRegistry: serviceRegistry,
	}
}

func (h posts) GetListLogPosts(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		payload model.RequestGetDetailPostModel
	)

	if err := c.ShouldBindUri(&payload); err != nil {
		logger.Error(ctx, err.Error(), err)
		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      err.Error(),
		})
		return
	}

	errValidate := h.common.GetValidator().Struct(payload)
	if errValidate != nil {
		dataError := validator.ToErrResponseV2(errValidate)

		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      "Failed Data",
			Data:         dataError,
		})
		return
	}

	data, err := h.serviceRegistry.GetPostsService().GetListLogPosts(ctx, int(payload.Id))
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.Response{
			Status:  commonModel.StatusError,
			Message: http.StatusText(http.StatusUnprocessableEntity),
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Success: true,
		Data:    data,
	})

}
func (h posts) UpdatePosts(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		payload model.RequestUpdatePostModel
	)

	if err := c.ShouldBindUri(&payload); err != nil {
		logger.Error(ctx, err.Error(), err)
		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      err.Error(),
		})
		return
	}

	err := c.ShouldBind(&payload)
	if err != nil {

		logger.Error(ctx, err.Error(), err)
		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      err.Error(),
		})
		return
	}

	errValidate := h.common.GetValidator().Struct(payload)
	if errValidate != nil {
		dataError := validator.ToErrResponseV2(errValidate)

		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      "Failed Data",
			Data:         dataError,
		})
		return
	}

	err = h.serviceRegistry.GetPostsService().UpdatePosts(ctx, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.Response{
			Status:  commonModel.StatusError,
			Message: http.StatusText(http.StatusUnprocessableEntity),
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "Posts Updated",
		Data:         payload,
	})

}
func (h posts) DeletePosts(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		payload model.RequestGetDetailPostModel
	)

	if err := c.ShouldBindUri(&payload); err != nil {
		logger.Error(ctx, err.Error(), err)
		c.AbortWithStatusJSON(response.BadRequest(ctx).WithMessage(err.Error()).ToHTTPCodeAndMap())
		return
	}

	errValidate := h.common.GetValidator().Struct(payload)
	if errValidate != nil {
		dataError := validator.ToErrResponseV2(errValidate)

		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      "Failed Data",
			Data:         dataError,
		})
		return
	}

	err := h.serviceRegistry.GetPostsService().DeletePosts(ctx, int(payload.Id))
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.Response{
			Status:  commonModel.StatusError,
			Message: http.StatusText(http.StatusUnprocessableEntity),
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "Posts Deleted",
		Data:         payload,
	})

}

func (h posts) GetDetailPosts(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		payload model.RequestGetDetailPostModel
	)

	if err := c.ShouldBindUri(&payload); err != nil {
		logger.Error(ctx, err.Error(), err)
		c.AbortWithStatusJSON(response.BadRequest(ctx).WithMessage(err.Error()).ToHTTPCodeAndMap())
		return
	}

	errValidate := h.common.GetValidator().Struct(payload)
	if errValidate != nil {
		dataError := validator.ToErrResponseV2(errValidate)

		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      "Failed Data",
			Data:         dataError,
		})
		return
	}

	list, err := h.serviceRegistry.GetPostsService().GetPostsDetail(ctx, int(payload.Id))
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.Response{
			Status:  commonModel.StatusError,
			Message: http.StatusText(http.StatusUnprocessableEntity),
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, commonModel.Response{
		Status:  commonModel.StatusSuccess,
		Message: http.StatusText(http.StatusOK),
		Data:    list,
	})

}
func (h posts) GetListPosts(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		payload model.RequestGetListPostModel
	)

	if err := c.ShouldBindUri(&payload); err != nil {
		logger.Error(ctx, err.Error(), err)
		c.AbortWithStatusJSON(response.BadRequest(ctx).WithMessage(err.Error()).ToHTTPCodeAndMap())
		return
	}
	err := c.ShouldBindQuery(&payload)
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		c.AbortWithStatusJSON(response.BadRequest(ctx).WithMessage(err.Error()).ToHTTPCodeAndMap())
		return
	}

	errValidate := h.common.GetValidator().Struct(payload)
	if errValidate != nil {
		dataError := validator.ToErrResponseV2(errValidate)

		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      "Failed Data",
			Data:         dataError,
		})
		return
	}

	if payload.Limit == 0 {
		payload.Limit = 10
	}

	if payload.Offset == 0 {
		payload.Offset = 1
	}

	list, totalCount, err := h.serviceRegistry.GetPostsService().GetPostsList(ctx, payload)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, commonModel.Response{
			Status:  commonModel.StatusError,
			Message: http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}

	totalPages, previousPage, nextPage := util.Pagination(int64(totalCount), int64(payload.Limit), int64(payload.Offset))

	c.JSON(http.StatusOK, commonModel.Response{
		Status:       commonModel.StatusSuccess,
		Message:      http.StatusText(http.StatusOK),
		Data:         list,
		TotalRecords: totalCount,
		CurrentPage:  payload.Offset,
		NextPage:     uint(nextPage),
		PreviousPage: uint(previousPage),
		TotalPages:   uint(totalPages),
	})

}
func (h posts) CreatePosts(c *gin.Context) {
	// const logCtx = "delivery.http.tnc.CreateTnC"
	var (
		ctx = c.Request.Context()
		// span    = h.common.GetSentry().StartSpan(ctx, logCtx)
		payload model.RequestPostModel
	)

	err := c.ShouldBind(&payload)
	if err != nil {

		logger.Error(ctx, err.Error(), err)
		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      err.Error(),
		})
		return
	}
	errValidate := h.common.GetValidator().Struct(payload)
	if errValidate != nil {
		dataError := validator.ToErrResponseV2(errValidate)
		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      "Failed Data",
			Data:         dataError,
		})
		return
	}

	id, err := h.serviceRegistry.GetPostsService().CreatePosts(ctx, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      err.Error(),
		})
		return
	}
	data := model.ResponsePostModel{
		Id:       id,
		Title:    payload.Title,
		Category: payload.Category,
		Content:  payload.Content,
		Status:   payload.Status,
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "Posts Created",
		Data:         data,
	})
}
