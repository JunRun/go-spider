package modle

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type sqlDB struct {
}

const (
	host     = "127.0.0.1"
	user     = "postgres"
	dbname   = "anime"
	password = "123456"
)

var (
	GDB    *gorm.DB
	config = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, dbname, password)
)

func GetConnect() *gorm.DB {
	if GDB == nil {
		DB, err := gorm.Open("postgres", config)
		if err != nil {
			panic(err)
		}
		GDB = DB
	}
	return GDB
}
