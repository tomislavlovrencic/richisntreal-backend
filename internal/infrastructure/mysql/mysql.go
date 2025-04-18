package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"richisntreal-backend/cmd/config"
)

type MySQL struct {
	DB *sqlx.DB
}

// NewMySQL dials MySQL using your config and returns a *MySQL.
func NewMySQL(cfg config.MySQL) (*MySQL, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &MySQL{DB: db}, nil
}
