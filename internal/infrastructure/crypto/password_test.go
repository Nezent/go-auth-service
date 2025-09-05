package crypto

import (
	"encoding/base64"
	"testing"

	"github.com/Nezent/auth-service/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestConfig() {
	testConfig := config.Config{
		Argon2id: config.Argon2idConfig{
			Time:    1,
			Memory:  64 * 1024, // 64 MB
			Threads: 4,
			KeyLen:  32,
			SaltLen: 16,
		},
	}
	InitArgon2Config(testConfig)
}

func TestGenerateSalt(t *testing.T) {
	setupTestConfig()

	t.Run("should generate valid salt", func(t *testing.T) {
		salt, err := GenerateSalt()

		require.NoError(t, err)
		assert.NotEmpty(t, salt)

		// Check if salt can be decoded
		saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
		require.NoError(t, err)
		assert.Equal(t, cfg.Argon2id.SaltLen, len(saltBytes))
	})

	t.Run("should generate different salts", func(t *testing.T) {
		salt1, err1 := GenerateSalt()
		salt2, err2 := GenerateSalt()

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, salt1, salt2)
	})
}

func TestHashPassword(t *testing.T) {
	setupTestConfig()

	t.Run("should hash password successfully", func(t *testing.T) {
		password := "testpassword123"

		hash, salt, err := HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEmpty(t, salt)
		assert.NotEqual(t, password, hash)
	})

	t.Run("should generate different hashes for same password", func(t *testing.T) {
		password := "testpassword123"

		hash1, salt1, err1 := HashPassword(password)
		hash2, salt2, err2 := HashPassword(password)

		require.NoError(t, err1)
		require.NoError(t, err2)

		// Different salts should produce different hashes
		assert.NotEqual(t, salt1, salt2)
		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("should handle empty password", func(t *testing.T) {
		hash, salt, err := HashPassword("")

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEmpty(t, salt)
	})
}

func TestHashPasswordWithSalt(t *testing.T) {
	setupTestConfig()

	t.Run("should hash password with provided salt", func(t *testing.T) {
		password := "testpassword123"
		salt := "dGVzdHNhbHQxMjM0NTY3OA" // base64 encoded test salt

		hash, err := HashPasswordWithSalt(password, salt)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})

	t.Run("should produce same hash for same password and salt", func(t *testing.T) {
		password := "testpassword123"
		salt := "dGVzdHNhbHQxMjM0NTY3OA"

		hash1, err1 := HashPasswordWithSalt(password, salt)
		hash2, err2 := HashPasswordWithSalt(password, salt)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.Equal(t, hash1, hash2)
	})

	t.Run("should return error for invalid salt", func(t *testing.T) {
		password := "testpassword123"
		invalidSalt := "invalid-base64!"

		_, err := HashPasswordWithSalt(password, invalidSalt)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid salt encoding")
	})

	t.Run("should produce different hashes for different passwords", func(t *testing.T) {
		password1 := "testpassword123"
		password2 := "differentpassword456"
		salt := "dGVzdHNhbHQxMjM0NTY3OA"

		hash1, err1 := HashPasswordWithSalt(password1, salt)
		hash2, err2 := HashPasswordWithSalt(password2, salt)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2)
	})
}

func TestVerifyPassword(t *testing.T) {
	setupTestConfig()

	t.Run("should verify correct password", func(t *testing.T) {
		password := "testpassword123"
		hash, salt, err := HashPassword(password)
		require.NoError(t, err)

		isValid := VerifyPassword(password, salt, hash)

		assert.True(t, isValid)
	})

	t.Run("should reject incorrect password", func(t *testing.T) {
		password := "testpassword123"
		wrongPassword := "wrongpassword456"
		hash, salt, err := HashPassword(password)
		require.NoError(t, err)

		isValid := VerifyPassword(wrongPassword, salt, hash)

		assert.False(t, isValid)
	})

	t.Run("should reject with wrong salt", func(t *testing.T) {
		password := "testpassword123"
		hash, _, err := HashPassword(password)
		require.NoError(t, err)

		wrongSalt := "d3JvbmdzYWx0MTIzNDU2Nzg"
		isValid := VerifyPassword(password, wrongSalt, hash)

		assert.False(t, isValid)
	})

	t.Run("should reject with wrong hash", func(t *testing.T) {
		password := "testpassword123"
		_, salt, err := HashPassword(password)
		require.NoError(t, err)

		wrongHash := "aW52YWxpZGhhc2gxMjM0NTY3OA"
		isValid := VerifyPassword(password, salt, wrongHash)

		assert.False(t, isValid)
	})

	t.Run("should handle invalid salt gracefully", func(t *testing.T) {
		password := "testpassword123"
		invalidSalt := "invalid-base64!"
		hash := "dGVzdGhhc2gxMjM0NTY3OA"

		isValid := VerifyPassword(password, invalidSalt, hash)

		assert.False(t, isValid)
	})
}

func TestInitArgon2Config(t *testing.T) {
	t.Run("should initialize config", func(t *testing.T) {
		testConfig := config.Config{
			Argon2id: config.Argon2idConfig{
				Time:    2,
				Memory:  128 * 1024,
				Threads: 8,
				KeyLen:  64,
				SaltLen: 32,
			},
		}

		InitArgon2Config(testConfig)

		assert.Equal(t, testConfig.Argon2id.Time, cfg.Argon2id.Time)
		assert.Equal(t, testConfig.Argon2id.Memory, cfg.Argon2id.Memory)
		assert.Equal(t, testConfig.Argon2id.Threads, cfg.Argon2id.Threads)
		assert.Equal(t, testConfig.Argon2id.KeyLen, cfg.Argon2id.KeyLen)
		assert.Equal(t, testConfig.Argon2id.SaltLen, cfg.Argon2id.SaltLen)
	})
}

// Benchmark tests
func BenchmarkHashPassword(b *testing.B) {
	setupTestConfig()
	password := "benchmarkpassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := HashPassword(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	setupTestConfig()
	password := "benchmarkpassword123"
	hash, salt, err := HashPassword(password)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyPassword(password, salt, hash)
	}
}
