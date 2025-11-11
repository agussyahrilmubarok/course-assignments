package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/catalog/internal/catalog"
	"github.com/agussyahrilmubarok/gox/pkg/xconfig/xviper"
	"github.com/agussyahrilmubarok/gox/pkg/xgorm"
	"github.com/go-faker/faker/v4"
)

func main() {
	configFlag := flag.String("config", "configs/catalog.yaml", "Path to config file")
	flag.Parse()

	vCfg, err := xviper.NewConfig(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	var cfg *catalog.Config
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
