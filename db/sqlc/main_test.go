package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	dsn := "postgres://root:secret@localhost:5432/db?sslmode=disable"

	testDB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
