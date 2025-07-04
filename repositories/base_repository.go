package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type BaseRepository struct {
	DB *gorm.DB
}

func LoadDBConfig() string {
	// Load environment variables from .env file (optional)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Get environment variables
	server := os.Getenv("DB_SERVER")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")

	// Check if all required variables are set
	if server == "" || port == "" || user == "" || password == "" || database == "" {
		return "" // Return empty string if variables are missing
	}

	// Build and return the connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		server, user, password, port, database)
	return connString
}

func (dal *BaseRepository) CreateConnection() (*gorm.DB, error) {
	// If the DB is already initialized, return it
	if dal.DB != nil {
		return dal.DB, nil
	}

	// Load the database connection string
	connString := LoadDBConfig()
	if connString == "" {
		return nil, fmt.Errorf("failed to load database configuration: environment variables missing")
	}

	// Connect to the database
	sqlDB, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	// Ensure we can connect by pinging the database
	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging the database: %w", err)
	}

	// Initialize GORM with the *sql.DB connection
	db, err := gorm.Open(sqlserver.New(sqlserver.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error initializing GORM: %w", err)
	}

	dal.DB = db
	return dal.DB, nil
}

func (dal *BaseRepository) CloseConnection() {
	if dal.DB != nil {
		sqlDB, _ := dal.DB.DB()
		sqlDB.Close()
	}
}
