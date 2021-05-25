package wellknown

import (
	"net/http"
	"strings"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// GetOpenIDProviderMetadata ...
func GetOpenIDProviderMetadata(ctx *gin.Context) {
	issuer := strings.TrimSuffix(viper.GetString("server.oauth2.issuer"), "/")
	metadata := OpenIDProviderMetadata{
		Issuer:                issuer,
		AuthorizationEndpoint: issuer + "/oauth2/authorize",
		TokenEndpoint:         issuer + "/oauth2/token",
		JwksURI:               issuer + "/.well-known/jwks.json",
	}
	if viper.GetBool("server.oauth2.device_authorization_endpoint_enabled") {
		metadata.DeviceAuthorizationEndpoint = issuer + "/oauth2/device/code"
	}
	if viper.GetBool("server.oauth2.introspection_endpoint_enabled") {
		metadata.IntrospectionEndpoint = issuer + "/oauth2/introspect"
	}
	if viper.GetBool("server.oauth2.revocation_endpoint_enabled") {
		metadata.RevocationEndpoint = issuer + "/oauth2/revoke"
	}
	if viper.GetBool("server.oidc.enabled") && viper.GetBool("server.oidc.userinfo_endpoint_enabled") {
		metadata.UserinfoEndpoint = issuer + "/oidc/userinfo"
	}
	metadata.ResponseTypesSupported = append(metadata.ResponseTypesSupported,
		"code",
		"token",
		"id_token",
		"code token",
		"code id_token",
		"token id_token",
		"code token id_token",
		"none",
	)
	metadata.SubjectTypesSupported = append(metadata.SubjectTypesSupported,
		"public",
	)
	metadata.IDTokenSigningAlgValuesSupported = append(metadata.IDTokenSigningAlgValuesSupported,
		"RS256",
	)
	if scopeCodes, err := service.OAuth2.AllScopeCode(contexts.WithGinContext(ctx)); err == nil {
		metadata.ScopesSupported = append(metadata.ScopesSupported, model.ConvertCodeSliceToStringSlice(scopeCodes)...)
	}
	metadata.TokenEndpointAuthMethodsSupported = append(metadata.TokenEndpointAuthMethodsSupported,
		"client_secret_post",
		"client_secret_basic",
	)
	metadata.ClaimsSupported = append(metadata.ClaimsSupported,
		"sub",
		"aud",
		"email",
		"email_verified",
		"exp",
		"nickname",
		"iat",
		"iss",
		"name",
		"picture",
	)
	metadata.CodeChallengeMethodsSupported = append(metadata.CodeChallengeMethodsSupported,
		"plain",
		"S256",
	)
	metadata.GrantTypesSupported = append(metadata.GrantTypesSupported,
		"authorization_code",
		"implicit",
		"password",
		"client_credentials",
		"refresh_token",
		"urn:ietf:params:oauth:grant-type:device_code",
		"urn:ietf:params:oauth:grant-type:jwt-bearer",
	)
	ctx.JSON(http.StatusOK, metadata)
}
