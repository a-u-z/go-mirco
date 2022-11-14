package main

import (
	"time"

	corss "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

// 之後要要反轉注入什麼一次性的 singleton
func NewServer() *Server {
	r := gin.New()
	r.Use(gin.Recovery())
	s := &Server{
		router: r,
	}
	s.routesV2()
	return s
}

// 回傳的 Type 是 *gin.Engine
func (s *Server) routesV2() {

	s.router.Use(corss.New(corss.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "authorization", "Referer"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour, // pre-flight request cache
	}))
	s.router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	s.router.POST("/", s.Broker)
	s.router.POST("/handle", s.HandleSubmission)
}
