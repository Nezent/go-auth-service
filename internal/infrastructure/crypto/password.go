package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/Nezent/auth-service/internal/infrastructure/config"
	"golang.org/x/crypto/argon2"
)

var cfg config.Config

func InitArgon2Config(c config.Config) {
	cfg = c
}

// GenerateSalt creates a random salt for password hashing.
func GenerateSalt() (string, error) {
	salt := make([]byte, cfg.Argon2id.SaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.RawStdEncoding.EncodeToString(salt), nil
}

// HashPassword generates a random salt and hashes the password using Argon2id.
// Returns the hash and the salt.
func HashPassword(password string) (hash string, salt string, err error) {
	salt, err = GenerateSalt()
	if err != nil {
		return "", "", err
	}
	hash, err = HashPasswordWithSalt(password, salt)
	return hash, salt, err
}

// HashPasswordWithSalt hashes a password using Argon2id and a provided salt.
func HashPasswordWithSalt(password, salt string) (string, error) {
	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return "", fmt.Errorf("invalid salt encoding: %w", err)
	}
	key := argon2.IDKey(
		[]byte(password),
		saltBytes,
		cfg.Argon2id.Time,
		cfg.Argon2id.Memory,
		cfg.Argon2id.Threads,
		cfg.Argon2id.KeyLen,
	)
	return base64.RawStdEncoding.EncodeToString(key), nil
}

// VerifyPassword checks if the password matches the hash using the same salt.
func VerifyPassword(password, salt, hash string) bool {
	computed, err := HashPasswordWithSalt(password, salt)
	if err != nil {
		return false
	}
	return computed == hash
}
