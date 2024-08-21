package integration_test

import (
	"crypto/tls"
	"log"

	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/network/standard"
)

var hertzClient *client.Client

func init() {
	var err error
	clientCfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	hertzClient, err = client.NewClient(
		client.WithTLSConfig(clientCfg),
		client.WithDialer(standard.NewDialer()),
	)
	if err != nil {
		log.Fatalf("Failed to create Hertz client: %v", err)
	}
}
