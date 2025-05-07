package db

import (
	"database/sql"
	"fmt"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/cirius-go/portfolio-server/util"
)

type PostgresConfig struct {
	DSN      string            `envconfig:"DSN"`
	Host     string            `envconfig:"HOST"`
	Port     int               `envconfig:"PORT"`
	Username string            `envconfig:"USERNAME"`
	Password string            `envconfig:"PASSWORD"`
	Database string            `envconfig:"DATABASE"`
	Args     util.QueryDecoder `envconfig:"ARGS"`
	LogLevel logger.LogLevel   `envconfig:"LOG_LEVEL"` // default 3
}

// buildDSN string.
func (c *PostgresConfig) buildDSN() string {
	u := url.URL{}

	u.Scheme = "postgres"
	if c.Port != 0 {
		u.Host = fmt.Sprintf("%s:%d", c.Host, c.Port)
	} else {
		u.Host = c.Host
	}

	u.User = url.UserPassword(c.Username, c.Password)
	u.Path = c.Database

	q := u.Query()
	for k, v := range c.Args {
		if q.Has(k) {
			q.Add(k, v)
		} else {
			q.Set(k, v)
		}
	}

	return u.String()
}

// Postgres contains configured gorm + connection.
type Postgres struct {
	Conn *sql.DB
	DB   *gorm.DB
}

// NewPostgres connects the database and config gorm.
func NewPostgres(cfg PostgresConfig) (*Postgres, error) {
	gormCfg := &gorm.Config{
		Logger:                                   logger.Default.LogMode(cfg.LogLevel),
		IgnoreRelationshipsWhenMigrating:         true,
		DisableForeignKeyConstraintWhenMigrating: true,
		CreateBatchSize:                          1000,
	}

	addr := cfg.DSN
	if cfg.DSN == "" {
		addr = cfg.buildDSN()
	}
	if addr == "" {
		return nil, fmt.Errorf("empty dsn")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  addr,
		PreferSimpleProtocol: true,
	}), gormCfg)
	if err != nil {
		return nil, err
	}

	p := &Postgres{}
	p.Conn, err = db.DB()
	if err != nil {
		return nil, err
	}

	p.DB = db
	return p, nil
}
