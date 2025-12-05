package main

import (
	"flag"
	"log"
	"time"

	"example.com.backend/internal/domain"
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

	// Connect To Postgres
	db, err := gorm.Open(postgres.Open(*dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// Migrate Postgres Tables
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Campaign{},
		&domain.CampaignImage{},
		&domain.Transaction{},
	); err != nil {
		log.Fatal("failed to migrate database: ", err)
	}
	log.Println("Migration success")

	if err := seedAdmin(db); err != nil {
		log.Fatalf("failed to seed admin: %v", err)
	}

	if err := seedUserCase(db); err != nil {
		log.Fatalf("failed to seed users: %v", err)
	}

	if err := seedCampaignTxCase1(db); err != nil {
		log.Fatalf("failed to seed campaign and transaction case 1: %v", err)
	}
}

func seedAdmin(db *gorm.DB) error {
	occupation := "Administrator"
	defaultImage := "default.png"
	hashed, err := bcrypt.GenerateFromPassword([]byte("P@ssw0rd!"), bcrypt.DefaultCost)
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
		return err
	}

	return nil
}

func seedUserCase(db *gorm.DB) error {
	// Create User John Doe, Jane Smith, Bob Johnson, Alice Williams, Michael Brown
	defaultPassword := "P@ssw0rd!"
	defaultImage := "default.png"

	users := []struct {
		Name       string
		Email      string
		Password   string
		Role       string
		Occupation *string
		ImageName  *string
	}{
		{"John Doe", "john.doe@example.com", defaultPassword, string(domain.RoleUser), strPtr("Software Engineer"), &defaultImage},
		{"Jane Smith", "jane.smith@example.com", defaultPassword, string(domain.RoleUser), strPtr("Designer"), &defaultImage},
		{"Bob Johnson", "bob.johnson@example.com", defaultPassword, string(domain.RoleUser), strPtr("Product Manager"), &defaultImage},
		{"Alice Williams", "alice.williams@example.com", defaultPassword, string(domain.RoleUser), strPtr("Content Writer"), &defaultImage},
		{"Michael Brown", "michael.brown@example.com", defaultPassword, string(domain.RoleUser), strPtr("Marketer"), &defaultImage},
	}

	for _, u := range users {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := domain.User{
			Name:       u.Name,
			Email:      u.Email,
			Password:   string(hashed),
			Role:       u.Role,
			Occupation: u.Occupation,
			ImageName:  u.ImageName,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := db.FirstOrCreate(&user, domain.User{Email: u.Email}).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedCampaignTxCase1(db *gorm.DB) error {
	// John Doe create 2 campaigns
	// Jane Smith, Alice Williams create transaction to john doe campaign with status pending
	// Bob Johnson, Bob Johnson create transaction to john doe campaign with status paid

	var john, jane, alice, bob domain.User
	if err := db.Where("email = ?", "john.doe@example.com").First(&john).Error; err != nil {
		return err
	}
	if err := db.Where("email = ?", "jane.smith@example.com").First(&jane).Error; err != nil {
		return err
	}
	if err := db.Where("email = ?", "alice.williams@example.com").First(&alice).Error; err != nil {
		return err
	}
	if err := db.Where("email = ?", "bob.johnson@example.com").First(&bob).Error; err != nil {
		return err
	}

	campaigns := []domain.Campaign{
		{
			Title:            "Open Source Platform",
			ShortDescription: "Help build an open source platform",
			Description:      "Support developers worldwide",
			GoalAmount:       10000.00,
			CurrentAmount:    500.00,
			BackerCount:      2,
			Perks:            "Sticker, T-Shirt",
			Slug:             "open-source-platform",
			UserID:           john.ID,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
		{
			Title:            "Community Library",
			ShortDescription: "Build a community library",
			Description:      "Provide books for children",
			GoalAmount:       5000.00,
			CurrentAmount:    200.00,
			BackerCount:      1,
			Perks:            "Bookmark, Thank You Note",
			Slug:             "community-library",
			UserID:           john.ID,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	}

	for _, c := range campaigns {
		if err := db.FirstOrCreate(&c, domain.Campaign{Slug: c.Slug}).Error; err != nil {
			return err
		}

		// Campaign Image
		img := domain.CampaignImage{
			ImageName:  "default.jpeg",
			IsPrimary:  true,
			CampaignID: c.ID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := db.FirstOrCreate(&img, domain.CampaignImage{CampaignID: c.ID}).Error; err != nil {
			return err
		}
	}

	var campaign1 domain.Campaign
	if err := db.Where("slug = ?", "open-source-platform").First(&campaign1).Error; err != nil {
		return err
	}

	transactions := []domain.Transaction{
		{
			Amount:     100.00,
			Status:     string(domain.StatusPending),
			Note:       "Donation from Jane",
			Reference:  "ref_jane_1",
			UserID:     jane.ID,
			CampaignID: campaign1.ID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Amount:     150.00,
			Status:     string(domain.StatusPending),
			Note:       "Donation from Alice",
			Reference:  "ref_alice_1",
			UserID:     alice.ID,
			CampaignID: campaign1.ID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Amount:     200.00,
			Status:     string(domain.StatusPaid),
			Note:       "Donation from Bob",
			Reference:  "ref_bob_1",
			UserID:     bob.ID,
			CampaignID: campaign1.ID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Amount:     50.00,
			Status:     string(domain.StatusPaid),
			Note:       "Second donation from Bob",
			Reference:  "ref_bob_2",
			UserID:     bob.ID,
			CampaignID: campaign1.ID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	for _, tx := range transactions {
		if err := db.FirstOrCreate(&tx, domain.Transaction{Reference: tx.Reference}).Error; err != nil {
			return err
		}
	}

	return nil
}

func strPtr(s string) *string {
	return &s
}
