package mysql

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbOnce sync.Once
	DB     *gorm.DB
)

func ConnectDatabase() *gorm.DB {
	dbOnce.Do(func() {
		port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
		portTest := os.Getenv("DB_PORT")
		_ = portTest

		portStr := os.Getenv("DB_PORT")
		if portStr == "" {
			panic("DB_PORT is not set")
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			panic("DB_PORT must be a number")
		}

		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			port,
			os.Getenv("DB_DATABASE"),
		)

		database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("Failed to connect database: " + err.Error())
		}

		DB = database
		fmt.Println("MySQL connected")
	})

	return DB
}
