package main

const webPort = ":80" // 要記得 ":" ，很重要阿，沒有的話就不會通

type Config struct{}

func main() {
	s := NewServer()
	s.router.Run(webPort)
}
