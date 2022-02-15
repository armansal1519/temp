package jwt

import (
	"bamachoub-backend-go-v1/config"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type AdminTokenPayload struct {
	Key    string
	Access string
}

func GenerateAdminToken(payload *AdminTokenPayload, isRefresh bool) string {
	v, err := time.ParseDuration(config.ACCESS_TOKEN_EXP)

	if err != nil {
		panic("Invalid time duration. Should be time.ParseDuration string")
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":    time.Now().Add(v).Unix(),
		"Key":    payload.Key,
		"Access": payload.Access,
	})
	var key string
	if isRefresh {
		key = config.ADMIN_REFRESH_TOKEN_KEY
	} else {
		key = config.ADMIN_ACCESS_TOKEN_KEY

	}

	token, err := t.SignedString([]byte(key))

	if err != nil {
		panic(err)
	}

	return token
}

func AdminParse(token string, isRefreshToken bool) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println(6)
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		if isRefreshToken {
			return []byte(config.ADMIN_REFRESH_TOKEN_KEY), nil
		}
		return []byte(config.ADMIN_ACCESS_TOKEN_KEY), nil
	})
}

//Verify verifies the jwt token against the secret
func VerifyAdmin(token string, isRefresh bool) (*AdminTokenPayload, error) {
	parsed, err := AdminParse(token, isRefresh)

	if err != nil {
		log.Println(5)
		return nil, err
	}

	// Parsing token claims
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}
	fmt.Println(1, claims)
	// Getting ID, it's an interface{} so I need to cast it to uint
	key, ok := claims["Key"].(string)
	if !ok {
		return nil, errors.New("Something went wrong in VerifyAdmin1")
	}

	access, ok := claims["Access"].(string)
	if !ok {
		log.Println(access)
		return nil, errors.New("Something went wrong in VerifyAdmin3")
	}

	return &AdminTokenPayload{
		Key:    key,
		Access: access,
	}, nil
}
