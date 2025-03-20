package psql

import (
	"fmt"
	"memesearch/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func connect(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname,
	)
	fmt.Println("Connection string:", connStr) // Отладочный вывод
	return sqlx.Connect("postgres", connStr)
}
