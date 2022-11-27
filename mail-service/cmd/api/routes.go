package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	Mailer Mail
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

	s.router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "authorization", "Referer"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour, // pre-flight request cache
	}))
	s.router.GET("/ping", func(c *gin.Context) {
		c.String(200, "mailer-pong")
	})
	s.router.POST("/send", s.SendMail)

}
