package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

const ( //для отладки
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "qwerty"
	dbname   = "db_valutevalues"
)

// Создание соединения с БД
func Pgdbconnect() (*sql.DB, error) {
	var pgdb *sql.DB
	var err error
	pgdbconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname) //для отладки
	//pgdbconn := "postgresql://postgres:qwerty@clair_postgres:5432?sslmode=disable"
	pgdb, err = sql.Open("postgres", pgdbconn)
	if err != nil {
		return pgdb, err
	}
	pgdb.SetConnMaxIdleTime(1 * time.Second)
	pgdb.SetConnMaxLifetime(1 * time.Second)
	fmt.Println(pgdb)
	return pgdb, nil
}
