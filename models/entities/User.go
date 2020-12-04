package entities

import "DigitalRegionAPI/models/auxiliary"

type User struct {
	auxiliary.BaseModel
	Username     string
	Password     string
	School       string
	DataUploads  uint   `gorm:"default: 0; not null;"`
}
