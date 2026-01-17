package persistence

import (
	"fmt"
	"time"

	"github.com/alexduzi/challengepismo/internal/infrastructure/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const driverName = "pgx"

func ConnectDB(conf *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.DBName,
		conf.Database.SSLMode,
	)

	db, err := sqlx.Connect(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(conf.Database.MaxConns)
	db.SetMaxIdleConns(conf.Database.MinConns)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func CloseDB(db *sqlx.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
