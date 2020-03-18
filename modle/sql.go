package modle

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"sync"
)

type sqlDB struct {
}

const (
	host     = "127.0.0.1"
	user     = "mysql"
	dbname   = "shows"
	password = "123456"
	port     = 3306
)

var (
	err    error
	gDB    *gorm.DB
	config = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbname)
	single = sync.Once{}
)

func GetConnect() *gorm.DB {
	single.Do(func() {
		if gDB, err = gorm.Open("mysql", config); err != nil {
			panic(err)
		}
	})
	return gDB
}
