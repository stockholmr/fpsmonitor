package main

import (
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

type Computer struct {
	gorm.Model

	Name string `gorm:"type:text not null;"`
}

type NetworkAdapter struct {
	gorm.Model

	ComputerID int    `gorm:"type:int not null;"`
	Name       string `gorm:"type:text not null;"`
	MacAddress string
	IPAddress  string
}

type User struct {
	gorm.Model

	ComputerID int    `gorm:"type:int not null;"`
	Username   string `gorm:"type:text not null;"`
}
