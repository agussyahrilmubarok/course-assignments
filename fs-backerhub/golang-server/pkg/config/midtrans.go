package config

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type Midtrans struct {
	ServerKey string `json:"server_key" mapstructure:"server_key"`
	ClientKey string `json:"client_key" mapstructure:"client_key"`
	IsProd    bool   `json:"is_prod" mapstructure:"is_prod"`
}

func NewMidtransSnapClient(cfg Midtrans) snap.Client {
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
