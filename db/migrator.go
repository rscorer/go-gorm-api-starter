package db

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"io/ioutil"
)

func (s *Storage) migrations(db *sql.DB) error {
	// are there any files to migrate?
	dir := "./db/migrations"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+dir,
		"mysql", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		return err
	}
	return nil
}
