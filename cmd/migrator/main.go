package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"trash_project/migrations"
	appconfig "trash_project/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func dsnFromCfg(cfg *appconfig.Configuration) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PostgreSQL.User, cfg.PostgreSQL.Pass, cfg.PostgreSQL.Host, cfg.PostgreSQL.Port, cfg.PostgreSQL.Database, cfg.PostgreSQL.SSLMode,
	)
}

func main() {
	configPath := flag.String("config", "data/config.yml", "path to config yaml")
	flag.Parse()

	appconfig.Setup(*configPath)
	cfg := appconfig.GetConfig()

	if err := up(cfg); err != nil {
		log.Printf("migration failed: %v", err)
		os.Exit(1)
	}

	log.Println("migration applied")
}

func up(cfg *appconfig.Configuration) error {
	dsn := dsnFromCfg(cfg)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db: %v", err)
	}
	defer db.Close()

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("iofs source: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres driver: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate init: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %v", err)
	}

	return nil
}
