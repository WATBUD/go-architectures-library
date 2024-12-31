package db

import (
	"context"
	"fmt"
	"log"
	"spt/config"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm/schema"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Service represents a service that interacts with a database.
type Service interface {
	Health() map[string]string
	GetDB() *gorm.DB
	//Close() error
}

type service struct {
	db *gorm.DB
}

var (
	dbInstance *service
	once       sync.Once
)

func Instance() Service {
	once.Do(func() {
		// Load database configuration
		dbConfig := config.LoadDBConfig()

		// DSN (Data Source Name) from the config
		dsn := dbConfig.DSN()

		// Initialize GORM DB connection
		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         256,   // default size for string fields
			DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
			DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			DontSupportRenameColumn:   true,  // use change when rename column, rename column not supported before MySQL 8, MariaDB
			SkipInitializeWithVersion: false, // autoconfigure based on currently MySQL version
		}), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})

		if err != nil {
			log.Fatal("failed to connect to database: ", err)
		}

		// Get the underlying sql.DB to access its stats
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal("failed to get underlying sql.DB: ", err)
		}

		// Set connection pool settings
		sqlDB.SetConnMaxLifetime(0)
		sqlDB.SetMaxIdleConns(50)
		sqlDB.SetMaxOpenConns(50)

		dbInstance = &service{
			db: db,
		}
	})
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	sqlDB, err := s.db.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err)
		return stats
	}

	// Ping the database
	err = sqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err)
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}
func (s *service) GetDB() *gorm.DB {
	return s.db
}

// Close closes the database connection.
// Note: GORM automatically manages the connection pool, so explicit closing is not required.
//func (s *service) Close() error {
//	log.Printf("Disconnected from database: %s", dbname)
//	return s.sqlDB.Close()
//}
