package jwtc

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	sctx "github.com/viettranx/service-context"
	"time"
)

const (
	defaultSecret               = "very-important-please-change-it!" // in 32 bytes
	defaultExpireTokenInSeconds = 60 * 60 * 24 * 7                   // 7d
)

var (
	ErrSecretKeyNotValid     = errors.New("secret key must be in 32 bytes")
	ErrTokenLifeTimeTooShort = errors.New("token life time too short")
)

type jwtx struct {
	id                   string
	secret               string
	expireTokenInSeconds int
}

func NewJWT(id string) *jwtx {
	return &jwtx{id: id}
}

func (j *jwtx) ID() string {
	return j.id
}

func (j *jwtx) InitFlags() {
	flag.StringVar(
		&j.secret,
		"jwt-secret",
		defaultSecret,
		"Secret key to sign JWT",
	)

	flag.IntVar(
		&j.expireTokenInSeconds,
		"jwt-exp-secs",
		defaultExpireTokenInSeconds,
		"Number of seconds token will expired",
	)
}

func (j *jwtx) Activate(_ sctx.ServiceContext) error {
	if len(j.secret) < 32 {
		return errors.WithStack(ErrSecretKeyNotValid)
	}

	if j.expireTokenInSeconds <= 60 {
		return errors.WithStack(ErrTokenLifeTimeTooShort)
	}

	return nil
}

func (j *jwtx) Stop() error {
	return nil
}

func (j *jwtx) IssueToken(ctx context.Context, id, sub string) (token string, expSecs int, err error) {
	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(j.expireTokenInSeconds))),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        id,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSignedStr, err := t.SignedString([]byte(j.secret))

	if err != nil {
		return "", 0, errors.WithStack(err)
	}

	return tokenSignedStr, j.expireTokenInSeconds, nil
}

func (j *jwtx) ParseToken(ctx context.Context, tokenString string) (claims *jwt.RegisteredClaims, err error) {
	var rc jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(tokenString, &rc, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.secret), nil
	})

	if !token.Valid {
		return nil, errors.WithStack(err)
	}

	return &rc, nil
}
