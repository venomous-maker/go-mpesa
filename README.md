# Go M-Pesa SDK

A comprehensive Go SDK for integrating with Safaricom's M-Pesa API. This package provides a clean, type-safe interface for M-Pesa services including STK Push, B2C, C2B, Account Balance, Transaction Status, and Reversals.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-Passing-green.svg)](Tests/)
[![Documentation](https://img.shields.io/badge/Documentation-Complete-brightgreen.svg)](docs/)

## Features

- 🚀 **STK Push (Lipa na M-Pesa Online)** - Initiate payments from customer's phone
- 💸 **B2C (Business to Customer)** - Send money to customers
- 💰 **C2B (Customer to Business)** - Receive payments from customers
- 📊 **Account Balance** - Check account balance
- 🔍 **Transaction Status** - Query transaction status
- ↩️ **Reversal** - Reverse transactions
- 🔒 **Secure Authentication** - Automatic token management
- 🧪 **Sandbox & Production** - Support for both environments
- ✅ **Comprehensive Testing** - Full test coverage with mocks
- 📝 **Type Safety** - Fully typed API responses
- 🔄 **Automatic Retries** - Built-in retry mechanism for failed requests
- 📱 **Phone Number Validation** - Automatic phone number formatting
- 🛡️ **Security** - Encrypted credentials and secure token handling

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Services](#services)
  - [STK Push](#stk-push-lipa-na-m-pesa-online)
  - [B2C Transactions](#b2c-business-to-customer)
  - [C2B Transactions](#c2b-customer-to-business)
  - [Account Balance](#account-balance)
  - [Transaction Status](#transaction-status)
  - [Transaction Reversal](#transaction-reversal)
- [Error Handling](#error-handling)
- [Testing](#testing)
- [API Reference](#api-reference)
- [Contributing](#contributing)
- [License](#license)

## Installation

```bash
go get github.com/venomous-maker/go-mpesa
```

## Quick Start

### 1. Initialize the M-Pesa Client

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/venomous-maker/go-mpesa/Mpesa"
)

func main() {
    // Initialize M-Pesa client
    mpesa, err := Mpesa.New(
        "your_consumer_key",
        "your_consumer_secret", 
        "sandbox", // or "production"
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Your M-Pesa client is ready to use
    fmt.Println("M-Pesa client initialized successfully")
}
```

### 2. Simple STK Push Example

```go
// Configure STK Push service
stkService := mpesa.STKPush()

// Set transaction details
response, err := stkService.
    SetAmount("100").
    SetPhoneNumber("254712345678").
    SetAccountReference("ORDER001").
    SetTransactionDesc("Payment for order").
    SetCallbackURL("https://yourdomain.com/callback").
    Send()

if err != nil {
    log.Fatal(err)
}

fmt.Printf("STK Push initiated: %+v\n", response)
```

## Configuration

### Basic Configuration

```go
mpesa, err := Mpesa.New(
    "your_consumer_key",
    "your_consumer_secret",
    "sandbox", // Environment: "sandbox" or "production"
)
```

### Advanced Configuration with Optional Parameters

```go
import "github.com/venomous-maker/go-mpesa/Abstracts"

// Create configuration with optional parameters
businessCode := "174379"
passKey := "your_lipa_na_mpesa_passkey"
securityCredential := "your_security_credential"
queueTimeoutURL := "https://yourdomain.com/timeout"
resultURL := "https://yourdomain.com/result"

cfg, err := Abstracts.NewMpesaConfig(
    "consumer_key",
    "consumer_secret",
    Abstracts.Sandbox,
    &businessCode,
    &passKey,
    &securityCredential,
    &queueTimeoutURL,
    &resultURL,
)

if err != nil {
    log.Fatal(err)
}

// Create M-Pesa client with custom config
client := Abstracts.NewApiClient(cfg)
mpesa := &Mpesa.Mpesa{
    Config: cfg,
    Client: client,
}
```

## Services

### STK Push (Lipa na M-Pesa Online)

STK Push allows you to initiate M-Pesa payments directly from a customer's mobile phone.

#### Basic STK Push

```go
stkService := mpesa.STKPush()

response, err := stkService.
    SetAmount("1000").
    SetPhoneNumber("254712345678").
    SetAccountReference("INV001").
    SetTransactionDesc("Payment for invoice INV001").
    SetCallbackURL("https://yourdomain.com/stk-callback").
    Send()

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Checkout Request ID: %s\n", response["CheckoutRequestID"])
fmt.Printf("Merchant Request ID: %s\n", response["MerchantRequestID"])
```

#### STK Push with Custom Transaction Type

```go
response, err := stkService.
    SetTransactionType("CustomerBuyGoodsOnline").
    SetAmount("500").
    SetPhoneNumber("254712345678").
    SetAccountReference("PROD123").
    SetTransactionDesc("Purchase of Product 123").
    SetCallbackURL("https://yourdomain.com/stk-callback").
    Send()
```

#### Query STK Push Status

```go
checkoutRequestID := "ws_CO_123456789"
merchantRequestID := "12345-67890-1"

status, err := stkService.Query(checkoutRequestID, merchantRequestID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Transaction Status: %+v\n", status)
```

### B2C (Business to Customer)

Send money from your business account to customer accounts.

```go
b2cService := mpesa.B2C()

response, err := b2cService.
    SetAmount("1000").
    SetPhoneNumber("254712345678").
    SetRemarks("Salary payment").
    SetOccasion("Monthly salary").
    SetCommandID("BusinessPayment"). // or "SalaryPayment", "PromotionPayment"
    Send()

if err != nil {
    log.Fatal(err)
}

fmt.Printf("B2C Transaction ID: %s\n", response["ConversationID"])
```

#### Available B2C Command IDs

- `BusinessPayment` - General business payments
- `SalaryPayment` - Salary payments to employees
- `PromotionPayment` - Promotional payments and rewards

### C2B (Customer to Business)

Register URLs and simulate customer payments to your business.

#### Register C2B URLs

```go
c2bService := mpesa.C2B()

response, err := c2bService.
    SetValidationURL("https://yourdomain.com/c2b-validation").
    SetConfirmationURL("https://yourdomain.com/c2b-confirmation").
    RegisterURLs()

if err != nil {
    log.Fatal(err)
}

fmt.Printf("URLs registered: %+v\n", response)
```

#### Simulate C2B Transaction (Sandbox Only)

```go
response, err := c2bService.
    SetAmount("1000").
    SetPhoneNumber("254712345678").
    SetBillRefNumber("INV001").
    SetCommandID("CustomerPayBillOnline"). // or "CustomerBuyGoodsOnline"
    Simulate()

if err != nil {
    log.Fatal(err)
}

fmt.Printf("C2B Simulation: %+v\n", response)
```

### Account Balance

Query your M-Pesa account balance.

```go
balanceService := mpesa.AccountBalance()

response, err := balanceService.
    SetRemarks("Account balance inquiry").
    SetQueueTimeoutURL("https://yourdomain.com/timeout").
    SetResultURL("https://yourdomain.com/balance-result").
    Query()

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Balance Query ID: %s\n", response["ConversationID"])
```

### Transaction Status

Check the status of any M-Pesa transaction.

```go
statusService := mpesa.TransactionStatus()

response, err := statusService.
    SetTransactionID("LHG31AA5TX").
    SetRemarks("Transaction status inquiry").
    SetOccasion("Status check").
    SetQueueTimeoutURL("https://yourdomain.com/timeout").
    SetResultURL("https://yourdomain.com/status-result").
    Query()

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status Query ID: %s\n", response["ConversationID"])
```

### Transaction Reversal

Reverse a completed M-Pesa transaction.

```go
reversalService := mpesa.Reversal()

response, err := reversalService.
    SetTransactionID("LHG31AA5TX").
    SetAmount("1000").
    SetRemarks("Reversal for duplicate payment").
    SetOccasion("Duplicate transaction").
    SetQueueTimeoutURL("https://yourdomain.com/timeout").
    SetResultURL("https://yourdomain.com/reversal-result").
    Reverse()

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Reversal ID: %s\n", response["ConversationID"])
```

## Error Handling

The SDK provides comprehensive error handling with detailed error messages.

```go
response, err := stkService.
    SetAmount("100").
    SetPhoneNumber("254712345678").
    Send()

if err != nil {
    // Handle different types of errors
    switch {
    case strings.Contains(err.Error(), "authentication"):
        log.Println("Authentication failed - check credentials")
    case strings.Contains(err.Error(), "insufficient funds"):
        log.Println("Insufficient funds in account")
    case strings.Contains(err.Error(), "invalid phone"):
        log.Println("Invalid phone number format")
    default:
        log.Printf("Transaction failed: %v", err)
    }
    return
}

// Check response for transaction status
if responseCode, ok := response["ResponseCode"].(string); ok {
    if responseCode == "0" {
        log.Println("Transaction initiated successfully")
    } else {
        log.Printf("Transaction failed with code: %s", responseCode)
    }
}
```

## Testing

The SDK includes comprehensive tests with mocks for all services.

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestSTKPush ./Tests/
```

### Example Test

```go
func TestSTKPush(t *testing.T) {
    // Create mock client
    mockClient := &MockMpesaClient{}
    
    // Set up expected response
    expectedResponse := map[string]any{
        "CheckoutRequestID": "ws_CO_123456789",
        "ResponseCode":      "0",
        "ResponseDescription": "Success",
    }
    
    mockClient.On("Post", mock.Anything, mock.Anything).Return(expectedResponse, nil)
    
    // Create service with mock
    cfg := createTestConfig()
    stkService := NewStkService(cfg, mockClient)
    
    // Execute test
    response, err := stkService.
        SetAmount("100").
        SetPhoneNumber("254712345678").
        Send()
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "ws_CO_123456789", response["CheckoutRequestID"])
}
```

## Webhook Handling

Handle M-Pesa callbacks in your application:

### STK Push Callback

```go
func handleSTKCallback(w http.ResponseWriter, r *http.Request) {
    var callback STKCallbackResponse
    
    if err := json.NewDecoder(r.Body).Decode(&callback); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Process the callback
    if callback.Body.StkCallback.ResultCode == 0 {
        // Payment successful
        log.Printf("Payment successful: %s", 
            callback.Body.StkCallback.CheckoutRequestID)
    } else {
        // Payment failed
        log.Printf("Payment failed: %s", 
            callback.Body.StkCallback.ResultDesc)
    }
    
    // Send response
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
```

### B2C Result Callback

```go
func handleB2CResult(w http.ResponseWriter, r *http.Request) {
    var result B2CResultResponse
    
    if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Process the result
    if result.Result.ResultCode == 0 {
        // Transaction successful
        log.Printf("B2C successful: %s", result.Result.ConversationID)
    } else {
        // Transaction failed
        log.Printf("B2C failed: %s", result.Result.ResultDesc)
    }
    
    w.WriteHeader(http.StatusOK)
}
```

## Best Practices

### 1. Environment Management

```go
// Use environment variables for credentials
import "os"

consumerKey := os.Getenv("MPESA_CONSUMER_KEY")
consumerSecret := os.Getenv("MPESA_CONSUMER_SECRET")
environment := os.Getenv("MPESA_ENVIRONMENT") // "sandbox" or "production"

mpesa, err := Mpesa.New(consumerKey, consumerSecret, environment)
```

### 2. Phone Number Formatting

```go
// The SDK automatically formats phone numbers, but ensure they start with 254
phoneNumber := "0712345678"   // Will be converted to 254712345678
phoneNumber := "712345678"    // Will be converted to 254712345678
phoneNumber := "254712345678" // Already in correct format
```

### 3. Callback URL Security

```go
// Always use HTTPS for callback URLs
callbackURL := "https://yourdomain.com/mpesa-callback"

// Implement signature verification for callbacks
func verifyCallback(r *http.Request) bool {
    // Implement signature verification logic
    return true
}
```

### 4. Error Logging

```go
import "github.com/sirupsen/logrus"

logger := logrus.New()

response, err := stkService.Send()
if err != nil {
    logger.WithFields(logrus.Fields{
        "service": "STK Push",
        "error":   err.Error(),
        "amount":  stkService.amount,
        "phone":   stkService.phoneNumber,
    }).Error("Transaction failed")
}
```

## API Reference

### Core Types

```go
// Environment constants
const (
    Sandbox    Environment = "sandbox"
    Production Environment = "production"
)

// Common response fields
type MpesaResponse struct {
    ResponseCode        string `json:"ResponseCode"`
    ResponseDescription string `json:"ResponseDescription"`
    ConversationID      string `json:"ConversationID"`
    OriginatorConversationID string `json:"OriginatorConversationID"`
}

// STK Push specific response
type STKResponse struct {
    MerchantRequestID   string `json:"MerchantRequestID"`
    CheckoutRequestID   string `json:"CheckoutRequestID"`
    ResponseCode        string `json:"ResponseCode"`
    ResponseDescription string `json:"ResponseDescription"`
    CustomerMessage     string `json:"CustomerMessage"`
}
```

### Service Methods

#### STK Service
- `SetAmount(amount string) *StkService`
- `SetPhoneNumber(phone string) *StkService`
- `SetAccountReference(ref string) *StkService`
- `SetTransactionDesc(desc string) *StkService`
- `SetCallbackURL(url string) *StkService`
- `SetTransactionType(txType string) *StkService`
- `Send() (map[string]any, error)`
- `Query(checkoutRequestID, merchantRequestID string) (map[string]any, error)`

#### B2C Service
- `SetAmount(amount string) *B2CService`
- `SetPhoneNumber(phone string) *B2CService`
- `SetRemarks(remarks string) *B2CService`
- `SetOccasion(occasion string) *B2CService`
- `SetCommandID(commandID string) *B2CService`
- `Send() (map[string]any, error)`

#### Account Balance Service
- `SetRemarks(remarks string) *AccountBalanceService`
- `SetQueueTimeoutURL(url string) *AccountBalanceService`
- `SetResultURL(url string) *AccountBalanceService`
- `Query() (map[string]any, error)`

For complete API documentation, see the [API Documentation](docs/api.md).

## Examples

Check out the [examples directory](examples/) for complete working examples:

- [Basic STK Push](examples/stk-push/main.go)
- [B2C Payment](examples/b2c/main.go)
- [Web Application Integration](examples/webapp/main.go)
- [Webhook Handling](examples/webhooks/main.go)

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/venomous-maker/go-mpesa.git
   cd go-mpesa
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   go test ./...
   ```

4. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes and version history.

## Support

- 📧 Email: support@venomous-maker.com
- 🐛 Issues: [GitHub Issues](https://github.com/venomous-maker/go-mpesa/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/venomous-maker/go-mpesa/discussions)
- 📖 Documentation: [docs/](docs/)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This SDK is not officially affiliated with Safaricom. M-Pesa is a trademark of Safaricom Limited. Use this SDK in accordance with Safaricom's terms of service and your agreement with them.

---

**Made with ❤️ by [Venomous Maker](https://github.com/venomous-maker)**
