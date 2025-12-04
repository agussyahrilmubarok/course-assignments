package config

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type Midtrans struct {
	ServerKey string `mapstructure:"server_key"`
	ClientKey string `mapstructure:"client_key"`
	IsProd    bool   `mapstructure:"is_prod"`
}

func NewMidtrans(cfg *Midtrans) snap.Client {
	// 1. Setup Midtrans Configuration
	midtrans.ServerKey = cfg.ServerKey // YOUR-VT-SERVER-KEY
	midtrans.ClientKey = cfg.ClientKey // YOUR-VT-CLIENT-KEY
	midtrans.Environment = midtrans.Sandbox
	if cfg.IsProd {
		midtrans.Environment = midtrans.Production
	}

	// 2. Create Snap Client
	snapClient := snap.Client{}
	snapClient.New(midtrans.ServerKey, midtrans.Environment)

	return snapClient
}
