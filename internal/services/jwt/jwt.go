package jwt

import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	// for now there's no reason for err segregation & uniq processing
	// but its good idea to have list of error which module can return
	ErrKeyParsing      = fmt.Errorf("parsing error")
	ErrTokenGeneration = fmt.Errorf("token generation error")
	ErrSigning         = fmt.Errorf("signing error")
	ErrValidation      = fmt.Errorf("token validation errror")
)

type JWTManager struct {
	issuer     string
	expiresIn  time.Duration
	publicKey  jwk.Key
	privateKey jwk.Key
	signAlg    jwa.SignatureAlgorithm
}

func NewJWTManager(issuer string, expiresIn time.Duration, publicKey, privateKey []byte, signAlg jwa.SignatureAlgorithm) (*JWTManager, error) {
	pubKey, err := jwk.ParseKey(publicKey, jwk.WithPEM(true))
	if err != nil {
		return nil, fmt.Errorf("public key %w", ErrKeyParsing)
	}
	privKey, err := jwk.ParseKey(privateKey, jwk.WithPEM(true))
	if err != nil {
		return nil, fmt.Errorf("private key %w", ErrKeyParsing)
	}

	return &JWTManager{
		issuer:     issuer,
		expiresIn:  expiresIn,
		publicKey:  pubKey,
		privateKey: privKey,
		signAlg:    signAlg,
	}, nil
}

func (j *JWTManager) IssueToken(userID string) (string, error) {
	t, err := jwt.NewBuilder().
		IssuedAt(time.Now()).
		Issuer(j.issuer).
		Expiration(time.Now().Add(j.expiresIn)).
		Subject(userID).
		Build()
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrTokenGeneration, err)
	}
	signed, err := jwt.Sign(t, jwt.WithKey(j.signAlg, j.privateKey))
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrSigning, err)
	}
	return string(signed), nil
}

func (j *JWTManager) VerifyToken(payload []byte) (jwt.Token, error) {
	t, err := jwt.Parse(payload,
		jwt.WithKey(j.signAlg, j.publicKey),
		jwt.WithVerify(true),
		jwt.WithIssuer(j.issuer),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidation, err)
	}
	return t, nil
}
