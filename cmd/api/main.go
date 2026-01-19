package main

import (
	"log"
	"os"
	"strconv"

	"github.com/ReyviRahman/social/internal/db"
	"github.com/ReyviRahman/social/internal/store"
	"github.com/joho/godotenv"
)

func getEnvInt(key string, fallback int) int {
	s := os.Getenv(key)
	if s == "" {
		return fallback
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return v
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: os.Getenv("ADDR"),
		db: dbConfig{
			addr:         os.Getenv("DB_ADDR"),
			maxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  os.Getenv("DB_MAX_IDLE_TIME"),
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
	log.Println("database connection pool established")
	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
