package database

import (
	"context"
	"fmt"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"grpc-server/ent"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	Client *ent.Client
}

func New(connString string) (*DB, error) {
	drv, err := entsql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := drv.DB()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))

	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL with Ent")

	return &DB{Client: client}, nil
}

func (db *DB) Close() error {
	return db.Client.Close()
}
