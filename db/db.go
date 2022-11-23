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

	err = createTicksTable()

	if err != nil {
		fmt.Println(fmt.Errorf("Creating table `%s` error: %s", "ticks", err.Error()))
	}

}

func GetDB() *sql.DB {
	return db
}

func createTicksTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS ticks (
		    id      SERIAL PRIMARY KEY,
		    symbol  VARCHAR(25) NOT NULL,
		    data    JSON NOT NULL,
		    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
`
	_, err := db.Exec(query)
	return err
}
