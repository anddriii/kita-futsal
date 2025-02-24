package config

import (
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	config := Config
	encodedPassword := url.QueryEscape(config.Database.Password)
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.Database.UserName, encodedPassword, config.Database.Host, config.Database.Port, config.Database.Name)

	fmt.Printf("Database Config: %+v\n", config.Database)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Database.MaxLifeTimeConn) * time.Second)
	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConn)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.Database.MaxIdleTime) * time.Second)

	return db, nil

}
