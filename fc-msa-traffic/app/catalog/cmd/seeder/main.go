package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/catalog/internal/catalog"
	"github.com/go-faker/faker/v4"
)

func main() {
	configFlag := flag.String("config", "configs/catalog.yaml", "Path to config file")
	flag.Parse()

	cfg, err := catalog.NewConfig(*configFlag)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		os.Exit(1)
	}

	db, err := catalog.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
		os.Exit(1)
	}

	// if err := db.AutoMigrate(&catalog.Product{}); err != nil {
	// 	log.Fatalf("Failed to auto-migrate: %v", err)
	// }

	// Seed products
	const numberOfProducts = 20
	for i := 0; i < numberOfProducts; i++ {
		// Generate random price
		priceInt, err := faker.RandomInt(10, 1000)
		if err != nil || len(priceInt) == 0 {
			log.Printf("failed to generate random price: %v", err)
			continue
		}

		// Generate random stock
		stockInt, err := faker.RandomInt(0, 500)
		if err != nil || len(stockInt) == 0 {
			log.Printf("failed to generate random stock: %v", err)
			continue
		}

		product := catalog.Product{
			ID:          faker.UUIDHyphenated(),
			Name:        faker.Word(),
			Description: faker.Sentence(),
			Price:       float64(priceInt[0]),
			Stock:       stockInt[0],
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := db.Create(&product).Error; err != nil {
			log.Printf("Failed to create product: %v", err)
		} else {
			fmt.Printf("Seeded product: %s (Price: %.2f, Stock: %d)\n", product.Name, product.Price, product.Stock)
		}
	}

	fmt.Println("Product seeding completed!")
}
