package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	corss "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	router   *gin.Engine
	Rabbitmq *amqp.Connection
}

// 之後要要反轉注入什麼一次性的 singleton
func NewServer() *Server {
	r := gin.New()
	rabbitmq, err := connect()
	if err != nil {
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		log.Printf("can't connect to rabbitmq")
		log.Println("err:", err)
		os.Exit(1)
	}

	r.Use(gin.Recovery())
	s := &Server{
		router:   r,
		Rabbitmq: rabbitmq,
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

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
