package handler

import (
	"net/http"

	"github.com/NidzamuddinMuzakki/test-sharing-session/common/response"
	"github.com/NidzamuddinMuzakki/test-sharing-session/go-lib-common/logger"
	common "github.com/NidzamuddinMuzakki/test-sharing-session/go-lib-common/registry"
	commonModel "github.com/NidzamuddinMuzakki/test-sharing-session/go-lib-common/response/model"
	"github.com/NidzamuddinMuzakki/test-sharing-session/model"
	"github.com/NidzamuddinMuzakki/test-sharing-session/service"
	"github.com/gin-gonic/gin"
)

type IPosts interface {
	CreatePosts(c *gin.Context)
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
		c.JSON(response.BadRequest(ctx).WithMessage(errValidate.Error()).ToHTTPCodeAndMap())
		return
	}

	err = h.serviceRegistry.GetPostsService().CreatePosts(ctx, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "Posts Created",
		Data:         payload,
	})
}
