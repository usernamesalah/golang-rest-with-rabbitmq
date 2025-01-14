package database

import (
	"context"
	"fmt"
	"tenant/infrastructure/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitializeDatabase to initialize database
func InitializeDatabase(conf *config.Config) (*pgxpool.Pool, error) {
	// Menyiapkan connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DBName,
		conf.Database.SSLMode,
	)

	// Membuat pool koneksi
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
