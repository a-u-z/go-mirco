package main

import (
	"fmt"
	"log"
	"log-service/logs"
	"net"
	"net/http"
	"net/rpc"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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

func (s *Server) listenFromRpc() error {
	log.Println("Starting RPC server on port ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func (s *Server) listenFromGRPC() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	ss := grpc.NewServer()

	logs.RegisterLogServiceServer(ss, &LogServer{Models: s.Models})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := ss.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
