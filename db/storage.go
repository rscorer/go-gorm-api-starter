package db

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type Storage struct {
	log *logrus.Logger
}

func CreateStorage(log *logrus.Logger) *Storage {
	return &Storage{
		log: log,
	}
}

func (s *Storage) OpenDB(dsn string) (*Repository, error) {
	s.log.WithField("dsn", dsn).Debug("opening")
	database, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &Repository{Db: database}, nil
}

func (s *Storage) InitDB(dsn string) (*Repository, error) {
	db, err := s.OpenDB(dsn)
	if err != nil {
		return nil, err
	}
	s.log.Debug("migrating")
	err = s.migrations(db.Db.DB())
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		s.log.WithError(err).Errorln("not errnochange")
		return nil, err
	}
	return db, nil
}

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
