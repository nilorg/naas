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

	"github.com/nilorg/oauth2"
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

// Structured version of Claims Section, as referenced at
// https://tools.ietf.org/html/rfc7519#section-4.1
// See examples for how to use this with your own claim types
type StandardClaims struct {
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	ID        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
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

	// cl := jwt.Claims{
	// 	Subject:   "subject",
	// 	Issuer:    "http://localhost:8080",
	// 	NotBefore: jwt.NewNumericDate(time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)),
	// 	Audience:  jwt.Audience{"naas-oidc-test", "1xxxx000"},
	// 	Expiry:    jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	// }
	cl := StandardClaims{
		Subject:   "subject",
		Issuer:    "http://localhost:8080",
		NotBefore: time.Now().Unix(),
		Audience:  "naas-oidc-test",
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}
	privateRaw, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	if err != nil {
		panic(err)
	}

	fmt.Printf("token: %s\n", privateRaw)
	tok, err := jose.ParseSigned(privateRaw)
	fmt.Println("======================")
	certPEM := `-----BEGIN CERTIFICATE-----
MIIDZjCCAk4CCQDPAI6eEEcsfDANBgkqhkiG9w0BAQUFADB0MQswCQYDVQQGEwJD
TjERMA8GA1UECAwIU2hhbmRvbmcxDjAMBgNVBAcMBUppbmFuMQ8wDQYDVQQKDAZk
ZXZvcHMxDzANBgNVBAsMBmRldm9wczEgMB4GA1UEAwwXYWNjb3VudHMuZGlhbmZl
bmc1OC5jb20wIBcNMjAwNTE3MTIxODQxWhgPMzAxOTA5MTgxMjE4NDFaMHQxCzAJ
BgNVBAYTAkNOMREwDwYDVQQIDAhTaGFuZG9uZzEOMAwGA1UEBwwFSmluYW4xDzAN
BgNVBAoMBmRldm9wczEPMA0GA1UECwwGZGV2b3BzMSAwHgYDVQQDDBdhY2NvdW50
cy5kaWFuZmVuZzU4LmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AKc8v+6nheA3aQB3Cl6eHbdOT+kzRVDmmh+mKO0PPWkoyXKP+2clRpSf6itqOg7H
pJVeE6xaQfdXAZsvIRf970C/bjre0m5zknlUqWLdFB0jH7XbmfYX88j22NF2wEJs
BiMDAxkKr3Vhb/7lVjw3o1PNMcaH8eUExWRP9fc7dwbMFiIw+8DTiUdXrgfcRqSD
pr+cyLCucNdFzRFYranDXaCNrVIgFuBBRoG/whi5/qqEAhJ4dtOqDSlZEJzbkLYR
ufTjXX15T6w8HtyRV5r7PZXemTv55n2xm7UmMpzQKvEp0TW4ZageZ92oyPDOApux
CuJuoe4++L0u1qm/WDZpyDUCAwEAATANBgkqhkiG9w0BAQUFAAOCAQEAK2IE2mBd
Y+JNzzdPqSPMHUuA8SYrgW3pm04rPzOM3KKTwc9hPuQoUeZBg2YAkH5GOi5CPHT/
qT7VJ3o0svN4a12Sd8muqPEbCR3LMbjst6rW192rvwo90+8wPc2eicFCK7vzCqIN
NRwZV784rMz9qMKzXxI36t6OGawsXn8BWKajh2ewxLr5WLSdrdYAl+7BqQLYV3yu
pcMtcIOnwg610vBBr4CsYl/lj5JnNwgLb2J23iK9zPF86iHmk+HzQy2SyqX/u1It
H402+bAq3Pp163qUjQYlfNG9ttB8xfkK53hU9SRgCQanNbCY8s9P/p604ysypkcG
UAmLvAHgYEWfSg==
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

func TestOAuth2Token(t *testing.T) {
	rsaPrivateKeyPEM := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEApzy/7qeF4DdpAHcKXp4dt05P6TNFUOaaH6Yo7Q89aSjJco/7
ZyVGlJ/qK2o6DseklV4TrFpB91cBmy8hF/3vQL9uOt7SbnOSeVSpYt0UHSMftduZ
9hfzyPbY0XbAQmwGIwMDGQqvdWFv/uVWPDejU80xxofx5QTFZE/19zt3BswWIjD7
wNOJR1euB9xGpIOmv5zIsK5w10XNEVitqcNdoI2tUiAW4EFGgb/CGLn+qoQCEnh2
06oNKVkQnNuQthG59ONdfXlPrDwe3JFXmvs9ld6ZO/nmfbGbtSYynNAq8SnRNbhl
qB5n3ajI8M4Cm7EK4m6h7j74vS7Wqb9YNmnINQIDAQABAoIBAFYP3z2znCN8oF6K
5B05BVXVyS3bIqq1YU80NQ95rkK1qKV6DwhPmHjXqqxY6DO+7aWoWjtx30yny73O
jRtJpJwPZ2yISoZol1I1DU5BMx6jeqgdsKeCQASFc6Knl90WtjnCTQ/P/edME1R7
NNucTkLL7/eY8hTHVcV/mLZ4NZKbEybuZ5KxuoubqXqkUnL4RrZCbA/Q3+N6UvFN
w4aUd9zoPfATMwU7VL5pI4madJpQjqukcV3cna59fBiFenskoEOj7UrTtcZv8h/P
jQ3Hsg5c9L7PovUlBnmltxSoPBM0e066afG3sgkKpRGTNN8sj8mYQhd8M+2wLEdv
CEd1UUECgYEA1I3btTrhJBaqcQcwXKlYks1SAUGVnxnK+bTqTxZ1iEldEalunjmb
1jmVPD0Z2wRGa/Usis7O0B+XCNVDgpGbL2C89SIs/hIbqHPmBShQg/O+zpryXZvs
Kmtvxwpptjwwmnw9Q+Pzm5I7ruK8nYX4ogZPQsVGvT9bPM0eh5zMOUkCgYEAyWui
wZ50XfQpuxR0QGXPeig8w66y7PEvaub7+RoVHPhWT8f1M4sDi/1PBM/3agnhhTN7
pNUJLhtMat35Z4ojWmBhfFROowxCzE5FinyW9eML3Yvo1xL9HWBkTLxxKxE9ilFq
jIjPwinHkMUez8Zy15rzD241bp9caffeUOCPY40CgYAIf1NVP3FYu/88XYk1ax+7
XrH0kuakYaeXq//iAYfZVvV9i0R81tjAC7VHnzm1Y8pc7oRFWFc0Qs8K71uvkJqf
nkJvmloqHhc0+M0tT5tIayopoFAoJd+fIoRpdKUdP/LBek4ItMg8Y/A24aGguoZi
E9Z/WNunHS1MlPavfTk84QKBgHGWD3ycvQbW0Em96Sj/wRckZc/8Ts6r3I+unt4F
RW7G5PWsz6w3ctKZENyn4uCbneAd/lYgBUNJBbkmYKVxEyq+O3t/l7D/ExRf93t3
czJKzcAsTCwteyv71dQoWLFu0YOVEj8aT/8wzGfpocyOHulTakqDXgJ6QAVKUMbP
PE1pAoGAJvplll85jDa4TbH6kR/iCEyPHH4x19cUdf/50iS20gP8kK1rXz3R2Y6o
l/7rd1IUFq+7PvSm2y/quRbNDx/GduoTUjxwWJrDt6/x5KnD5l6B+iwJ0kMMfkCx
TG7xfFjm50rCcibX7oQEhlZgavYHvnhTCqb1cYKfOtvd6Gv87qY=
-----END RSA PRIVATE KEY-----`)
	rsaPrivatePEMBlock, _ := pem.Decode(rsaPrivateKeyPEM)
	if rsaPrivatePEMBlock == nil {
		panic("failed to parse certificate PEM")
	}
	rsaPrivateKey, rsaPrivateKeyErr := x509.ParsePKCS1PrivateKey(rsaPrivatePEMBlock.Bytes)
	if rsaPrivateKeyErr != nil {
		fmt.Printf("x509.ParsePKCS1PrivateKey Error: %s\n", rsaPrivateKeyErr)
	}
	cl := oauth2.JwtClaims{
		JwtStandardClaims: oauth2.JwtStandardClaims{
			Subject:   "subject",
			Issuer:    "http://localhost:8080",
			NotBefore: time.Now().Unix(),
			Audience:  []string{"1001"},
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token, tokenErr := oauth2.NewJwtToken(&cl, "RS256", rsaPrivateKey)
	if tokenErr != nil {
		t.Error(tokenErr)
		return
	}
	t.Logf("rsa token: %s\n", token)
	t.Log("===================")
	cl2, cl2Err := oauth2.ParseJwtClaimsToken(token, rsaPrivateKey.Public())
	if cl2Err != nil {
		t.Error(cl2Err)
		return
	}
	t.Logf("rsa token cl2: %+v\n", cl2)
}
