package oidc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
)

func TestCreateJwtOpenID_H256(t *testing.T) {
	key := []byte("secret")
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: key}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		panic(err)
	}

	cl := jwt.Claims{
		Subject:   "subject",
		Issuer:    "issuer",
		NotBefore: jwt.NewNumericDate(time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)),
		Audience:  jwt.Audience{"leela", "fry"},
	}
	raw, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	if err != nil {
		panic(err)
	}

	fmt.Println(raw)
}

// Decode base64-encoded string into byte array. Strips whitespace (for testing).
func fromBase64Bytes(b64 string) []byte {
	re := regexp.MustCompile(`\s+`)
	val, err := base64.StdEncoding.DecodeString(re.ReplaceAllString(b64, ""))
	if err != nil {
		panic("Invalid test data")
	}
	return val
}

// GenerateSigningKey generates a keypair for corresponding SignatureAlgorithm.
func GenerateSigningKey(bits int) (crypto.PublicKey, crypto.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return key.Public(), key, err
}

func TestCert(t *testing.T) {

	rsaPrivateKeyPEM := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA20St6pqB4LQvqT1Aq2jZPbrkpSiwFeQwiu6AA2eBz3oYveYA
SCDzl/jXfPsY36b8VahDWmhgB/ie5Ku+R6yXiZcY9SYDiu8sMONwdkhlIL4nP1oC
97CffWf4vkt4mH7i5/rJWCd/MMLzjSmrMPdUOh9Jd2awNjUZ9QiVTBogZeMo8b5i
nVBRfRcKAQDZYlo5/VkpaRBTqahh+RoIReX1MHy/LuPMJywPaqHpIh3dlwOvnY6Q
uFrPo3cF4B7mi/ofTeRX7xzm6z+uxVZGkUHAxgm4VMAYmiP0dLSzyagA5IHUaPHV
ex8luTSR6DcbINm0bw9skUzI8zYPIGzI/rchSQIDAQABAoIBAQDazaAXOfNcvbHJ
2jvMUKZn+TXssbt1PO5L1U+dFg7tcVN7PCcP0wIBpumx6AecNtAa0fvUHc+mZKx6
V/9bGpllTYg0KajjXWPlrTAueHOhxt73UuUfMfsVc0k+66T917Cp+RIui8taZ1AO
j4QrKsO79Dilk61HipnKcLQ66t9liv4Uf/oxOjfvjaw0+mRDgD2eulTNE+pSIw6L
uZXduUcpZkYenXCIS+YfRjKMJGHdCiy0bj8887vg0JiqF+mPxGo1UrOMrkWtC4am
Fht7IMUO5KnfBveL1rMB3ed8LRie9B5EOopRoBZ7PhZ31sqlimYargHGnZwYH8BH
HzazCGwBAoGBAO8N14JcbqEcs0VpGqyuuBffheu3+6waGt90MhYEMVJsL07qLkIw
8P4zvPDthXMncrLBC7VJzKkZ7hmww3/qZX5xYjeSVggxG149I1Kncqn9l9BW/Qes
IEmTUfDE8Js6mQfJVxf7qKDsN9E5N90Oj2j4XZK2ECfaLKbwWfDv3IBBAoGBAOrP
x/jm9s6Y6KBzxBkXK0jtx2PGM1KxwJFcH9TKgz1A5yue0I1gVdU5Yf3HQowkUGJK
lT2sUHh1JXUWd2gSrZ5ba6Fc7yITIRUYjAJaW4JKvGtk59QsdRUsHiKsMxmM1GJl
/uDuZem+EiSA4R9ZZZSHAIfQY2VJD3MLDWVMvt8JAoGAJDebo/NvC1e2zVhMI0dh
OrSxrHG2Xm+iDKKlB/LgqhUb4b/W/E4/5LNf97x0kGq0lOJsbK3epOv5x8ihBds0
P0DcWYEBKcKO2+s1U8tsstZpzrWvJh9s0NjR/EFKFqp9DtHxMP/+n0rKdhdOIF6Z
WZTvUE/nCLKkOzKE3dzpMkECgYAYkkmwyCqHkAS31aVtorkK1qcIz9LLEoK+M0+5
ar+1BzepnuLgCHay62BPuCxEkgA/aOKZI5EAKfITgJhaMaotag+nQRxdCndpx7nO
/TmaNsvkyRhhYY2W+5jjs/Vc9Rm8ekPjsc7EWPl5DGuCZk507nOlwq7ECJMvTLbI
JPHMUQKBgF9O0xzJu7NwR1njqeU1MWdo8nzmb9F2itsYRXmOtC+rjTs3uqWBqlu3
TE+L0j3o3S6navSHhzzcZLwozW6otHfDcmfFBQG48zbH7YgBVuTnSQyegEpSUHRa
Pk78NMGbTCMJ65lA96vscXaSk0hF9Y83YY9Jjiju+uwWdnx74khb
-----END RSA PRIVATE KEY-----`)
	rsaPrivatePEMBlock, _ := pem.Decode(rsaPrivateKeyPEM)
	if rsaPrivatePEMBlock == nil {
		panic("failed to parse certificate PEM")
	}
	rsaPrivateKey, rsaPrivateKeyErr := x509.ParsePKCS1PrivateKey(rsaPrivatePEMBlock.Bytes)
	if rsaPrivateKeyErr != nil {
		fmt.Printf("x509.ParsePKCS1PrivateKey Error: %s\n", rsaPrivateKeyErr)
	}
	sig, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.RS256,
			Key:       rsaPrivateKey,
		},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		panic(err)
	}

	cl := jwt.Claims{
		Subject:   "subject",
		Issuer:    "issuer",
		NotBefore: jwt.NewNumericDate(time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)),
		Audience:  jwt.Audience{"leela", "fry"},
	}
	privateRaw, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	if err != nil {
		panic(err)
	}

	fmt.Println(privateRaw)
	tok, err := jose.ParseSigned(privateRaw)
	fmt.Println("======================")
	certPEM := `-----BEGIN CERTIFICATE-----
MIIDSjCCAjICCQDWXqh/wC9VZjANBgkqhkiG9w0BAQUFADBnMQswCQYDVQQGEwJD
TjERMA8GA1UECAwIU2hhbmRvbmcxDjAMBgNVBAcMBUppbmFuMQ8wDQYDVQQKDAZk
ZXZvcHMxDzANBgNVBAsMBmRldm9wczETMBEGA1UEAwwKbmlsb3JnLmNvbTAeFw0y
MDA1MTYxMjA5MjNaFw0yMTA1MTYxMjA5MjNaMGcxCzAJBgNVBAYTAkNOMREwDwYD
VQQIDAhTaGFuZG9uZzEOMAwGA1UEBwwFSmluYW4xDzANBgNVBAoMBmRldm9wczEP
MA0GA1UECwwGZGV2b3BzMRMwEQYDVQQDDApuaWxvcmcuY29tMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEA20St6pqB4LQvqT1Aq2jZPbrkpSiwFeQwiu6A
A2eBz3oYveYASCDzl/jXfPsY36b8VahDWmhgB/ie5Ku+R6yXiZcY9SYDiu8sMONw
dkhlIL4nP1oC97CffWf4vkt4mH7i5/rJWCd/MMLzjSmrMPdUOh9Jd2awNjUZ9QiV
TBogZeMo8b5inVBRfRcKAQDZYlo5/VkpaRBTqahh+RoIReX1MHy/LuPMJywPaqHp
Ih3dlwOvnY6QuFrPo3cF4B7mi/ofTeRX7xzm6z+uxVZGkUHAxgm4VMAYmiP0dLSz
yagA5IHUaPHVex8luTSR6DcbINm0bw9skUzI8zYPIGzI/rchSQIDAQABMA0GCSqG
SIb3DQEBBQUAA4IBAQAxCCdWsJjI0BNja2VhW4UjN+E2NiE5YQU0wZWtoPtc//lt
RziOGrZP82W6uh6BreonBu9JdNOJ0z+FYO957OrCrk6YBoFHe3l38KkQa13Vc4yG
2I4s1QPwor9rPRLcRQv4rB/ZS42IXXQBaCEHg+RfQ6oOX8E8YVpmRI8i3fBL4Zcf
KPiaI5i2Ey9p7ncV+7LhZ9+rZvMeA10v1jdXhl0rRphJjN+EyC+pHCu01NAaQKAo
Cj3vnvAfK8f8dEsZ9hUHLw1olVz0PbdsoUwdvULvVU5weVNyIGFfFMQeoZESrhxr
B36K98eWEdm2Wc3IY6OL2xj+DaYm8Tuyh9KzL9hU
-----END CERTIFICATE-----`
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	certs, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}
	x5tSHA1 := sha1.Sum(certs[0].Raw)
	x5tSHA256 := sha256.Sum256(certs[0].Raw)

	jwk := jose.JSONWebKey{
		Key:                         certs[0].PublicKey,
		KeyID:                       "bar",
		Algorithm:                   "RS256",
		Use:                         "sig",
		Certificates:                certs,
		CertificateThumbprintSHA1:   x5tSHA1[:],
		CertificateThumbprintSHA256: x5tSHA256[:],
	}
	if jwk.IsPublic() {
		fmt.Println("jwk.IsPublic():", true)
	}
	raw1, err := jwk.MarshalJSON()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", raw1)
	ok, okErr := tok.Verify(jwk.Public())
	if okErr != nil {
		fmt.Println("tok.Verify: ", okErr)
	}
	fmt.Printf("验证通过：%s\n", ok)
}
