// Package Abstracts provides core abstractions, configurations, and interfaces for the M-Pesa SDK.
// This package contains the fundamental types and utilities needed for M-Pesa API integration.
package Abstracts

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

// Environment represents the M-Pesa API environment (sandbox or production).
type Environment string

const (
	// Sandbox represents the M-Pesa sandbox environment for testing and development.
	// Use this environment during development and testing phases.
	Sandbox Environment = "sandbox"

	// Production represents the live M-Pesa production environment.
	// Use this environment for real transactions with actual money.
	Production Environment = "live"
)

// MpesaConfig holds all configuration settings required for M-Pesa API operations.
// This includes credentials, environment settings, URLs, and security parameters.
type MpesaConfig struct {
	consumerKey        string      // Consumer key from Safaricom Developer Portal
	consumerSecret     string      // Consumer secret from Safaricom Developer Portal
	environment        Environment // Target environment (sandbox or production)
	baseURL            string      // Base URL for M-Pesa API endpoints
	businessCode       string      // Business shortcode for transactions
	passKey            string      // Lipa na M-Pesa Online passkey
	securityCredential string      // Security credential for B2C and other operations
	queueTimeoutURL    string      // URL for queue timeout notifications
	resultURL          string      // URL for transaction result notifications
}

// NewMpesaConfig creates a new M-Pesa configuration with the provided parameters.
// This is the primary constructor for creating M-Pesa configuration instances.
// Optional parameters can be provided as pointers, allowing nil values for unused settings.
//
// Parameters:
//   - consumerKey: Consumer key obtained from Safaricom Developer Portal
//   - consumerSecret: Consumer secret obtained from Safaricom Developer Portal
//   - environment: Target environment (Sandbox or Production)
//   - businessCode: Optional business shortcode for transactions
//   - passKey: Optional Lipa na M-Pesa Online passkey for STK Push
//   - securityCredential: Optional security credential for B2C operations
//   - queueTimeoutURL: Optional URL for queue timeout notifications
//   - resultURL: Optional URL for transaction result notifications
//
// Returns:
//   - *MpesaConfig: A configured M-Pesa configuration instance
//   - error: An error if configuration creation fails
//
// Example:
//
//	// Basic configuration
//	cfg, err := NewMpesaConfig(
//	    "consumer_key",
//	    "consumer_secret",
//	    Sandbox,
//	    nil, nil, nil, nil, nil,
//	)
//
//	// Configuration with optional parameters
//	businessCode := "174379"
//	passKey := "your_passkey"
//	cfg, err := NewMpesaConfig(
//	    "consumer_key",
//	    "consumer_secret",
//	    Sandbox,
//	    &businessCode,
//	    &passKey,
//	    nil, nil, nil,
//	)
func NewMpesaConfig(
	consumerKey, consumerSecret string,
	environment Environment,
	businessCode, passKey, securityCredential, queueTimeoutURL, resultURL *string,
) (*MpesaConfig, error) {
	env := strings.ToLower(string(environment))
	baseURL := "https://sandbox.safaricom.co.ke"
	if env == string(Production) {
		baseURL = "https://api.safaricom.co.ke"
	}

	cfg := &MpesaConfig{
		consumerKey:        consumerKey,
		consumerSecret:     consumerSecret,
		environment:        Environment(env),
		baseURL:            baseURL,
		businessCode:       getOrDefault(businessCode, ""),
		passKey:            getOrDefault(passKey, ""),
		securityCredential: getOrDefault(securityCredential, ""),
		queueTimeoutURL:    getOrDefault(queueTimeoutURL, ""),
		resultURL:          getOrDefault(resultURL, ""),
	}

	return cfg, nil
}

// Getters

// GetConsumerKey returns the consumer key used for API authentication.
//
// Returns:
//   - string: The consumer key obtained from Safaricom Developer Portal
func (cfg *MpesaConfig) GetConsumerKey() string {
	return cfg.consumerKey
}

// GetConsumerSecret returns the consumer secret used for API authentication.
//
// Returns:
//   - string: The consumer secret obtained from Safaricom Developer Portal
func (cfg *MpesaConfig) GetConsumerSecret() string {
	return cfg.consumerSecret
}

// GetEnvironment returns the current M-Pesa API environment setting.
//
// Returns:
//   - Environment: The configured environment (Sandbox or Production)
func (cfg *MpesaConfig) GetEnvironment() Environment {
	return cfg.environment
}

// GetBaseURL returns the base URL for M-Pesa API endpoints.
// The URL is automatically determined based on the environment setting.
//
// Returns:
//   - string: The base URL for API requests
//
// Example:
//   - Sandbox: "https://sandbox.safaricom.co.ke"
//   - Production: "https://api.safaricom.co.ke"
func (cfg *MpesaConfig) GetBaseURL() string {
	return cfg.baseURL
}

// GetBusinessCode returns the business shortcode for M-Pesa transactions.
//
// Returns:
//   - string: The business shortcode, or empty string if not set
func (cfg *MpesaConfig) GetBusinessCode() string {
	return cfg.businessCode
}

// GetPassKey returns the Lipa na M-Pesa Online passkey for STK Push operations.
//
// Returns:
//   - string: The passkey, or empty string if not set
func (cfg *MpesaConfig) GetPassKey() string {
	return cfg.passKey
}

// GetSecurityCredential returns the encrypted security credential for B2C and reversal operations.
//
// Returns:
//   - string: The base64-encoded encrypted security credential, or empty string if not set
func (cfg *MpesaConfig) GetSecurityCredential() string {
	return cfg.securityCredential
}

// GetQueueTimeoutURL returns the URL for queue timeout notifications.
//
// Returns:
//   - string: The queue timeout URL, or empty string if not set
func (cfg *MpesaConfig) GetQueueTimeoutURL() string {
	return cfg.queueTimeoutURL
}

// GetResultURL returns the URL for transaction result notifications.
//
// Returns:
//   - string: The result URL, or empty string if not set
func (cfg *MpesaConfig) GetResultURL() string {
	return cfg.resultURL
}

// Setters

// SetBusinessCode sets the business shortcode for M-Pesa transactions.
// This shortcode identifies your business in the M-Pesa system.
//
// Parameters:
//   - code: The business shortcode (e.g., "174379" for sandbox)
//
// Example:
//
//	cfg.SetBusinessCode("174379") // Sandbox
//	cfg.SetBusinessCode("123456") // Production
func (cfg *MpesaConfig) SetBusinessCode(code string) {
	cfg.businessCode = code
}

// SetPassKey sets the Lipa na M-Pesa Online passkey for STK Push operations.
// The passkey is obtained from Safaricom and is required for initiating STK Push requests.
//
// Parameters:
//   - key: The Lipa na M-Pesa Online passkey
//
// Example:
//
//	cfg.SetPassKey("bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919")
func (cfg *MpesaConfig) SetPassKey(key string) {
	cfg.passKey = key
}

// SetQueueTimeoutURL sets the URL where M-Pesa will send queue timeout notifications.
//
// Parameters:
//   - url: The fully qualified queue timeout URL
//
// Example:
//
//	cfg.SetQueueTimeoutURL("https://yourdomain.com/mpesa/timeout")
func (cfg *MpesaConfig) SetQueueTimeoutURL(url string) {
	cfg.queueTimeoutURL = url
}

// SetResultURL sets the URL where M-Pesa will send transaction result notifications.
//
// Parameters:
//   - url: The fully qualified result URL
//
// Example:
//
//	cfg.SetResultURL("https://yourdomain.com/mpesa/result")
func (cfg *MpesaConfig) SetResultURL(url string) {
	cfg.resultURL = url
}

// SetSecurityCredential encrypts an initiator password and sets it as the security credential.
// This credential is required for B2C transactions, reversals, and other operations that
// require initiator authentication. The password is encrypted using AES-256-CBC encryption.
//
// Parameters:
//   - initiatorPassword: The plain text initiator password
//
// Returns:
//   - error: An error if encryption fails
//
// Example:
//
//	err := cfg.SetSecurityCredential("myInitiatorPassword")
//	if err != nil {
//	    log.Printf("Failed to set security credential: %v", err)
//	}
func (cfg *MpesaConfig) SetSecurityCredential(initiatorPassword string) error {
	encryptionKey := []byte("mypasswordmypasswordmypassword12") // 32 bytes for AES-256
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	encrypter := cipher.NewCBCEncrypter(block, iv)

	plaintext := pad([]byte(fmt.Sprintf("%s + Certificate", initiatorPassword)), aes.BlockSize)
	ciphertext := make([]byte, len(plaintext))
	encrypter.CryptBlocks(ciphertext, plaintext)

	combined := append(iv, ciphertext...)
	cfg.securityCredential = base64.StdEncoding.EncodeToString(combined)

	return nil
}

// Helper function to return dereferenced pointer or default
func getOrDefault(val *string, fallback string) string {
	if val != nil {
		return *val
	}
	return fallback
}

// PKCS7 padding for AES
func pad(src []byte, blockSize int) []byte {
	padLen := blockSize - len(src)%blockSize
	pad := bytesRepeat(byte(padLen), padLen)
	return append(src, pad...)
}

// bytesRepeat returns a new byte slice repeating b n times
func bytesRepeat(b byte, count int) []byte {
	buf := make([]byte, count)
	for i := range buf {
		buf[i] = b
	}
	return buf
}
