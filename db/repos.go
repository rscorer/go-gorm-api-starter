package db

import (
	"github.com/jinzhu/gorm"
)

type Repository struct {
	Db *gorm.DB
}
