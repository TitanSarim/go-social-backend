package main

import (
	"time"

	"github.com/TitanSarim/go-social-backend/internal/db"
	"github.com/TitanSarim/go-social-backend/internal/env"
	"github.com/TitanSarim/go-social-backend/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			GoSocial API
//	@description	API FOR GGpSocial, a social network for Developers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
// 
//	@securityDefinitions.apiKey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description


func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8100"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8100"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "user=postgres dbname=go_social password=12345 sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime: env.GetString("DB_CONN_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp: time.Hour * 24 * 3, // 3 days
		},
	}

	// Logger

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
        logger.Fatal(err)
    }

	defer db.Close()
	logger.Info("Connected to the database")

	store := store.NewPostgresStorage(db)

	app := &application{
		config: cfg,
		store: store,
		logger: logger,
	}


	mux := app.mount()
	logger.Fatal(app.run(mux))
}