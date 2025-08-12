# Go M-Pesa SDK

A comprehensive Go SDK for integrating with Safaricom's M-Pesa API. This package provides a clean, type-safe interface for M-Pesa services including STK Push, B2C, C2B, Account Balance, Transaction Status, and Reversals.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-Passing-green.svg)](Tests/)

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
    
    "github.com/venomous-maker/go-mpesa/Abstracts"
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
    
    // Set your business shortcode
    mpesa.SetBusinessCode("174379")
    mpesa.SetPassKey("your_passkey")
}
```

### 2. STK Push (Lipa na M-Pesa Online)

```go
// Create STK service
stkService := mpesa.STK()

// Configure the payment request
stkService.
    SetAmount(100).
    SetTransactionType("CustomerPayBillOnline").
    SetCallbackUrl("https://yourdomain.com/callback").
    SetAccountReference("ORDER123").
    SetTransactionDesc("Payment for Order #123")

// Set phone number (supports multiple formats)
stkService, err = stkService.SetPhoneNumber("254711223344")
if err != nil {
    log.Fatal(err)
}

// Initiate the STK push
response, err := stkService.Push()
if err != nil {
    log.Fatal(err)
}

fmt.Println("STK Push initiated successfully!")

// Get the checkout request ID for tracking
checkoutID, err := stkService.GetCheckoutRequestID()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Checkout Request ID: %s\n", checkoutID)
```

### 3. Query STK Push Status

```go
// Query the status of an STK push
status, err := stkService.Query(checkoutID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Transaction Status: %+v\n", status)
```

## Advanced Configuration

### Custom Configuration

```go
import "github.com/venomous-maker/go-mpesa/Abstracts"

// Create custom configuration
cfg, err := Abstracts.NewMpesaConfig(
    "consumer_key",
    "consumer_secret",
    Abstracts.Sandbox, // or Abstracts.Production
    &businessCode,     // Business shortcode
    &passKey,          // Lipa na M-Pesa passkey
    &securityCredential, // For B2C, Reversal, etc.
    &queueTimeoutURL,  // Queue timeout URL
    &resultURL,        // Result URL
)
if err != nil {
    log.Fatal(err)
}

// Create client with custom config
client := Abstracts.NewApiClient(cfg)
```

### Environment Configuration

```go
// Sandbox environment (for testing)
mpesa, err := Mpesa.New(consumerKey, consumerSecret, "sandbox")

// Production environment
mpesa, err := Mpesa.New(consumerKey, consumerSecret, "production")
```

## Services

### STK Push Service

STK Push allows you to initiate M-Pesa payments from a customer's phone.

```go
stkService := mpesa.STK()

// Method chaining for clean configuration
response, err := stkService.
    SetAmount(1000).
    SetTransactionType("CustomerPayBillOnline").
    SetPhoneNumber("254711223344").
    SetCallbackUrl("https://yourdomain.com/callback").
    SetAccountReference("INV123").
    SetTransactionDesc("Invoice payment").
    Push()
```

**Supported Amount Types:**
- `int`: `SetAmount(100)`
- `string`: `SetAmount("100")`
- `int64`: `SetAmount(int64(100))`
- `float`: `SetAmount(99.99)`

**Phone Number Formats:**
- `254711223344` (with country code)
- `0711223344` (without country code)
- `+254711223344` (international format)

### B2C Service

Send money from business to customer.

```go
b2cService := mpesa.B2C()

response, err := b2cService.
    SetAmount(1000).
    SetPhoneNumber("254711223344").
    SetCommandID("BusinessPayment").
    SetRemarks("Salary payment").
    SetOccasion("Monthly salary").
    Send()
```

### C2B Service

Register URLs and simulate C2B transactions.

```go
c2bService := mpesa.C2B()

// Register URLs
err := c2bService.
    SetValidationURL("https://yourdomain.com/validation").
    SetConfirmationURL("https://yourdomain.com/confirmation").
    RegisterURLs()

// Simulate payment
response, err := c2bService.
    SetAmount(1000).
    SetPhoneNumber("254711223344").
    SetBillRefNumber("REF123").
    Simulate()
```

### Account Balance

Check your M-Pesa account balance.

```go
balanceService := mpesa.AccountBalance()

balance, err := balanceService.
    SetCommandID("AccountBalance").
    SetRemarks("Balance inquiry").
    Query()
```

### Transaction Status

Query the status of any M-Pesa transaction.

```go
statusService := mpesa.TransactionStatus()

status, err := statusService.
    SetTransactionID("ABC123XYZ").
    SetCommandID("TransactionStatusQuery").
    SetRemarks("Status check").
    Query()
```

### Reversal

Reverse a completed M-Pesa transaction.

```go
reversalService := mpesa.Reversal()

response, err := reversalService.
    SetTransactionID("ABC123XYZ").
    SetAmount(1000).
    SetCommandID("TransactionReversal").
    SetRemarks("Refund processing").
    SetOccasion("Customer refund").
    Reverse()
```

## Error Handling

The SDK provides detailed error messages for common issues:

```go
stkService := mpesa.STK()

// This will return a validation error
_, err := stkService.Push()
if err != nil {
    switch {
    case strings.Contains(err.Error(), "amount is required"):
        fmt.Println("Please set the amount")
    case strings.Contains(err.Error(), "phone number is required"):
        fmt.Println("Please set the phone number")
    case strings.Contains(err.Error(), "callback URL is required"):
        fmt.Println("Please set the callback URL")
    default:
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Testing

The package includes comprehensive tests with mocking support.

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test package
go test ./Tests

# Run with verbose output
go test -v ./Tests
```

### Writing Tests

The package provides mock interfaces for testing:

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestSTKPush(t *testing.T) {
    // Create mock client
    mockClient := &MockMpesaInterface{}
    
    // Setup expected response
    expectedResponse := map[string]any{
        "CheckoutRequestID": "ws_CO_12345678",
        "ResponseCode": "0",
    }
    
    mockClient.On("ExecuteRequest", mock.Anything, mock.Anything).
        Return(expectedResponse, nil)
    
    // Test your service
    cfg := createTestConfig()
    service := Services.NewStkService(cfg, mockClient)
    
    // Configure and test
    service.SetAmount(100)
    // ... configure other fields
    
    response, err := service.Push()
    assert.NoError(t, err)
    assert.NotNil(t, response)
}
```

## Callback Handling

### STK Push Callback

Handle STK Push callbacks in your application:

```go
type STKCallback struct {
    Body struct {
        StkCallback struct {
            MerchantRequestID   string `json:"MerchantRequestID"`
            CheckoutRequestID   string `json:"CheckoutRequestID"`
            ResultCode          int    `json:"ResultCode"`
            ResultDesc          string `json:"ResultDesc"`
            CallbackMetadata    struct {
                Item []struct {
                    Name  string      `json:"Name"`
                    Value interface{} `json:"Value"`
                } `json:"Item"`
            } `json:"CallbackMetadata,omitempty"`
        } `json:"stkCallback"`
    } `json:"Body"`
}

func handleSTKCallback(w http.ResponseWriter, r *http.Request) {
    var callback STKCallback
    
    if err := json.NewDecoder(r.Body).Decode(&callback); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    stkCallback := callback.Body.StkCallback
    
    if stkCallback.ResultCode == 0 {
        // Payment successful
        fmt.Printf("Payment successful: %s\n", stkCallback.CheckoutRequestID)
        
        // Extract payment details from CallbackMetadata
        for _, item := range stkCallback.CallbackMetadata.Item {
            switch item.Name {
            case "Amount":
                fmt.Printf("Amount: %v\n", item.Value)
            case "MpesaReceiptNumber":
                fmt.Printf("Receipt: %v\n", item.Value)
            case "PhoneNumber":
                fmt.Printf("Phone: %v\n", item.Value)
            }
        }
    } else {
        // Payment failed
        fmt.Printf("Payment failed: %s\n", stkCallback.ResultDesc)
    }
    
    w.WriteHeader(http.StatusOK)
}
```

## Configuration Reference

### Environment Variables

You can use environment variables for configuration:

```bash
export MPESA_CONSUMER_KEY="your_consumer_key"
export MPESA_CONSUMER_SECRET="your_consumer_secret"
export MPESA_ENVIRONMENT="sandbox"
export MPESA_BUSINESS_CODE="174379"
export MPESA_PASSKEY="your_passkey"
```

```go
import "os"

mpesa, err := Mpesa.New(
    os.Getenv("MPESA_CONSUMER_KEY"),
    os.Getenv("MPESA_CONSUMER_SECRET"),
    os.Getenv("MPESA_ENVIRONMENT"),
)
```

### API Endpoints

The SDK automatically handles API endpoints based on environment:

**Sandbox:**
- Base URL: `https://sandbox.safaricom.co.ke`

**Production:**
- Base URL: `https://api.safaricom.co.ke`

## Common Use Cases

### E-commerce Checkout

```go
func processPayment(orderID string, amount int, phoneNumber string) error {
    mpesa, err := Mpesa.New(consumerKey, consumerSecret, "production")
    if err != nil {
        return err
    }
    
    mpesa.SetBusinessCode(businessCode)
    mpesa.SetPassKey(passKey)
    
    stkService := mpesa.STK()
    
    _, err = stkService.
        SetAmount(amount).
        SetTransactionType("CustomerPayBillOnline").
        SetPhoneNumber(phoneNumber).
        SetCallbackUrl("https://yoursite.com/mpesa/callback").
        SetAccountReference(orderID).
        SetTransactionDesc(fmt.Sprintf("Payment for order %s", orderID)).
        Push()
    
    return err
}
```

### Salary Payments

```go
func paySalary(employeePhone string, amount int) error {
    mpesa, err := Mpesa.New(consumerKey, consumerSecret, "production")
    if err != nil {
        return err
    }
    
    b2cService := mpesa.B2C()
    
    _, err = b2cService.
        SetAmount(amount).
        SetPhoneNumber(employeePhone).
        SetCommandID("SalaryPayment").
        SetRemarks("Monthly salary payment").
        SetOccasion("Salary").
        Send()
    
    return err
}
```

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make your changes and add tests
4. Ensure tests pass: `go test ./...`
5. Commit your changes: `git commit -m "Add new feature"`
6. Push to the branch: `git push origin feature/new-feature`
7. Submit a pull request

### Development Setup

```bash
# Clone the repository
git clone https://github.com/venomous-maker/go-mpesa.git
cd go-mpesa

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📧 **Email**: support@example.com
- 🐛 **Issues**: [GitHub Issues](https://github.com/venomous-maker/go-mpesa/issues)
- 📖 **Documentation**: [API Documentation](https://developer.safaricom.co.ke/)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/venomous-maker/go-mpesa/discussions)

## Changelog

### v1.0.0
- Initial release
- STK Push implementation
- B2C, C2B services
- Account Balance and Transaction Status
- Reversal service
- Comprehensive testing
- Full documentation

## Acknowledgments

- [Safaricom](https://safaricom.co.ke) for providing the M-Pesa API
- [Go](https://golang.org) team for the excellent programming language
- All contributors who help improve this package

---

**Disclaimer**: This is an unofficial SDK. Please refer to the official Safaricom M-Pesa API documentation for the most up-to-date information.
