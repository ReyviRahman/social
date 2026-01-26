package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ReyviRahman/social/internal/db"
	"github.com/ReyviRahman/social/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	addr := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable", 
			dbUser, 
			dbPassword, 
			dbName,
	)
	
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store)
}
