package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/rpc"

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
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type RPCPayload struct {
	Name string
	Data string
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
	case "log":
		// s.logItem(c, requestPayload.Log)
		// s.logEventViaRabbit(c, requestPayload.Log)
		s.logItemViaRPC(c, requestPayload.Log)
	case "mail":
		s.sendMail(c, requestPayload.Mail)
	default:
		c.JSON(http.StatusBadRequest, errors.New("unknown action"))
	}
}

func (s *Server) sendMail(c *gin.Context, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// call the mail service
	mailServiceURL := "http://mailer-service/send"

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// send back json
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	c.JSON(http.StatusAccepted, payload)

}

func (s *Server) logItem(c *gin.Context, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"
	c.JSON(http.StatusAccepted, payload)

}
func (s *Server) authenticate(c *gin.Context, a AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

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

// logEventViaRabbit logs an event using the logger-service. It makes the call by pushing the data to RabbitMQ.
func (s *Server) logEventViaRabbit(c *gin.Context, l LogPayload) {
	err := s.pushToQueue(l.Name, l.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	c.JSON(http.StatusAccepted, payload)
}

// pushToQueue pushes a message into RabbitMQ
func (s *Server) pushToQueue(name, msg string) error {
	emitter, err := NewEventEmitter(s.Rabbitmq)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) logItemViaRPC(c *gin.Context, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	rpcPayload := RPCPayload(l)
	// rpcPayload := RPCPayload{
	// 	Name: l.Name,
	// 	Data: l.Data,
	// }

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	c.JSON(http.StatusAccepted, payload)
}
