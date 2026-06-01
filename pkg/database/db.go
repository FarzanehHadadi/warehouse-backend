package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"warehouse/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbConfig struct {
	Host            string
	Port            string
	Username        string
	password        string
	DbName          string
	SSLMode         string
	MaxOpenConn     int
	MaxIdleConn     int
	ConnMaxLifeTime time.Duration
	ConnMaxIdleTime time.Duration
}

func DefaultConfig() dbConfig {
	return dbConfig{
		Host:            "localhost",
		Port:            "5432",
		DbName:          "warehouse",
		Username:        "postgres",
		password:        "123456",
		SSLMode:         "disable",
		MaxOpenConn:     25,
		MaxIdleConn:     10,
		ConnMaxLifeTime: time.Minute * 5,
		ConnMaxIdleTime: time.Minute * 10,
	}
}

func NewPostgresConfiguration(cfg *dbConfig) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.Port, cfg.Username, cfg.password, cfg.DbName, cfg.SSLMode)
	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true,  // Improves performance by 30%+
		PrepareStmt:            true,  // Caches prepared statements
		AllowGlobalUpdate:      false, // Prevents accidental mass updates
		DisableAutomaticPing:   false, // Keeps connection alive check
	}
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	// log.Println("db,err", db, err)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Category{}, &models.Unit{}, &models.Department{}); err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql db %w", err)

	}
	// Connection pool configuration for production
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifeTime)

	//ٰverify connection is alive
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)

	}
	log.Printf("✅ Database connection established (max_open=%d, max_idle=%d)",
		cfg.MaxOpenConn, cfg.MaxIdleConn)
	return db, nil
}
