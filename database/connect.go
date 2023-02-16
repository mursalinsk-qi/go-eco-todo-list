package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var database *sql.DB
var err error

func ConnectDatabase() {
	fmt.Println("connecting")
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"postgres", 5432, "user", "mypassword", "user")
	database, err = sql.Open("postgres", postgresInfo)
	if err != nil {
		log.Fatal(err)
	}
	if err = database.Ping(); err != nil {
		panic(err)
	} else {
		_,err:=database.Exec("CREATE  TABLE IF NOT EXISTS todos(id SERIAL NOT NULL PRIMARY KEY,title VARCHAR(50) NOT NULL,iscomplete BOOLEAN DEFAULT FALSE)")
		if err!=nil{
			fmt.Println(err)
		}
		fmt.Println("Database Connected")
	}

}

func DB() *sql.DB {
	return database
}
