package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Broker(c *gin.Context) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	c.JSON(http.StatusCreated, payload)
}
