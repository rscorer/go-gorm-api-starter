package db

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type Storage struct {
	log *log.Logger
}

func CreateStorage(log *log.Logger) *Storage {
	return &Storage{
		log: log,
	}
}

func (storage *Storage) OpenDB(dsn string) (*Repository, error) {
	log.WithField("dsn", dsn).Debug("opening")
	database, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &Repository{Db: database}, nil
}

func (storage *Storage) InitDB(dsn string) (*Repository, error) {
	db, err := storage.OpenDB(dsn)
	if err != nil {
		return nil, err
	}
	storage.log.Debug("migrating")
	err = storage.migrations(db.Db.DB())
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}
	return db, nil
}

func (storage *Storage) migrations(db *sql.DB) error {
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
	return m.Up()
}
