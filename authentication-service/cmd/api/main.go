package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	bindAddr = "80"
	maxTries = 10
)

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")
	db := connectToDB()
	if db == nil {
		log.Panic("Can't connect to Postgres!")
	}

	app := &Config{
		DB:     db,
		Models: data.New(db),
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", bindAddr),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func dbCon(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
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
	dsn := os.Getenv("DSN")

	count := 0
	for {
		db, err := dbCon(dsn)

		if err != nil {
			log.Println("Postgres not yet ready...")
			count++
		} else {
			log.Println("Connected to Postgres!")
			return db
		}

		if count > maxTries {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}

}
