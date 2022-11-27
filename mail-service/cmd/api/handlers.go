package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) SendMail(c *gin.Context) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage
	if err := c.ShouldBind(&requestPayload); err != nil {
		// 錯誤處理
		c.JSON(http.StatusBadRequest, err)
		return
	}
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err := s.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	c.JSON(http.StatusAccepted, payload)
}
