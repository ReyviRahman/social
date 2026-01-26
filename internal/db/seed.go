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
	log.Println("Seeding completed successfully!")
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
