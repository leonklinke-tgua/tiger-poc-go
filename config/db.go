package config

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewDBConnection() (*sqlx.DB, error) {
	// add configs
	config := DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "postgres",
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Printf("Error connecting to database: %v\n", err)
		// return nil, err
		// todo config the db connection
	}

	// if err = db.Ping(); err != nil {
	// 	log.Printf("Error pinging database: %v\n", err)
	// 	return nil, err
	// }

	return db, nil
}
