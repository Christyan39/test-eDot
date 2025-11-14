package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Christyan39/test-eDot/pkg/config"
	_ "github.com/go-sql-driver/mysql"
)

// MySQL creates a new MySQL database connection
func MySQL(dsn string) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// InitMySQL initializes database connection with configuration from environment
func InitMySQL(serviceName string) (*sql.DB, error) {
	dbUser := config.GetEnv("DB_USER", "root")
	dbPassword := config.GetEnv("DB_PASSWORD", "")
	dbHost := config.GetEnv("DB_HOST", "localhost")
	dbPort := config.GetEnv("DB_PORT", "3306")
	dbName := config.GetEnv("DB_NAME", "ecommerce")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	log.Printf("[STARTUP] Connecting to MySQL: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
