package jwt

import (
	"bamachoub-backend-go-v1/config"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenPayload defines the payload for the token
type TokenPayload struct {
	Key string
}

// Generate generates the jwt token based on payload
func GenerateAccessToken(payload *TokenPayload) string {
	v, err := time.ParseDuration(config.ACCESS_TOKEN_EXP)

	if err != nil {
		panic("Invalid time duration. Should be time.ParseDuration string")
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(v).Unix(),
		"Key": payload.Key,
	})

	token, err := t.SignedString([]byte(config.ACCESS_TOKEN_KEY))

	if err != nil {
		panic(err)
	}

	return token
}
func GenerateRefreshToken(payload *TokenPayload) string {
	v, err := time.ParseDuration(config.REFRESH_TOKEN_EXP)

	if err != nil {
		panic("Invalid time duration. Should be time.ParseDuration string")
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(v).Unix(),
		"Key": payload.Key,
	})

	token, err := t.SignedString([]byte(config.REFRESH_TOKEN_KEY))

	if err != nil {
		panic(err)
	}

	return token
}

func parse(token string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.ACCESS_TOKEN_KEY), nil
	})
}

//Verify verifies the jwt token against the secret
func Verify(token string) (*TokenPayload, error) {
	parsed, err := parse(token)

	if err != nil {
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
		return nil, errors.New("Something went wrong")
	}

	return &TokenPayload{
		Key: key,
	}, nil
}

func VerifyRefreshToken(token string) (*TokenPayload, error) {
	parsed, err := parseRefreshToken(token)

	if err != nil {
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
		return nil, errors.New("Something went wrong")
	}

	return &TokenPayload{
		Key: key,
	}, nil
}
func parseRefreshToken(token string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.REFRESH_TOKEN_KEY), nil
	})
}
