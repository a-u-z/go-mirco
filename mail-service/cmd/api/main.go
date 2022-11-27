package main

import (
	"log"
	"os"
	"strconv"
)

const webPort = ":80" // 要記得 ":" ，很重要阿，沒有的話就不會通

type Config struct{}

func main() {
	s := NewServer()
	s.Mailer = createMail()
	s.router.Run(webPort)
	log.Println("Starting mail service on port", webPort)

}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	return m
}
