package database

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)
var database *sql.DB
func ConnectDatabase(){
	db, err := sql.Open("postgres", "user=postgres password=ms310 dbname=Todolist sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	} else {
		fmt.Println("Database Connected")
		database=db
	}

}

func DB() *sql.DB {
	return database
}
