package jwt

import (
	"errors"
	"time"

	"github.com/pascaldekloe/jwt"
)

// NewToken -
func NewToken(jwtSecret, subject string, expiration time.Time, set map[string]interface{}) (string, error) {
	claims := jwt.Claims{}
	claims.Subject = subject
	claims.Expires = jwt.NewNumericTime(expiration)
	claims.Set = set
	token, err := claims.HMACSign(jwt.HS256, []byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return string(token), nil
}

// ParseToken -
func ParseToken(jwtSecret, token string) (*jwt.Claims, error) {
	return jwt.HMACCheck([]byte(token), []byte(jwtSecret))
}

// ValidateToken -
func ValidateToken(jwtSecret, token string) (*jwt.Claims, error) {
	claims, err := ParseToken(jwtSecret, token)
	if err != nil {
		return nil, err
	}
	if valid := claims.Valid(time.Now()); !valid {
		return nil, errors.New("token is expired")
	}
	return claims, nil
}
