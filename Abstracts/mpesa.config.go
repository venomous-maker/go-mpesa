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

type Environment string

const (
	Sandbox    Environment = "sandbox"
	Production Environment = "live"
)

type MpesaConfig struct {
	ConsumerKey        string
	ConsumerSecret     string
	Environment        Environment
	BaseURL            string
	BusinessCode       string
	PassKey            string
	SecurityCredential string
	QueueTimeoutURL    string
	ResultURL          string
}

// NewMpesaConfig initializes the config with optional values.
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
		ConsumerKey:        consumerKey,
		ConsumerSecret:     consumerSecret,
		Environment:        Environment(env),
		BaseURL:            baseURL,
		BusinessCode:       getOrDefault(businessCode, ""),
		PassKey:            getOrDefault(passKey, ""),
		SecurityCredential: getOrDefault(securityCredential, ""),
		QueueTimeoutURL:    getOrDefault(queueTimeoutURL, ""),
		ResultURL:          getOrDefault(resultURL, ""),
	}

	return cfg, nil
}

// SetSecurityCredential encrypts a password with AES-256-CBC and sets the security credential.
func (cfg *MpesaConfig) SetSecurityCredential(initiatorPassword string) error {
	// Dummy password and IV, you should replace with actual logic/certificate-based key
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

	// Pad plaintext to block size
	plaintext := pad([]byte(fmt.Sprintf("%s + Certificate", initiatorPassword)), aes.BlockSize)
	ciphertext := make([]byte, len(plaintext))
	encrypter.CryptBlocks(ciphertext, plaintext)

	combined := append(iv, ciphertext...)
	cfg.SecurityCredential = base64.StdEncoding.EncodeToString(combined)

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
