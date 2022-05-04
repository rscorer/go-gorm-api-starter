package db

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Repository struct {
	Db *gorm.DB
}

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
	database, err := gorm.Open(mysql.Open(dsn))
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
	mdb, err := db.Db.DB()
	if err != nil {
		return nil, err
	}
	err = s.migrations(mdb)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}
	return db, nil
}
