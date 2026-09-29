package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/sowncns/k3s-deploy-platform/internal/config"
)

const migrationsPath = "file://migrations"

func main() {
	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	cfg := config.Load()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("create migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		log.Fatalf("init migrate: %v", err)
	}

	switch cmd {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "drop":
		err = m.Drop()
	case "version":
		version, dirty, verErr := m.Version()
		if verErr != nil {
			log.Fatalf("version: %v", verErr)
		}
		fmt.Printf("version=%d dirty=%v\n", version, dirty)
		return
	case "goto":
		if len(os.Args) < 3 {
			log.Fatal("usage: migrate goto <version>")
		}
		v, parseErr := strconv.ParseUint(os.Args[2], 10, 64)
		if parseErr != nil {
			log.Fatalf("invalid version: %v", parseErr)
		}
		err = m.Migrate(uint(v))
	default:
		log.Fatalf("unknown command %q (expected: up | down | drop | version | goto <n>)", cmd)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrate %s: %v", cmd, err)
	}

	fmt.Printf("migrate %s: done\n", cmd)
}
