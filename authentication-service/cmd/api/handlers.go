package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Authenticate(c *gin.Context) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 要放一個 swagger 的 required
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, requestPayload)
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
