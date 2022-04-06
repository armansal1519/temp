package password

import (
	"golang.org/x/crypto/bcrypt"
)

// Generate return a hashed password
//func Generate(raw string) string {
//	//hash, err := bcrypt.GenerateFromPassword([]byte(raw), 10)
//	//
//	//if err != nil {
//	//	panic(err)
//	//}
//	hash, err := argon2id.CreateHash(raw, argon2id.DefaultParams)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return hash
//}

// CheckPasswordHash compares a hashed password with plaintext password
//func CheckPasswordHash( raw string,hash string) bool {
//	match, err := argon2id.ComparePasswordAndHash(raw, hash)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return match
//}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
