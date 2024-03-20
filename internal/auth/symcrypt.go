package auth

import (
	"crypto/sha256"
	"log"
	"sync"
	"tasklify/internal/config"

	"github.com/gtank/cryptopasta"
)

const (
	symcryptSalt = "Fd7nQfuc8Q3Zo3XgpxhDVvyqYeL0trkYJg0"
)

type Symcrypt interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type symcrypt struct {
	Key *[32]byte
}

var (
	onceSymcrypt sync.Once

	symcryptClient *symcrypt
)

func GetSymcrypt(config ...*config.Config) Symcrypt {

	onceSymcrypt.Do(func() { // <-- atomic, does not allow repeating
		config := config[0]

		symcryptClient = loadSymcrypt(config.Auth)
	})

	return symcryptClient
}

func loadSymcrypt(config config.Auth) *symcrypt {
	symcryptKey := sha256.Sum256([]byte(symcryptSalt + config.SymcryptKey))

	log.Println("Symcrypt loaded")

	return &symcrypt{Key: &symcryptKey}
}

func (s *symcrypt) Encrypt(plaintext string) (string, error) {
	ciphertext, err := cryptopasta.Encrypt([]byte(plaintext), s.Key)
	if err != nil {
		return "", err
	}

	return string(ciphertext), nil
}

func (s *symcrypt) Decrypt(ciphertext string) (string, error) {
	plaintext, err := cryptopasta.Decrypt([]byte(ciphertext), s.Key)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
