package cryptohandler

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// TestLoadEnv tests loading the environment variable.
func TestLoadEnv(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	key := os.Getenv("AES_KEY")
	if key == "" {
		t.Fatal("AES_KEY is not set in the .env file")
	}

	expectedKeyLength := 32 // For AES-256
	if len(key) != expectedKeyLength {
		t.Fatalf("Invalid AES key length: got %d, expected %d", len(key), expectedKeyLength)
	}
}

// TestEncryptDecrypt tests encryption and decryption.
func TestEncryptDecrypt(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	key := os.Getenv("AES_KEY")
	if len(key) != 32 {
		t.Fatalf("Invalid AES key length: got %d, expected 32", len(key))
	}

	plaintext := "This is a test string."

	// Encrypt the plaintext
	encrypted, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt the ciphertext
	decrypted, err := Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Verify the decrypted text matches the original plaintext
	if decrypted != plaintext {
		t.Fatalf("Decrypted text does not match plaintext: got %q, expected %q", decrypted, plaintext)
	}
}

// TestInvalidKey tests decryption with an invalid key.
func TestInvalidKey(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	key := os.Getenv("AES_KEY")
	if len(key) != 32 {
		t.Fatalf("Invalid AES key length: got %d, expected 32", len(key))
	}

	plaintext := "This is a test string."

	// Encrypt the plaintext
	encrypted, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Attempt to decrypt with a wrong key
	wrongKey := "thisisnottherightkey!!!!!!!!!"
	_, err = Decrypt(wrongKey, encrypted)
	if err == nil {
		t.Fatal("Decryption should have failed with an invalid key")
	}
}
