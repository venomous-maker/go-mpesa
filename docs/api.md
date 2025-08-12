# API Reference

This document provides a comprehensive reference for the Go M-Pesa SDK API.

## Table of Contents

- [Configuration](#configuration)
- [Core Types](#core-types)
- [Services](#services)
- [Response Types](#response-types)
- [Error Types](#error-types)

## Configuration

### MpesaConfig

The main configuration struct for the M-Pesa SDK.

```go
type MpesaConfig struct {
    // Private fields - access via getters
}
```

#### Constructor

```go
func NewMpesaConfig(
    consumerKey, consumerSecret string,
    environment Environment,
    businessCode, passKey, securityCredential, queueTimeoutURL, resultURL *string,
) (*MpesaConfig, error)
```

#### Methods

| Method | Return Type | Description |
|--------|-------------|-------------|
| `GetConsumerKey()` | `string` | Returns the consumer key |
| `GetConsumerSecret()` | `string` | Returns the consumer secret |
| `GetEnvironment()` | `Environment` | Returns the environment setting |
| `GetBaseURL()` | `string` | Returns the API base URL |
| `GetBusinessCode()` | `string` | Returns the business shortcode |
| `GetPassKey()` | `string` | Returns the Lipa na M-Pesa passkey |
| `GetSecurityCredential()` | `string` | Returns the security credential |
| `GetQueueTimeoutURL()` | `string` | Returns the queue timeout URL |
| `GetResultURL()` | `string` | Returns the result URL |

### Environment

```go
type Environment string

const (
    Sandbox    Environment = "sandbox"
    Production Environment = "production"
)
```

## Core Types

### Mpesa Client

```go
type Mpesa struct {
    Config *Abstracts.MpesaConfig
    Client *Abstracts.ApiClient
}
```

#### Constructor

```go
func New(consumerKey, consumerSecret, environment string) (*Mpesa, error)
```

#### Service Factory Methods

| Method | Return Type | Description |
|--------|-------------|-------------|
| `STKPush()` | `*StkService` | Creates STK Push service |
| `B2C()` | `*B2CService` | Creates B2C service |
| `C2B()` | `*C2BService` | Creates C2B service |
| `AccountBalance()` | `*AccountBalanceService` | Creates Account Balance service |
| `TransactionStatus()` | `*TransactionStatusService` | Creates Transaction Status service |
| `Reversal()` | `*ReversalService` | Creates Reversal service |

## Services

### STK Push Service

```go
type StkService struct {
    *BaseService
    // Private fields
}
```

#### Methods

| Method | Parameters | Return Type | Description |
|--------|------------|-------------|-------------|
| `SetAmount(amount string)` | `amount`: Transaction amount | `*StkService` | Sets the transaction amount |
| `SetPhoneNumber(phone string)` | `phone`: Customer phone number | `*StkService` | Sets the customer phone number |
| `SetAccountReference(ref string)` | `ref`: Account reference | `*StkService` | Sets the account reference |
| `SetTransactionDesc(desc string)` | `desc`: Transaction description | `*StkService` | Sets the transaction description |
| `SetCallbackURL(url string)` | `url`: Callback URL | `*StkService` | Sets the callback URL |
| `SetTransactionType(txType string)` | `txType`: Transaction type | `*StkService` | Sets the transaction type |
| `Send()` | - | `(map[string]any, error)` | Initiates the STK Push |
| `Query(checkoutRequestID, merchantRequestID string)` | `checkoutRequestID`, `merchantRequestID`: Request IDs | `(map[string]any, error)` | Queries STK Push status |

#### Transaction Types

- `CustomerPayBillOnline` - Pay Bill transactions
- `CustomerBuyGoodsOnline` - Buy Goods transactions

### B2C Service

```go
type B2CService struct {
    *BaseService
    // Private fields
}
```

#### Methods

| Method | Parameters | Return Type | Description |
|--------|------------|-------------|-------------|
| `SetAmount(amount string)` | `amount`: Transaction amount | `*B2CService` | Sets the transaction amount |
| `SetPhoneNumber(phone string)` | `phone`: Recipient phone number | `*B2CService` | Sets the recipient phone number |
| `SetRemarks(remarks string)` | `remarks`: Transaction remarks | `*B2CService` | Sets the transaction remarks |
| `SetOccasion(occasion string)` | `occasion`: Transaction occasion | `*B2CService` | Sets the transaction occasion |
| `SetCommandID(commandID string)` | `commandID`: Command ID | `*B2CService` | Sets the command ID |
| `SetQueueTimeoutURL(url string)` | `url`: Queue timeout URL | `*B2CService` | Sets the queue timeout URL |
| `SetResultURL(url string)` | `url`: Result URL | `*B2CService` | Sets the result URL |
| `Send()` | - | `(map[string]any, error)` | Initiates the B2C transaction |

#### Command IDs

- `BusinessPayment` - General business payments
- `SalaryPayment` - Salary payments
- `PromotionPayment` - Promotional payments

### C2B Service

```go
type C2BService struct {
    *BaseService
    // Private fields
}
```

#### Methods

| Method | Parameters | Return Type | Description |
|--------|------------|-------------|-------------|
| `SetValidationURL(url string)` | `url`: Validation URL | `*C2BService` | Sets the validation URL |
| `SetConfirmationURL(url string)` | `url`: Confirmation URL | `*C2BService` | Sets the confirmation URL |
| `SetAmount(amount string)` | `amount`: Transaction amount | `*C2BService` | Sets the transaction amount |
| `SetPhoneNumber(phone string)` | `phone`: Customer phone number | `*C2BService` | Sets the customer phone number |
| `SetBillRefNumber(ref string)` | `ref`: Bill reference number | `*C2BService` | Sets the bill reference number |
| `SetCommandID(commandID string)` | `commandID`: Command ID | `*C2BService` | Sets the command ID |
| `RegisterURLs()` | - | `(map[string]any, error)` | Registers C2B URLs |
| `Simulate()` | - | `(map[string]any, error)` | Simulates C2B transaction (sandbox only) |

### Account Balance Service

```go
type AccountBalanceService struct {
    *BaseService
    // Private fields
}
```

#### Methods

| Method | Parameters | Return Type | Description |
|--------|------------|-------------|-------------|
| `SetRemarks(remarks string)` | `remarks`: Query remarks | `*AccountBalanceService` | Sets the query remarks |
| `SetQueueTimeoutURL(url string)` | `url`: Queue timeout URL | `*AccountBalanceService` | Sets the queue timeout URL |
| `SetResultURL(url string)` | `url`: Result URL | `*AccountBalanceService` | Sets the result URL |
| `Query()` | - | `(map[string]any, error)` | Queries the account balance |

### Transaction Status Service

```go
type TransactionStatusService struct {
    *BaseService
    // Private fields
}
```

#### Methods

| Method | Parameters | Return Type | Description |
|--------|------------|-------------|-------------|
| `SetTransactionID(id string)` | `id`: Transaction ID | `*TransactionStatusService` | Sets the transaction ID |
| `SetRemarks(remarks string)` | `remarks`: Query remarks | `*TransactionStatusService` | Sets the query remarks |
| `SetOccasion(occasion string)` | `occasion`: Query occasion | `*TransactionStatusService` | Sets the query occasion |
| `SetQueueTimeoutURL(url string)` | `url`: Queue timeout URL | `*TransactionStatusService` | Sets the queue timeout URL |
| `SetResultURL(url string)` | `url`: Result URL | `*TransactionStatusService` | Sets the result URL |
| `Query()` | - | `(map[string]any, error)` | Queries the transaction status |

### Reversal Service

```go
type ReversalService struct {
    *BaseService
    // Private fields
}
```

#### Methods

| Method | Parameters | Return Type | Description |
|--------|------------|-------------|-------------|
| `SetTransactionID(id string)` | `id`: Transaction ID | `*ReversalService` | Sets the transaction ID |
| `SetAmount(amount string)` | `amount`: Reversal amount | `*ReversalService` | Sets the reversal amount |
| `SetRemarks(remarks string)` | `remarks`: Reversal remarks | `*ReversalService` | Sets the reversal remarks |
| `SetOccasion(occasion string)` | `occasion`: Reversal occasion | `*ReversalService` | Sets the reversal occasion |
| `SetQueueTimeoutURL(url string)` | `url`: Queue timeout URL | `*ReversalService` | Sets the queue timeout URL |
| `SetResultURL(url string)` | `url`: Result URL | `*ReversalService` | Sets the result URL |
| `Reverse()` | - | `(map[string]any, error)` | Initiates the reversal |

## Response Types

### Common Response Fields

All M-Pesa API responses contain these common fields:

```go
type BaseResponse struct {
    ResponseCode        string `json:"ResponseCode"`
    ResponseDescription string `json:"ResponseDescription"`
    ConversationID      string `json:"ConversationID"`
    OriginatorConversationID string `json:"OriginatorConversationID"`
}
```

### STK Push Response

```go
type STKResponse struct {
    MerchantRequestID   string `json:"MerchantRequestID"`
    CheckoutRequestID   string `json:"CheckoutRequestID"`
    ResponseCode        string `json:"ResponseCode"`
    ResponseDescription string `json:"ResponseDescription"`
    CustomerMessage     string `json:"CustomerMessage"`
}
```

### STK Push Query Response

```go
type STKQueryResponse struct {
    ResponseCode        string `json:"ResponseCode"`
    ResponseDescription string `json:"ResponseDescription"`
    MerchantRequestID   string `json:"MerchantRequestID"`
    CheckoutRequestID   string `json:"CheckoutRequestID"`
    ResultCode          string `json:"ResultCode"`
    ResultDesc          string `json:"ResultDesc"`
}
```

### Response Codes

| Code | Description |
|------|-------------|
| `0` | Success |
| `1` | Insufficient funds |
| `2` | Less than minimum transaction value |
| `3` | More than maximum transaction value |
| `4` | Would exceed daily transfer limit |
| `5` | Would exceed minimum balance |
| `11` | Invalid phone number |
| `12` | Invalid account number |
| `13` | Invalid amount |
| `17` | System internal error |

## Error Types

### Custom Error Types

```go
type MpesaError struct {
    Code    string
    Message string
    Details map[string]interface{}
}

func (e *MpesaError) Error() string {
    return fmt.Sprintf("M-Pesa Error [%s]: %s", e.Code, e.Message)
}
```

### Common Error Scenarios

| Error Type | Description | Handling |
|------------|-------------|----------|
| Authentication Error | Invalid credentials | Check consumer key/secret |
| Validation Error | Invalid request parameters | Validate input data |
| Network Error | Connection issues | Implement retry logic |
| Rate Limit Error | Too many requests | Implement backoff strategy |
| Server Error | M-Pesa service unavailable | Retry with exponential backoff |

## Usage Examples

### Complete STK Push Flow

```go
// Initialize client
mpesa, err := Mpesa.New("key", "secret", "sandbox")
if err != nil {
    log.Fatal(err)
}

// Create STK service
stkService := mpesa.STKPush()

// Set parameters and send
response, err := stkService.
    SetAmount("100").
    SetPhoneNumber("254712345678").
    SetAccountReference("ORDER001").
    SetTransactionDesc("Payment for order").
    SetCallbackURL("https://example.com/callback").
    Send()

if err != nil {
    log.Fatal(err)
}

// Extract important response fields
checkoutRequestID := response["CheckoutRequestID"].(string)
merchantRequestID := response["MerchantRequestID"].(string)

// Query status
status, err := stkService.Query(checkoutRequestID, merchantRequestID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %+v\n", status)
```

### Error Handling Pattern

```go
response, err := service.Send()
if err != nil {
    if mpesaErr, ok := err.(*MpesaError); ok {
        switch mpesaErr.Code {
        case "1":
            log.Println("Insufficient funds")
        case "11":
            log.Println("Invalid phone number")
        default:
            log.Printf("M-Pesa error: %s", mpesaErr.Message)
        }
    } else {
        log.Printf("Network error: %v", err)
    }
    return
}

// Check response code
if responseCode := response["ResponseCode"].(string); responseCode != "0" {
    log.Printf("Transaction failed: %s", response["ResponseDescription"])
}
```

## Rate Limits

The M-Pesa API has the following rate limits:

- **Sandbox**: 10 requests per second
- **Production**: Varies by business agreement

Implement appropriate rate limiting in your application:

```go
import "golang.org/x/time/rate"

// Create rate limiter (10 requests per second)
limiter := rate.NewLimiter(rate.Limit(10), 1)

// Use before making requests
if err := limiter.Wait(context.Background()); err != nil {
    log.Fatal(err)
}

// Make your M-Pesa request
response, err := service.Send()
```
