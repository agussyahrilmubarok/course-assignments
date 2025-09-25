package password

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Encrypt hashes a plain text password using bcrypt
func Encrypt(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("cannot encrypt the password")
		return ""
	}
	return string(hashedPassword)
}

// Matchers verifies if a plain text password matches the hashed password
func Matchers(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
