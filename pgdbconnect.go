package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

//const (
//	host     = "localhost"
//	port     = 5432
//	user     = "postgres"
//	password = "qwerty"
//	dbname   = "db_valutevalues"
//)

// Создание соединения с БД
func pgdbconnect() (*sql.DB, error) {
	var pgdb *sql.DB
	//pgdbconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	pgdbconn := "postgresql://postgres:qwerty@clair_postgres:5432?sslmode=disable"
	pgdb, err = sql.Open("postgres", pgdbconn)
	if err != nil {
		return pgdb, err
	}
	pgdb.SetConnMaxIdleTime(1 * time.Second)
	pgdb.SetConnMaxLifetime(1 * time.Second)
	return pgdb, nil
}
