package database

import (
	"fmt"
	"log"
	"micro-site/api/model"
	"micro-site/config"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBClient represents our database connection pool
var DB *gorm.DB

// InitDB initializes the database connection pool
func InitDB(cfg *config.Config) {
	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	fmt.Println("Attempting to connect to database:", cfg.DBName)

	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	//DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	//migration

	err = migration(db)

	if err != nil {
		log.Fatal("Migration error:", err)
	}

	DB = db
	// Ping the database to verify the connection

	fmt.Println("Successfully connected to MySQL database!")
}

// CloseDB closes the database connection pool
func CloseDB() {
	sqlDB, _ := DB.DB() // Get *sql.DB from GORM

	if sqlDB != nil {
		err := sqlDB.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			fmt.Println("Database connection closed.")
		}
	}
}

func migration(db *gorm.DB) error {
	db.AutoMigrate(model.User{})
	db.AutoMigrate(model.Role{})
	db.AutoMigrate(model.RoleUser{})

	return nil
}
