package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JoaoVitorStein/folder-manager/pkg/importer"
	"github.com/JoaoVitorStein/folder-manager/pkg/manager"
	"github.com/JoaoVitorStein/folder-manager/pkg/server"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()
	databaseConnectionString := os.Getenv("DB_CONNECTION_STRING")
	db, err := sqlx.Open("postgres", databaseConnectionString)
	if err != nil {
		fmt.Printf("error creating db %v \n", err)
		os.Exit(1)
	}
	// wait db up
	for i := 0; i < 5; i++ {
		err := db.Ping()
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	migration, err := migrate.New("file:///app/migrations/", databaseConnectionString)
	if err != nil {
		fmt.Printf("could not migrate: %s", err)
		os.Exit(1)
	}
	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		fmt.Printf("could not migrate: %s", err)
		os.Exit(1)
	}
	m := manager.New(db)
	imp := importer.New(m)
	file, err := os.Open("/app/csv_file.csv")
	if err != nil {
		fmt.Printf("error reading file %v \n", err)
		os.Exit(1)
	}
	defer file.Close()
	err = imp.ImportFile(bufio.NewReader(file))
	if err != nil {
		fmt.Printf("error importing file %v \n", err)
		os.Exit(1)
	}
	s := server.New(m)
	s.Start(ctx)
}
