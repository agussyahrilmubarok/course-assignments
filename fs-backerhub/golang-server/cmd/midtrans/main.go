package main

import (
	"crypto/sha512"
	"flag"
	"fmt"

	"example.com/backend/pkg/config"
	"example.com/backend/pkg/logger"
)

func main() {
	orderID := flag.String("order_id", "", "Midtrans Order ID")
	grossAmount := flag.String("gross_amount", "", "Gross Amount from Midtrans")

	configPath := flag.String("config", "backerhub_stage.json", "Path to config file")
	flag.Parse()

	// Load configuration file
	cfg := config.Load(*configPath)

	// Initialize logger
	log := logger.NewZeroLogger(cfg.Logging)

	if *orderID == "" || *grossAmount == "" {
		log.Error().Msgf("Both --order_id and --gross_amount must be provided")
		return
	}

	signature := GenerateMidtransSignature(*orderID, *grossAmount, cfg.Midtrans.ServerKey)
	log.Info().Msgf("Signature Key: %s", signature)
}

func GenerateMidtransSignature(orderID, grossAmount, serverKey string) string {
	statusCode := "200" // Midtrans always uses "200" for successful notification
	raw := orderID + statusCode + grossAmount + serverKey

	hash := sha512.Sum512([]byte(raw))
	signature := fmt.Sprintf("%x", hash[:])
	return signature
}
