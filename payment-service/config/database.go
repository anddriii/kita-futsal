package config

import (
	"fmt"
	"log"
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
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Database.MaxLifeTimeConn) * time.Second)
	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConn)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.Database.MaxIdleTime) * time.Second)

	// // ==================================FOR DB MYSQL===================================
	// //for local
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Database.UserName, encodedPassword, config.Database.Host, config.Database.Port, config.Database.Name)

	// // Open connection
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 	Logger: logger.Default.LogMode(logger.Info), // Enable SQL logging
	// })

	// if err != nil {
	// 	log.Fatalf("failed to connect database: %v", err)
	// }

	// sqlDB, err := db.DB()
	// if err != nil {
	// 	log.Fatalf("failed to get database instance: %v", err)
	// }
	// sqlDB.SetMaxOpenConns(config.Database.MaxOpenConn)
	// sqlDB.SetConnMaxLifetime(time.Duration(config.Database.MaxLifeTimeConn) * time.Second)
	// sqlDB.SetMaxIdleConns(config.Database.MaxIdleConn)

	return db, nil

}
