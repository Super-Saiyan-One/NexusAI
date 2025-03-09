package controller

import "github.com/gin-gonic/gin"

type RelayImageController interface {
	RelayImage(c *gin.Context)
}

type relayImageController struct {
}

func NewRelayImageController() RelayImageController {
	return &relayImageController{}
}

func (r *relayImageController) RelayImage(c *gin.Context) {
	// TODO
}
