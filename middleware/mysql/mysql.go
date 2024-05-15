package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Addr         string `toml:"addr"`
	User         string `toml:"user"`
	Pwd          string `toml:"pwd"`
	Dbname       string `toml:"dbname"`
	MaxOpenConns int    `toml:"max_open_conns"`
	MaxIdleConns int    `toml:"max_idle_conns"`
}

func NewMysql(c *Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", formatDataSource(c.Addr, c.User, c.Pwd, c.Dbname))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)
	return db, nil
}

func formatDataSource(addr string, user string, pwd string, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&timeout=2s&parseTime=true&useSSL=no", user, pwd, addr, dbname)
}
