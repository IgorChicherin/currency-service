package db

import (
	"database/sql"
	"fmt"
	"github.com/IgorChicherin/currency-service/config"
	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() {
	c := config.GetConfig()
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.GetString("database.host"),
		c.GetInt("database.port"),
		c.GetString("database.username"),
		c.GetString("database.password"),
		c.GetString("database.dbName"))

	var err error

	db, err = sql.Open("postgres", psqlconn)

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

}

func GetDB() *sql.DB {
	return db
}

func createTable() error {
	return nil
}
