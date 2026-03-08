package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"example.com/coupon/internal/coupon"
	"example.com/coupon/pkg/config"
	"gorm.io/gorm"
)

func main() {
	configFlag := flag.String("config", "configs/config.yaml", "Path to config file")
	outputDir := flag.String("output", "exports", "path ke direktori output CSV")
	flag.Parse()

	cfg, err := config.NewConfig(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	db, err := config.NewPostgres(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(*outputDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory: %v", err)
		os.Exit(1)
	}

	timestamp := time.Now().Format("20060102150405")
	policiesFile := filepath.Join(*outputDir, fmt.Sprintf("%s-coupon-policies.csv", timestamp))
	couponsFile := filepath.Join(*outputDir, fmt.Sprintf("%s-coupons.csv", timestamp))

	if err := exportCouponPolicies(db, policiesFile); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to export coupon policies: %v\n", err)
		os.Exit(1)
	}

	if err := exportCoupons(db, couponsFile); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to export coupons: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Export selesai.\nPolicies: %s\nCoupons: %s\n", policiesFile, couponsFile)
}

func exportCouponPolicies(db *gorm.DB, outputPath string) error {
	var policies []coupon.CouponPolicy
	if err := db.Find(&policies).Error; err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"id", "code", "name", "description",
		"total_quantity", "start_time", "end_time",
		"discount_type", "discount_value",
		"minimum_order_amount", "maximum_discount_amount",
		"created_at", "updated_at",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, p := range policies {
		record := []string{
			p.ID,
			p.Code,
			p.Name,
			p.Description,
			fmt.Sprintf("%d", p.TotalQuantity),
			p.StartTime.Format(time.RFC3339),
			p.EndTime.Format(time.RFC3339),
			string(p.DiscountType),
			fmt.Sprintf("%d", p.DiscountValue),
			fmt.Sprintf("%d", p.MinimumOrderAmount),
			fmt.Sprintf("%d", p.MaximumDiscountAmount),
			p.CreatedAt.Format(time.RFC3339),
			p.UpdatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func exportCoupons(db *gorm.DB, outputPath string) error {
	var coupons []coupon.Coupon
	if err := db.Preload("CouponPolicy").Find(&coupons).Error; err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"id", "code", "status", "used_at",
		"user_id", "order_id", "coupon_policy_id",
		"created_at", "updated_at",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, c := range coupons {
		usedAt := ""
		if c.UsedAt != nil {
			usedAt = c.UsedAt.Format(time.RFC3339)
		}
		orderId := ""
		if c.OrderID != nil {
			orderId = *c.OrderID
		}
		record := []string{
			c.ID,
			c.Code,
			string(c.Status),
			usedAt,
			c.UserID,
			orderId,
			c.CouponPolicyID,
			c.CreatedAt.Format(time.RFC3339),
			c.UpdatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
