package repository

import (
	"database/sql"
	"fmt"

	"github.com/seyedmo30/order_management/internal/config"
	"github.com/seyedmo30/order_management/internal/dto"
	"github.com/seyedmo30/order_management/pkg"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database connection pool and thread-safety mechanisms
var (
	databaseConfig config.DatabaseConfig
	db             *gorm.DB
)

// DB ensures a singleton instance of the database connection.
func DB() *gorm.DB {
	return db
}

// orderManagementRepository is the structure holding DB configuration for repository operations.
type orderManagementRepository struct {
	config config.DatabaseConfig
}

// Ensure orderManagementRepository implements the Repository interface.
// var _ interfaces.Repository = (*orderManagementRepository)(nil)

// NewOrderManagementRepository initializes a new orderManagementRepository with the provided configuration.
func NewOrderManagementRepository(config config.DatabaseConfig) *orderManagementRepository {
	databaseConfig = config // Store the config for the singleton
	var err error
	// Initialize the database only once
	db, err = SetupDB(databaseConfig)
	if err != nil {
		pkg.GetLogger().Error("Error initializing DB: ", err)
	} // Ensure DB instance is initialized

	return &orderManagementRepository{config: config}
}

// SetupDB configures and returns a GORM DB connection using the provided database config.
func SetupDB(config config.DatabaseConfig) (*gorm.DB, error) {

	// Open the underlying SQL connection (SQLite doesn't need credentials)
	sqlDb, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open SQL connection: %w", err)
	}

	// Ping the database to verify connection
	if err := sqlDb.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	logMode := logger.Info
	switch config.LogLevel {
	case "INFO":
		logMode = logger.Info
	case "ERROR":
		logMode = logger.Error
	case "WARNING":
		logMode = logger.Warn

	}
	// Initialize GORM with the SQLite connection
	db, err := gorm.Open(sqlite.New(sqlite.Config{Conn: sqlDb}), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logMode),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open GORM connection: %w", err)
	}

	// Auto-migrate the BaseOrder model
	err = db.AutoMigrate(&dto.BaseOrder{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	return db, nil
}
