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
	consumerKey        string
	consumerSecret     string
	environment        Environment
	baseURL            string
	businessCode       string
	passKey            string
	securityCredential string
	queueTimeoutURL    string
	resultURL          string
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

func (cfg *MpesaConfig) GetConsumerKey() string {
	return cfg.consumerKey
}

func (cfg *MpesaConfig) GetConsumerSecret() string {
	return cfg.consumerSecret
}

func (cfg *MpesaConfig) GetEnvironment() Environment {
	return cfg.environment
}

func (cfg *MpesaConfig) GetBaseURL() string {
	return cfg.baseURL
}

func (cfg *MpesaConfig) GetBusinessCode() string {
	return cfg.businessCode
}

func (cfg *MpesaConfig) GetPassKey() string {
	return cfg.passKey
}

func (cfg *MpesaConfig) GetSecurityCredential() string {
	return cfg.securityCredential
}

func (cfg *MpesaConfig) GetQueueTimeoutURL() string {
	return cfg.queueTimeoutURL
}

func (cfg *MpesaConfig) GetResultURL() string {
	return cfg.resultURL
}

// Setters

func (cfg *MpesaConfig) SetBusinessCode(code string) {
	cfg.businessCode = code
}

func (cfg *MpesaConfig) SetPassKey(key string) {
	cfg.passKey = key
}

func (cfg *MpesaConfig) SetQueueTimeoutURL(url string) {
	cfg.queueTimeoutURL = url
}

func (cfg *MpesaConfig) SetResultURL(url string) {
	cfg.resultURL = url
}

// SetSecurityCredential encrypts a password with AES-256-CBC and sets the security credential.
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
