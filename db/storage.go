package db

import (
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

func (s *Storage) InitDB(dsn string) (*Repository, error) {
	s.log.WithField("dsn", dsn).Debugln("opening")
	database, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}
	return &Repository{Db: database}, nil
}
