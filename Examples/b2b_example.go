//go:build ignore
// +build ignore

// Example: how to use the B2B helpers, BusinessBuyGoodsService and the webhook handler.
// This file is marked with a build ignore tag so it's not compiled with `go test` or `go build`.

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/venomous-maker/go-mpesa/Abstracts"
	"github.com/venomous-maker/go-mpesa/Services"
)

func main() {
	// Create config (replace placeholders with real credentials for real runs)
	cfg, err := Abstracts.NewMpesaConfig(
		"YOUR_CONSUMER_KEY",
		"YOUR_CONSUMER_SECRET",
		Abstracts.Sandbox,
		ptr("174379"), // business shortcode (sandbox example)
		nil,           // passKey not required for B2B
		nil,           // security credential will be set per-initiator below
		nil, nil,
	)
	if err != nil {
		log.Fatalf("failed to create config: %v", err)
	}

	client := Abstracts.NewApiClient(cfg)

	// Create BusinessBuyGoods service
	buy := Services.NewBusinessBuyGoodsService(cfg, client)

	// Set initiator and security credential (encrypts internally)
	buy.SetInitiator("API_Username")
	if err := buy.SetSecurityCredential("InitiatorPlainPassword"); err != nil {
		log.Fatalf("failed to set security credential: %v", err)
	}

	// Configure payment
	buy.SetAmount(100)
	buy.SetPartyA("123456") // your shortcode (or set cfg business code earlier)
	buy.SetPartyB("600000") // destination merchant till/store
	buy.SetAccountReference("INV-1001")
	buy.SetRemarks("Payment for goods")
	buy.SetResultURL("https://your.domain/b2b/buygoods/result")
	buy.SetQueueTimeoutURL("https://your.domain/b2b/buygoods/timeout")

	// NOTE: Send() will perform a network call and requires valid credentials and network.
	// Uncomment the following lines when ready to perform real requests.
	/*
		resp, err := buy.Send()
		if err != nil {
			log.Fatalf("buy goods failed: %v", err)
		}
		fmt.Printf("B2B request response: %+v\n", resp)
	*/

	// Start webhook server for callbacks (ParseB2BCallback will be used inside handler)
	svc := buy // reuse service instance for parsing callbacks
	http.HandleFunc("/webhook/b2b", Services.B2PayBillCallbackHandler(&Services.BusinessToPayBillService{Config: cfg, Client: client}))

	fmt.Println("Example server listening on :8080 (webhook route /webhook/b2b)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ptr(s string) *string { return &s }
