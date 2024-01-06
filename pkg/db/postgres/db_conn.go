package postgres

import (
	"fmt"

	"github.com/pcaokhai/scraper/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func NewPsqlDBConnection(config *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%v password=%s dbname=%s port=%s sslmode=%s",
					config.Postgres.Host,
					config.Postgres.User,
					config.Postgres.Password,
					config.Postgres.Dbname,
					config.Postgres.Port,
					config.Postgres.Sslmode)

	pdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Connect to postgres failed %v\n", err)
		return nil, err
	}

	fmt.Println("Successfully connected to DB")
	db = pdb
	return db, nil
}

func GetDb() *gorm.DB {
	return db
}