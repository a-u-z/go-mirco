package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (s *Server) WriteLog(c *gin.Context) {
	// read json into var
	var requestPayload JSONPayload
	if err := c.ShouldBind(&requestPayload); err != nil {
		// 錯誤處理
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// insert data
	event := LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := s.Models.LogEntry.Insert(event)
	if err != nil {
		// 錯誤處理
		c.JSON(http.StatusBadRequest, err)
		return
	}
	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	c.JSON(http.StatusAccepted, resp)
}
