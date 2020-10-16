package global

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/square/go-jose/v3"
)

var (
	// JwtPrivateKey ...
	JwtPrivateKey *rsa.PrivateKey
	// JwtPublicKey ...
	JwtPublicKey *rsa.PublicKey
	// JwtCertificates ...
	JwtCertificates []*x509.Certificate
	// Jwk ...
	Jwk jose.JSONWebKey
)

// Init ...
func Init() {
	initPrivate()
	initCert()
	initJwk()
}

func initPrivate() {
	var (
		rsaPrivatePEMBlock *pem.Block
		err                error
	)
	rsaPrivatePEMBlock, _ = pem.Decode([]byte(viper.GetString("jwt.rsa.private")))
	if rsaPrivatePEMBlock == nil {
		logrus.Fatalln("failed to parse certificate PEM")
		return
	}
	JwtPrivateKey, err = x509.ParsePKCS1PrivateKey(rsaPrivatePEMBlock.Bytes)
	if err != nil {
		logrus.Fatalf("x509.ParsePKCS1PrivateKey Error: %s", err)
		return
	}
	JwtPublicKey = &JwtPrivateKey.PublicKey
}

func initCert() {
	var (
		rsaCertPEMBlock *pem.Block
		err             error
	)
	rsaCertPEMBlock, _ = pem.Decode([]byte(viper.GetString("jwt.rsa.cert")))
	if rsaCertPEMBlock == nil {
		logrus.Fatalln("failed to parse certificate PEM")
		return
	}
	JwtCertificates, err = x509.ParseCertificates(rsaCertPEMBlock.Bytes)
	if err != nil {
		logrus.Fatalln("failed to parse certificate: %s", err)
		return
	}
}

func initJwk() {
	x5tSHA1 := sha1.Sum(JwtCertificates[0].Raw)
	x5tSHA256 := sha256.Sum256(JwtCertificates[0].Raw)

	Jwk = jose.JSONWebKey{
		Key:                         JwtCertificates[0].PublicKey,
		KeyID:                       "naas",
		Algorithm:                   "RS256",
		Use:                         "sig",
		Certificates:                JwtCertificates,
		CertificateThumbprintSHA1:   x5tSHA1[:],
		CertificateThumbprintSHA256: x5tSHA256[:],
	}
	if !Jwk.Valid() {
		logrus.Fatalf("Jwk.Valid: false")
	}
}
