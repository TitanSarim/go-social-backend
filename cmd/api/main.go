package main

import (
	"log"

	"github.com/TitanSarim/go-social-backend/internal/db"
	"github.com/TitanSarim/go-social-backend/internal/env"
	"github.com/TitanSarim/go-social-backend/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8100"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "user=postgres dbname=go_social password=12345 sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime: env.GetString("DB_CONN_MAX_IDLE_TIME", "15m"),
		},
	}
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
        log.Panic(err)
    }

	defer db.Close()
	log.Println("Connected to the database")

	store := store.NewPostgresStorage(db)

	app := &application{
		config: cfg,
		store: store,
	}


	mux := app.mount()
	log.Fatal(app.run(mux))
}