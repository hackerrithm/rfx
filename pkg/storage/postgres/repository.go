package asjson

import (
	"database/sql"
	"fmt"
)

// Storage stores user data in JSON files
type Storage struct {
	db     *sql.DB
	config Config
}

// Config ...
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// NewStorage returns a new JSON  storage
func NewStorage(cfg Config) (*Storage, error) {

	s := new(Storage)

	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port,
		cfg.User, cfg.Password, cfg.Database)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	s.db = db
	s.config = cfg

	fmt.Println("Successfully connected!")

	return s, nil
}
