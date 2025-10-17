package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/account/internal/account"
	"github.com/go-faker/faker/v4"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	configFlag := flag.String("config", "configs/account.yaml", "Path to config file")
	flag.Parse()

	cfg, err := account.NewConfig(*configFlag)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		os.Exit(1)
	}

	db, err := account.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	const numberOfUsers = 10
	for i := 0; i < numberOfUsers; i++ {
		password := "password123"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
			continue
		}

		user := account.User{
			ID:        faker.UUIDHyphenated(),
			Name:      faker.Name(),
			Email:     fmt.Sprintf("user%d_%s@example.com", i, faker.Username()),
			Password:  string(hashedPassword),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create user: %v", err)
		} else {
			fmt.Printf("Seeded user: %s <%s>\n", user.Name, user.Email)
		}
	}

	fmt.Println("Seeding completed")
}
