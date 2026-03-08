package password

import "golang.org/x/crypto/bcrypt"

func GenerateHash(plainText string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CompareHash(plainText, hashText string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashText), []byte(plainText)) == nil
}
