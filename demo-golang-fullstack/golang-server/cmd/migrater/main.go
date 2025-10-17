package main

import (
	"flag"
	"log"
	"time"

	"example.com/backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := flag.String(
		"dsn",
		"host=127.0.0.1 user=backeruser password=backerpass dbname=backerhub port=5432 sslmode=disable TimeZone=Asia/Jakarta",
		"PostgreSQL DSN",
	)
	flag.Parse()

	db, err := gorm.Open(postgres.Open(*dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Campaign{},
		&domain.CampaignImage{},
		&domain.Transaction{},
	); err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	log.Println("✅ Migration success")
	occupation := "Administrator"
	defaultImage := "default.png"
	hashed, err := bcrypt.GenerateFromPassword([]byte("hashedpassword123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	user := domain.User{
		Name:       "Administrator",
		Email:      "admin@backerhub.org",
		Password:   string(hashed),
		Role:       string(domain.RoleAdmin),
		Occupation: &occupation,
		ImageName:  &defaultImage,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := db.FirstOrCreate(&user, domain.User{Email: user.Email}).Error; err != nil {
		log.Fatalf("failed to seed admin: %v", err)
	}
}

func Seed(db *gorm.DB) {
	hashed, err := bcrypt.GenerateFromPassword([]byte("hashedpassword123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	user := domain.User{
		ID:        "user_1",
		Name:      "Potato Human",
		Email:     "potato@example.com",
		Password:  string(hashed),
		Role:      "USER",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.FirstOrCreate(&user, domain.User{Email: user.Email}).Error; err != nil {
		log.Fatalf("failed to seed user: %v", err)
	}

	// Seed Campaign
	campaign := domain.Campaign{
		ID:               "camp_1",
		Title:            "Build Open Source Platform",
		ShortDescription: "Help us build an amazing open source project",
		Description:      "This project aims to support developers worldwide.",
		GoalAmount:       10000.00,
		CurrentAmount:    500.00,
		BackerCount:      5,
		Perks:            "Sticker, T-Shirt, Early Access",
		Slug:             "build-open-source-platform",
		UserID:           user.ID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := db.FirstOrCreate(&campaign, domain.Campaign{Slug: campaign.Slug}).Error; err != nil {
		log.Fatalf("failed to seed campaign: %v", err)
	}

	// Seed Campaign Image
	campaignImage := domain.CampaignImage{
		ID:         "img_1",
		ImageName:  "https://example.com/campaign1.png",
		IsPrimary:  true,
		CampaignID: campaign.ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := db.FirstOrCreate(&campaignImage, domain.CampaignImage{ID: campaignImage.ID}).Error; err != nil {
		log.Fatalf("failed to seed campaign image: %v", err)
	}

	// Seed Transaction
	tx := domain.Transaction{
		ID:         "tx_1",
		Amount:     100.00,
		Status:     "SUCCESS",
		Note:       "First donation",
		Reference:  "ref_123",
		UserID:     user.ID,
		CampaignID: campaign.ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := db.FirstOrCreate(&tx, domain.Transaction{ID: tx.ID}).Error; err != nil {
		log.Fatalf("failed to seed transaction: %v", err)
	}

	log.Println("✅ Seeding success")
}
