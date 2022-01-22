package manager

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
)

// Database connection
const (
	dbName = "test"
	user   = "postgres"
	pass   = "s3cr3tp4ssw0rd"
	host   = "localhost"
)

type dockerPg struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (d *dockerPg) Start(t *testing.T) *sqlx.DB {

	pool, resource, err := createDockerImage()
	d.pool = pool
	d.resource = resource

	if err != nil {
		t.Fatalf("could not connect to docker: %s", err)
	}
	databaseConnectionString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, pass, host, d.resource.GetPort("5432/tcp"), dbName)
	db, _ := sqlx.Open("postgres", databaseConnectionString)
	err = waitDockerUp(pool, db)
	if err != nil {
		_ = tearDownDocker(pool, resource)
		t.Fatalf("could not connect to docker: %s", err)
	}
	migrationPath, _ := filepath.Abs("../../migrations")
	migrate, err := migrate.New(fmt.Sprintf("file://%v", migrationPath), databaseConnectionString)
	if err != nil {
		_ = tearDownDocker(pool, resource)
		t.Fatalf("could not migrate: %s", err)
	}
	err = migrate.Up()
	if err != nil {
		_ = tearDownDocker(pool, resource)
		t.Fatalf("could not migrate: %s", err)
	}
	dockerDb, err := sqlx.Connect("postgres", databaseConnectionString)
	if err != nil {
		t.Fatalf("could not connect to docker db: %s", err)
	}

	return dockerDb
}

func (d dockerPg) Stop(t *testing.T) {
	err := tearDownDocker(d.pool, d.resource)

	if err != nil {
		t.Fatalf("could not connect to docker: %s", err)
	}

}

func createDockerImage() (pool *dockertest.Pool, resource *dockertest.Resource, err error) {
	pool, err = dockertest.NewPool("")
	if err != nil {
		return nil, nil, err
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11.11",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%v", user),
			fmt.Sprintf("POSTGRES_PASSWORD=%v", pass),
			fmt.Sprintf("POSTGRES_DB=%v", dbName),
		},
	}

	resource, err = pool.RunWithOptions(&opts)

	if err != nil {
		return nil, nil, err
	}
	return
}

func waitDockerUp(pool *dockertest.Pool, dbInstance *sqlx.DB) error {
	if err := pool.Retry(func() error {
		return dbInstance.Ping()
	}); err != nil {
		return err
	}
	return nil
}

func tearDownDocker(pool *dockertest.Pool, resource *dockertest.Resource) error {
	return pool.Purge(resource)
}
