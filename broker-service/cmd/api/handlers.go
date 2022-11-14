package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Broker(c *gin.Context) {

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the brokerr",
	}
	c.JSON(http.StatusCreated, payload)
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) HandleSubmission(c *gin.Context) {

	var requestPayload RequestPayload
	if err := c.ShouldBind(&requestPayload); err != nil {
		// 錯誤處理
		return
	}

	// 對於帶進去的參數做處理
	switch requestPayload.Action {
	case "auth":
		s.authenticate(c, requestPayload.Auth)
	default:
		c.JSON(http.StatusBadRequest, errors.New("unknown action"))
	}
}

func (s *Server) authenticate(c *gin.Context, a AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	log.Printf("here is :%+v", string(jsonData))

	// call the service
	request, err := http.NewRequest(http.MethodPost, "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	// 為什麼這個 url 可以打到東西，我在本機瀏覽器會導到奇怪的外網？
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("here is err:%+v", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	defer response.Body.Close()
	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		c.JSON(http.StatusBadRequest, errors.New("invalid credentials"))

		return
	} else if response.StatusCode != http.StatusAccepted {
		c.JSON(http.StatusBadRequest, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	if jsonFromService.Error {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	c.JSON(http.StatusAccepted, payload)

}
