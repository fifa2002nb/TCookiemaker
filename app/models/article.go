package models

import (
	"github.com/jinzhu/gorm"
)

type Article struct {
	gorm.Model
	Category string
	HashID   string
	Url      string
	Title    string
	Content  string `sql:"size:10000"`
	Enabled  bool
}
