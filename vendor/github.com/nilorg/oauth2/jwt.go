package oauth2

import (
	"crypto/subtle"
	"fmt"
	"time"

	"github.com/nilorg/pkg/slice"
	"github.com/nilorg/sdk/strings"
	jose "github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
)

// 参考 https://github.com/dgrijalva/jwt-go/blob/master/claims.go

// ----- helpers

func verifyAud(aud []string, cmp []string, required bool) bool {
	if len(aud) == 0 {
		return !required
	}
	return slice.IsSubset(cmp, aud)
}

func verifyExp(exp int64, now int64, required bool) bool {
	if exp == 0 {
		return !required
	}
	return now <= exp
}

func verifyIat(iat int64, now int64, required bool) bool {
	if iat == 0 {
		return !required
	}
	return now >= iat
}

func verifyIss(iss string, cmp string, required bool) bool {
	if iss == "" {
		return !required
	}
	return subtle.ConstantTimeCompare([]byte(iss), []byte(cmp)) != 0
}

func verifyNbf(nbf int64, now int64, required bool) bool {
	if nbf == 0 {
		return !required
	}
	return now >= nbf
}

func verifyScope(scope []string, cmp []string, required bool) bool {
	if len(scope) == 0 {
		return !required
	}
	return slice.IsSubset(cmp, scope)
}

// JwtStandardClaims as referenced at
// https://tools.ietf.org/html/rfc7519#section-4.1
type JwtStandardClaims struct {
	Audience  []string `json:"aud,omitempty"`
	ExpiresAt int64    `json:"exp,omitempty"`
	ID        string   `json:"jti,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
	Issuer    string   `json:"iss,omitempty"`
	NotBefore int64    `json:"nbf,omitempty"`
	Subject   string   `json:"sub,omitempty"`
}

// Valid time based claims "exp, iat, nbf".
// There is no accounting for clock skew.
// As well, if any of the above claims are not in the token, it will still
// be considered a valid claim.
func (c JwtStandardClaims) Valid() error {
	now := time.Now().Unix()
	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if !c.VerifyExpiresAt(now, false) {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		return fmt.Errorf("token is expired by %v", delta)
	}
	if !c.VerifyIssuedAt(now, false) {
		return fmt.Errorf("Token used before issued")
	}
	if !c.VerifyNotBefore(now, false) {
		return fmt.Errorf("token is not valid yet")
	}
	return nil
}

// VerifyAudience Compares the aud claim against cmp.
// If required is false, this method will return true if the value matches or is unset
// 如果required为false，如果值匹配或未设置，此方法将返回true
func (c *JwtStandardClaims) VerifyAudience(cmp []string, req bool) bool {
	return verifyAud(c.Audience, cmp, req)
}

// VerifyExpiresAt Compares the exp claim against cmp.
// If required is false, this method will return true if the value matches or is unset
// 如果required为false，如果值匹配或未设置，此方法将返回true
func (c *JwtStandardClaims) VerifyExpiresAt(cmp int64, req bool) bool {
	return verifyExp(c.ExpiresAt, cmp, req)
}

// VerifyIssuedAt Compares the iat claim against cmp.
// If required is false, this method will return true if the value matches or is unset
// 如果required为false，如果值匹配或未设置，此方法将返回true
func (c *JwtStandardClaims) VerifyIssuedAt(cmp int64, req bool) bool {
	return verifyIat(c.IssuedAt, cmp, req)
}

// VerifyIssuer Compares the iss claim against cmp.
// If required is false, this method will return true if the value matches or is unset
// 如果required为false，如果值匹配或未设置，此方法将返回true
func (c *JwtStandardClaims) VerifyIssuer(cmp string, req bool) bool {
	return verifyIss(c.Issuer, cmp, req)
}

// VerifyNotBefore Compares the nbf claim against cmp.
// If required is false, this method will return true if the value matches or is unset
// 如果required为false，如果值匹配或未设置，此方法将返回true
func (c *JwtStandardClaims) VerifyNotBefore(cmp int64, req bool) bool {
	return verifyNbf(c.NotBefore, cmp, req)
}

// JwtClaims 在jwt标准上的扩展
type JwtClaims struct {
	JwtStandardClaims
	Scope string `json:"scope,omitempty"`
}

// VerifyScope Compares the aud claim against cmp.
// If required is false, this method will return true if the value matches or is unset
// 如果required为false，如果值匹配或未设置，此方法将返回true
func (c *JwtClaims) VerifyScope(scope string, req bool) bool {
	source := strings.Split(c.Scope, " ")
	array := strings.Split(scope, " ")
	return verifyScope(source, array, req)
}

// NewJwtClaims ...
func NewJwtClaims(issuer, audience, scope, openID string) *JwtClaims {
	currTime := time.Now()
	return &JwtClaims{
		JwtStandardClaims: JwtStandardClaims{
			// Issuer = iss,令牌颁发者。它表示该令牌是由谁创建的
			Issuer: issuer,
			// Subject = sub,令牌的主体。它表示该令牌是关于谁的
			Subject: openID,
			// Audience = aud,令牌的受众。它表示令牌的接收者
			Audience: []string{audience},
			// ExpiresAt = exp,令牌的过期时间戳。它表示令牌将在何时过期
			ExpiresAt: currTime.Add(AccessTokenExpire).Unix(),
			// NotBefore = nbf,令牌的生效时的时间戳。它表示令牌从什么时候开始生效
			NotBefore: currTime.Unix(),
			// IssuedAt = iat,令牌颁发时的时间戳。它表示令牌是何时被创建的
			IssuedAt: currTime.Unix(),
		},
		Scope: scope,
	}
}

func newJwtToken(v interface{}, algorithm string, key interface{}) (token string, err error) {
	var sig jose.Signer
	sig, err = jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.SignatureAlgorithm(algorithm),
			Key:       key,
		},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return
	}
	token, err = jwt.Signed(sig).Claims(v).CompactSerialize()
	return
}

// NewJwtToken ...
func NewJwtToken(v interface{}, algorithm string, key interface{}) (string, error) {
	return newJwtToken(v, algorithm, key)
}

// NewJwtClaimsToken ...
func NewJwtClaimsToken(claims *JwtClaims, algorithm string, key interface{}) (string, error) {
	return newJwtToken(claims, algorithm, key)
}

// NewJwtStandardClaimsToken ...
func NewJwtStandardClaimsToken(claims *JwtStandardClaims, algorithm string, key interface{}) (string, error) {
	return newJwtToken(claims, algorithm, key)
}

// NewHS256JwtClaimsToken ...
func NewHS256JwtClaimsToken(claims *JwtClaims, jwtVerifyKey []byte) (string, error) {
	return newJwtToken(claims, "HS256", jwtVerifyKey)
}

// parseJwtToken ...
func parseJwtToken(token string, key interface{}, dest ...interface{}) (err error) {
	var (
		tok *jwt.JSONWebToken
	)
	tok, err = jwt.ParseSigned(token)
	if err != nil {
		return
	}
	err = tok.Claims(key, dest...)
	return
}

// ParseJwtClaimsToken ...
func ParseJwtClaimsToken(token string, key interface{}) (claims *JwtClaims, err error) {
	claims = new(JwtClaims)
	err = parseJwtToken(token, key, claims)
	return
}

// ParseJwtStandardClaimsToken ...
func ParseJwtStandardClaimsToken(token string, key interface{}) (claims *JwtStandardClaims, err error) {
	claims = new(JwtStandardClaims)
	err = parseJwtToken(token, key, claims)
	return
}

// ParseHS256JwtClaimsToken ...
func ParseHS256JwtClaimsToken(token string, jwtVerifyKey []byte) (claims *JwtClaims, err error) {
	return ParseJwtClaimsToken(token, jwtVerifyKey)
}
