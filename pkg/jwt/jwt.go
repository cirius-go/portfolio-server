package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/cirius-go/portfolio-server/pkg/errors"
)

var (
	ErrInvalidToken = errors.NewUnauthorized(nil, "invalid token")
)

type Config struct {
	alg    *jwt.SigningMethodHMAC
	secret []byte
}

func C() *Config {
	return &Config{}
}

func (c *Config) Alg(alg *jwt.SigningMethodHMAC) *Config {
	c.alg = alg
	return c
}

func (c *Config) Secret(secret []byte) *Config {
	c.secret = secret
	return c
}

// JWT represents the jwt service.
type JWT struct {
	cfg *Config
}

// NewJWTWithConfig creates a new JWT service with custom config.
func NewJWTWithConfig(c *Config) *JWT {
	if c.alg == nil || c.alg.Alg() == "" {
		panic("alg must not be empty")
	}

	if len(string(c.secret)) < 16 {
		panic("secret length must be at least 16")
	}

	return &JWT{
		cfg: c,
	}
}

// NewToken generates a new JWT token.
func (s *JWT) NewToken(payload jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(s.cfg.alg, payload)
	return token.SignedString(s.cfg.secret)
}

// ParseToken parses a token and returns the claims.
func (s *JWT) ParseToken(token string, customClaims jwt.Claims) error {
	tk, err := jwt.ParseWithClaims(token, customClaims, func(t *jwt.Token) (any, error) {
		if m, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || m != s.cfg.alg {
			return nil, ErrInvalidToken
		}

		return s.cfg.secret, nil
	})
	if err != nil {
		return err
	}
	if !tk.Valid {
		return ErrInvalidToken
	}
	return nil
}
