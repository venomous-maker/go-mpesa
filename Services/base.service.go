package Services

import (
	"encoding/base64"
	"errors"
	_ "fmt"
	"github.com/venomous-maker/go-mpesa/Abstracts"
	"regexp"
	"strings"
	"time"
)

type BaseService struct {
	Config *Abstracts.MpesaConfig
	Client Abstracts.MpesaInterface
}

// NewBaseService creates a new base service
func NewBaseService(cfg *Abstracts.MpesaConfig, client Abstracts.MpesaInterface) *BaseService {
	return &BaseService{
		Config: cfg,
		Client: client,
	}
}

// GenerateTimestamp returns the current timestamp in "YmdHis" format
func (b *BaseService) GenerateTimestamp() string {
	return time.Now().Format("20060102150405")
}

// GeneratePassword creates a base64-encoded password using business code, passkey, and timestamp
func (b *BaseService) GeneratePassword() string {
	timestamp := b.GenerateTimestamp()
	plain := b.Config.BusinessCode + b.Config.PassKey + timestamp
	return base64.StdEncoding.EncodeToString([]byte(plain))
}

// CleanPhoneNumber formats and validates a phone number for M-Pesa API
func (b *BaseService) CleanPhoneNumber(phone, countryCode string) (string, error) {
	if strings.TrimSpace(phone) == "" {
		return "", errors.New("phone number cannot be empty")
	}

	phone = strings.TrimSpace(phone)

	if len(phone) < 9 {
		return "", errors.New("phone number is too short")
	}

	if strings.HasPrefix(phone, "+") {
		// Remove leading '+' and non-digit characters
		return regexp.MustCompile(`\D`).ReplaceAllString(phone[1:], ""), nil
	}

	if strings.HasPrefix(phone, "0") {
		// Replace leading 0 with country code
		return countryCode + regexp.MustCompile(`\D`).ReplaceAllString(phone[1:], ""), nil
	}

	// Clean all non-digit characters
	return regexp.MustCompile(`\D`).ReplaceAllString(phone, ""), nil
}
