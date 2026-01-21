package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ReyviRahman/social/internal/store" // Sesuaikan dengan module path Anda
)

func Seed(s store.Storage) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// 1. Seed Users
	users := seedUsers(ctx, s)
	
	// 2. Seed Posts
	seedPosts(ctx, s, users)

	log.Println("Seeding completed successfully!")
}

func seedUsers(ctx context.Context, s store.Storage) []*store.User {
	usernames := []string{"reyvi", "rahman", "alice", "bob", "charlie"}
	users := make([]*store.User, len(usernames))

	for i, name := range usernames {
		user := &store.User{
			Username: name,
			Email:    fmt.Sprintf("%s@example.com", name),
			Password: "password123", // Dalam realita, ini harus di-hash
		}

		if err := s.Users.Create(ctx, user); err != nil {
			log.Fatalf("error seeding user: %v", err)
		}
		users[i] = user
	}
	log.Printf("Seeded %d users", len(users))
	return users
}

func seedPosts(ctx context.Context, s store.Storage, users []*store.User) {
	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < 20; i++ {
		user := users[rand.Intn(len(users))]
		
		post := &store.Post{
			UserID:  user.ID,
			Title:   fmt.Sprintf("Judul Postingan ke-%d", i),
			Content: fmt.Sprintf("Ini adalah konten postingan menarik dari %s.", user.Username),
			Tags:    []string{"golang", "pamo", "uneeenda"},
		}

		if err := s.Posts.Create(ctx, post); err != nil {
			log.Fatalf("error seeding post: %v", err)
		}
	}
	log.Printf("Seeded 20 posts")
}