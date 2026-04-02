package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/fyk7/code-snippets-app/app/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Config) *gorm.DB {
	val := url.Values{}
	val.Add("charset", "utf8mb4")
	val.Add("parseTime", "true")
	connection := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		panic(err)
	}

	slog.Info("database connection established")
	return dbConn
}

// SQLDBFromGorm extracts the underlying *sql.DB for lifecycle management.
func SQLDBFromGorm(db *gorm.DB) (*sql.DB, error) {
	return db.DB()
}
