package main

import (
	"fmt"
	"log"
	"time"

	"github.com/venomous-maker/go-mpesa/Abstracts"
	"github.com/venomous-maker/go-mpesa/Services"
)

func main() {
	// Initialize config
	cfg, err := Abstracts.NewMpesaConfig(
		"rHZXmBkGz6Ne30cA923bp9G0rSAK41hsDVCq65x522WkVqCF",                 // consumer key
		"QC7BEvNXH9FfMATpduK1fTh1836XisZ9qG7cIZ15S9cDGzIsBMc2YAkAKsEr7wjo", // consumer secret
		Abstracts.Sandbox, // environment
		ptr("174379"),     // business short code
		ptr("bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919"), // passkey
		nil, nil, nil, // optional: securityCredential, queueTimeoutURL, resultURL
	)
	if err != nil {
		log.Fatalf("config init error: %v", err)
	}

	// Create STK service
	stk := Services.NewStkService(cfg, Abstracts.NewApiClient(cfg))

	// Set phone number
	stk, err = stk.SetPhoneNumber("254714544312")
	if err != nil {
		log.Fatalf("invalid phone: %v", err)
	}

	// Configure request
	stk.
		SetAmount("10").
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
	time.Sleep(25 * time.Second)
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
