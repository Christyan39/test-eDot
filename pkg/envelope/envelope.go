package envelope

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// KeySize for AES-256
	KeySize = 32
	// NonceSize for AES-GCM
	NonceSize = 12
)

// EnvelopeService handles secure envelope encryption/decryption
type EnvelopeService struct {
	secretKey []byte
}

// NewEnvelopeService creates a new envelope service
func NewEnvelopeService(secret string) *EnvelopeService {
	// Derive a key from the secret using PBKDF2
	salt := []byte("edot-envelope-salt") // In production, use a random salt
	key := pbkdf2.Key([]byte(secret), salt, 10000, KeySize, sha256.New)

	return &EnvelopeService{
		secretKey: key,
	}
}

// EncryptData encrypts data into a secure envelope (simplified without nonce and timestamp)
func (e *EnvelopeService) EncryptData(data interface{}) (string, error) {
	// Serialize data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to serialize data: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(e.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data (directly encrypt JSON without timestamp)
	ciphertext := gcm.Seal(nil, nonce, jsonData, nil)

	// Combine nonce and ciphertext
	combined := append(nonce, ciphertext...)
	encodedData := base64.StdEncoding.EncodeToString(combined)

	return encodedData, nil
}

// DecryptData decrypts data from a secure envelope (simplified without nonce and timestamp validation)
func (e *EnvelopeService) DecryptData(encodedData string, target interface{}) error {
	// Decode base64 data
	combined, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}

	if len(combined) < NonceSize {
		return errors.New("invalid envelope data")
	}

	// Extract nonce and ciphertext
	nonce := combined[:NonceSize]
	ciphertext := combined[NonceSize:]

	// Create AES cipher
	block, err := aes.NewCipher(e.secretKey)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt data
	jsonData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	// Unmarshal JSON data directly
	if err := json.Unmarshal(jsonData, target); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}
