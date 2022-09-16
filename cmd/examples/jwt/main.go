package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/subtle"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"time"
)

var issuer = "example.com"
var audience = "example.com"

func main() {
	bitSize := 1024

	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		logrus.Fatalln("failed generating random key:", err)
	}

	// Extract public component.
	pub := key.Public()
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)

	// Encode private key to PKCS#1 ASN.1 PEM.
	privPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)

	// Encode public key to PKCS#1 ASN.1 PEM.
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)

	// save pem to file

	// load pem from file and parse
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privPEM)
	if err != nil {
		logrus.Fatalln(err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubPEM)
	if err != nil {
		logrus.Fatalln(err)
	}

	token, err := newToken("user", privateKey)
	if err != nil {
		logrus.Fatalln(err)
	}

	logrus.Infoln("New Token:", token)

	j, claims, err := parseToken(token, pubKey)
	if err != nil {
		logrus.Fatalln(err)
	}

	logrus.Infof("Token valid=%v and issuer is %s", j.Valid, claims.Issuer)
}

func newToken(subject string, privKey interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, &jwt.RegisteredClaims{
		Issuer:    issuer,             // Auth Server who generated the token
		Subject:   subject,            // UserId
		Audience:  []string{audience}, // where the token will be validated
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		NotBefore: jwt.NewNumericDate(time.Now().Add(time.Millisecond * 50)), // adds leeway
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "TokenId",
	})
	acToken, err := token.SignedString(privKey)
	if err != nil {
		return "", err
	}
	return acToken, nil
}

func parseToken(token string, pubKey interface{}) (*jwt.Token, *jwt.RegisteredClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil {
		return nil, nil, err
	}
	if !parsedToken.Valid {
		return nil, nil, fmt.Errorf("token invalid")
	}
	if parsedToken.Method != jwt.SigningMethodRS512 {
		return nil, nil, fmt.Errorf("wrong signing method")
	}
	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, nil, fmt.Errorf("invalid claims")
	}
	now := time.Now()
	nbfOk := claims.VerifyNotBefore(now, true)
	if !nbfOk {
		return nil, nil, fmt.Errorf("invalid nbf")
	}

	iatOk := claims.VerifyIssuedAt(now, true)
	if !iatOk {
		return nil, nil, fmt.Errorf("invalid iat")
	}

	expOk := claims.VerifyExpiresAt(now, true)
	if !expOk {
		return nil, nil, fmt.Errorf("invalid exp")
	}

	audienceOk := claims.VerifyAudience(audience, true)
	if !audienceOk {
		return nil, nil, fmt.Errorf("invalid audience")
	}

	issuerOk := verifyIss(claims.Issuer, issuer)
	if !issuerOk {
		return nil, nil, fmt.Errorf("invalid issuer")
	}

	return parsedToken, claims, nil
}

func verifyIss(iss string, cmp string) bool {
	if iss == "" {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(iss), []byte(cmp)) != 0 {
		return true
	} else {
		return false
	}
}
