package main

import (
	"authentication/data"
	"database/sql"

	corss "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	DB     *sql.DB
	Models data.Models
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
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	s.router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	s.router.POST("/authenticate", s.Authenticate)
}
