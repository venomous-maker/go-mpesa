package main

import (
	"fmt"
	"log"

	"github.com/venomous-maker/go-mpesa/Abstracts"
	"github.com/venomous-maker/go-mpesa/Services"
)

func main() {
	// Initialize config
	cfg, err := Abstracts.NewMpesaConfig(
		"fEiZlnAbILFe9A0oKsODFEB8eAKRz9f4nb19YCGX8fsJbXgn",                 // consumer key
		"PckiOm357DQ4A3XzcPUGhgOHkU1ar6AgBSyGtU1QDx91JaMP3CQJSeAwYvSUBydi", // consumer secret
		Abstracts.Sandbox, // environment
		ptr("4148097"),    // business short code
		ptr("edad729a3b68c38fc282b2cdaa5196a3e0939d48c8cdc0f8a9a155bb750d8418"), // passkey
		nil, nil, nil, // optional: securityCredential, queueTimeoutURL, resultURL
	)
	if err != nil {
		log.Fatalf("config init error: %v", err)
	}

	// Create STK service
	stk := Services.NewStkService(cfg, Abstracts.NewApiClient(cfg))

	// Set phone number
	stk, err = stk.SetPhoneNumber("254111844429")
	if err != nil {
		log.Fatalf("invalid phone: %v", err)
	}

	// Configure request
	stk.
		SetAmount("1").
		SetTransactionType("CustomerPayBillOnline").
		SetCallbackUrl("https://5d007d77dea3.ngrok-free.app/callback").
		SetAccountReference("This is a test account reference").
		SetTransactionDesc("Test Push Mpesa")

	// Send STK Push
	_, err = stk.Push()
	if err != nil {
		log.Fatalf("stk push failed: %v", err)
	}

	fmt.Println("✔️ STK Push sent successfully.")
	fmt.Printf("Response:\n%+v\n", stk.GetResponse())

	// Query transaction
	id, err := stk.GetCheckoutRequestID()
	if err != nil {
		log.Fatalf("checkout id error: %v", err)
	}

	status, err := stk.Query(id)
	if err != nil {
		log.Fatalf("query error: %v", err)
	}

	fmt.Println("\n📦 STK Push Query Result:")
	fmt.Printf("%+v\n", status)
}

func ptr(s string) *string {
	return &s
}
