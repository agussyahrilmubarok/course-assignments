package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/account/internal/account"
	"github.com/agussyahrilmubarok/gox/pkg/xconfig/xviper"
	"github.com/agussyahrilmubarok/gox/pkg/xgorm"
	"github.com/go-faker/faker/v4"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	configFlag := flag.String("config", "configs/account.yaml", "Path to config file")
	flag.Parse()

	vCfg, err := xviper.NewConfig(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	var cfg *account.Config
	if err := vCfg.Unmarshal(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	db, err := xgorm.NewGorm("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DbName,
		cfg.Postgres.SslMode,
	), &xgorm.Options{
		MaxOpenConns:    cfg.Postgres.MaxOpenConns,
		MaxIdleConns:    cfg.Postgres.MaxIdleConns,
		ConnMaxLifetime: cfg.Postgres.ConnMaxLifetime,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
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
