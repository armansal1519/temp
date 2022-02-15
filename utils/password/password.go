package password

import (
	"github.com/alexedwards/argon2id"
	"log"
)

// Generate return a hashed password
func Generate(raw string) string {
	//hash, err := bcrypt.GenerateFromPassword([]byte(raw), 10)
	//
	//if err != nil {
	//	panic(err)
	//}
	hash, err := argon2id.CreateHash(raw, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}

	return hash
}

// Verify compares a hashed password with plaintext password
func Verify(hash string, raw string) bool {
	match, err := argon2id.ComparePasswordAndHash(raw, hash)
	if err != nil {
		log.Fatal(err)
	}
	return match
}
