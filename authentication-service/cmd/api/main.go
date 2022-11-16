package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

const webPort = ":80"

var counts int64

func main() {
	log.Println("Starting authentication service")

	// connect to DB 在 container 裡面用這個
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// set up config
	s := NewServer()
	s.DB = conn
	s.Models = New(db)
	s.router.Run(webPort)

}

func openDB(dsn string) (*sql.DB, error) {
	var err error
	db, err = sql.Open("pgx", dsn) // 這邊要用 = 而不是 := 不然會 panic ，不過不知為何，readme 第一篇
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	// dsn := os.Getenv("DSN")

	dsn := "host=postgres port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5"
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			log.Printf("here is :%+v", connection.Stats())
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
