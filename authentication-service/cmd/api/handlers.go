package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) Authenticate(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("here is err:%+v", err)
		// Handle error
	}

	var requestPayload RequestPayload

	json.Unmarshal(jsonData, &requestPayload)

	// 不知道為何這個會失效
	// if err := c.ShouldBind(&requestPayload); err != nil {
	// 	// 錯誤處理
	// 	log.Printf("here is err:%+v", err)
	// 	return
	// }

	// validate the user against the database
	user, err := s.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("invalid credentials"))
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		c.JSON(http.StatusBadRequest, errors.New("invalid credentials"))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	c.JSON(http.StatusAccepted, payload)
}
