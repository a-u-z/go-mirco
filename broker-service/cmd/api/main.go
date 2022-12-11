package main

import "log"

const webPort = ":80" // 要記得 ":" ，很重要阿，沒有的話就不會通

type Config struct{}

func main() {
	s := NewServer()
	err := s.router.Run(webPort)
	if err != nil {
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		log.Panic(err)
	}

	defer s.Rabbitmq.Close() // 超級重要，這個要在 main 裡面 defer 不然，跑完那個 func rabbitmq 就被 close 了。 https://stackoverflow.com/questions/36579759/golang-rabbitmq-channel-connection-is-not-open
}
