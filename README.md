# 🇰🇪 go-mpesa – Safaricom M-Pesa Daraja API SDK for Golang

[![Go Reference](https://pkg.go.dev/badge/github.com/venomous-maker/go-mpesa.svg)](https://pkg.go.dev/github.com/venomous-maker/go-mpesa)

> A fully featured, type-safe, and extensible Go SDK for integrating with the Safaricom M-Pesa Daraja API (v1). Supports STK Push, C2B, B2C, Transaction Status, Reversals, and Account Balance queries.

---

## ✨ Features

- 🔐 Token caching and auto-refresh
- 📱 STK Push (Lipa na M-Pesa Online)
- 🔁 Reversals and Transaction Status queries
- 👤 Customer to Business (C2B)
- 🏢 Business to Customer (B2C)
- 💰 Account balance inquiries
- 🔧 Fluent service configuration
- 🧪 Simple test examples included

---

## 📦 Installation

```bash
go get github.com/venomous-maker/go-mpesa
```

> **Ensure Go 1.18+** is installed to support generics and modern features.

---

## 🛠️ Basic Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	mpesa "github.com/venomous-maker/go-mpesa/Mpesa"
	"github.com/venomous-maker/go-mpesa/Abstracts"
)

func main() {
	client := mpesa.NewMpesa("your-consumer-key", "your-consumer-secret", Abstracts.Sandbox).
		SetBusinessCode("123456").
		SetPassKey("your-passkey")

	stk := client.Stk()
	stk, err := stk.
		SetPhoneNumber("2547XXXXXXXX").
		SetAmount(10).
		SetTransactionType("CustomerPayBillOnline").
		SetCallbackUrl("https://your-callback.url").
		SetAccountReference("Ref123").
		SetTransactionDesc("Test Transaction").
		Push()

	if err != nil {
		log.Fatalf("STK Push failed: %v", err)
	}

	fmt.Println("✔️ STK Push initiated!")
	fmt.Printf("Response: %+v\n", stk.GetResponse())
}
```

---

## 🧩 Available Services

| Service                     | Description                         |
|----------------------------|-------------------------------------|
| `Stk()`                    | Initiate STK Push (Lipa na M-Pesa)  |
| `CustomerToBusiness()`     | Handle C2B payment validations      |
| `BusinessToCustomer()`     | Send B2C payouts to users           |
| `TransactionStatus()`      | Query transaction status            |
| `AccountBalance()`         | Check account balance               |
| `Reversal()`               | Reverse a completed transaction     |

---

## 🧪 Testing Locally

Run the sample test script:

```bash
go run Tests/stk.push.test.go
```

Customize the environment, credentials, and phone numbers before running.

---

## 🧠 Architecture

- `Abstracts/` — Interfaces, config structs, and reusable utilities (e.g., token management)
- `Services/` — M-Pesa service implementations (STK, C2B, B2C, etc.)
- `Mpesa/` — Fluent entrypoint for consumers
- `Tests/` — Sample test files to verify integration

---

## 🔐 Environments

| Environment | Value       |
|-------------|-------------|
| Sandbox     | `Abstracts.Sandbox` |
| Production  | `Abstracts.Production` |

---

## 📄 License

MIT © 2025 [@venomous-maker](https://github.com/venomous-maker)

---

## 🤝 Contributing

1. Fork the repo
2. Create your feature branch (`git checkout -b feature/new-service`)
3. Commit changes (`git commit -am 'Add new service'`)
4. Push to branch (`git push origin feature/new-service`)
5. Open a pull request 🚀

---

## 🧷 Disclaimer

This SDK is not affiliated with Safaricom or M-Pesa. It is a community-maintained tool for developers. Use it responsibly and ensure compliance with Safaricom's API terms and security policies.