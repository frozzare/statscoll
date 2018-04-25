package db

import (
	"github.com/jinzhu/gorm"
)

// DB represents the database.
type DB struct {
	*gorm.DB
}

// Open opens a new database connection.
func Open(typ string, dsn string) (*DB, error) {
	conn, err := gorm.Open(typ, dsn)
	if err != nil {
		return nil, err
	}

	db := &DB{conn}

	return db, nil
}
