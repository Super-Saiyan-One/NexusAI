package controller // TODO

import (
	"net/http"
	"nexus-ai/model"
	"nexus-ai/repository"
	"nexus-ai/service"

	"github.com/gin-gonic/gin"
)

type RelayVideoController interface {
	GetTaskVideoRepo() repository.TaskVideoRepository
	Text2Video(c *gin.Context)
	Image2Video(c *gin.Context)
	VideoExtend(c *gin.Context)
	LipSync(c *gin.Context)
}

type relayVideoController struct {
	relayVideoService service.RelayVideoService
}

func NewRelayVideoController() RelayVideoController {
	relayVideoService := service.NewRelayVideoService()
	return &relayVideoController{relayVideoService: relayVideoService}
}

func (rvc *relayVideoController) GetTaskVideoRepo() repository.TaskVideoRepository {
	return repository.NewTaskVideoRepository(model.GetDB())
}

func (rvc *relayVideoController) Text2Video(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodPost:
		rvc.relayVideoService.Text2Video(rvc.GetTaskVideoRepo(), nil)
	case http.MethodGet:
		rvc.relayVideoService.Text2Video(rvc.GetTaskVideoRepo(), nil)
	}
}

func (rvc *relayVideoController) Image2Video(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodPost:
		rvc.relayVideoService.Image2Video(rvc.GetTaskVideoRepo(), nil)
	case http.MethodGet:
		rvc.relayVideoService.Image2Video(rvc.GetTaskVideoRepo(), nil)
	}
}

func (rvc *relayVideoController) VideoExtend(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodPost:
		rvc.relayVideoService.VideoExtend(rvc.GetTaskVideoRepo(), nil)
	case http.MethodGet:
		rvc.relayVideoService.VideoExtend(rvc.GetTaskVideoRepo(), nil)
	}
}

func (rvc *relayVideoController) LipSync(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodPost:
		rvc.relayVideoService.LipSync(rvc.GetTaskVideoRepo(), nil)
	case http.MethodGet:
		rvc.relayVideoService.LipSync(rvc.GetTaskVideoRepo(), nil)
	}
}
