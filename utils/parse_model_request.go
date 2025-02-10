package utils

import (
	"nexus-ai/dto"
	"strings"

	"github.com/gin-gonic/gin"
)

func ParseModelRequest(c *gin.Context) (*dto.ModelRequest, error) {
	var modelRequest dto.ModelRequest
	err := UnmarshalRequestBody(c, &modelRequest)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(c.Request.URL.Path, "/embeddings") {
		if modelRequest.Model == "" {
			modelRequest.Model = c.Param("model")
		}
	}

	return &modelRequest, nil
}
