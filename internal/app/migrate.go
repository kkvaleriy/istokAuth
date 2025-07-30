//go:build migrate

package app

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	attempts = 5
	timeout  = time.Second * 3
)

func init() {
	dbURL := os.Getenv("PG_MIGRATE_URL")
	if len(dbURL) < 1 {
		log.Fatal("migrate: env not declared: PG_MIGRATE_URL")
	}

	var (
		currentAttempt int
		err            error
		m              *migrate.Migrate
	)

	for currentAttempt < attempts {
		m, err = migrate.New("file://migrations", dbURL)
		if err == nil {
			break
		}

		log.Printf("migrate: postgres is trying to connect, attempts left: %d", attempts-currentAttempt-1)
		currentAttempt++
		<-time.After(timeout)
	}

	if err != nil {
		log.Fatalf("migrate: postgres connect error: %s", err.Error())
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("migrate: no change")
		return
	}

	log.Printf("migrate: up success")
}
