package oauth2

import (
	"errors"
	"net/http"
)

var (
	// ErrInvalidRequest 无效的请求
	ErrInvalidRequest = errors.New("invalid_request")
	// ErrUnauthorizedClient 未经授权的客户端
	ErrUnauthorizedClient = errors.New("unauthorized_client")
	// ErrAccessDenied 拒绝访问
	ErrAccessDenied = errors.New("access_denied")
	// ErrUnsupportedResponseType 不支持的response类型
	ErrUnsupportedResponseType = errors.New("unsupported_response_type")
	// ErrUnsupportedGrantType 不支持的grant类型
	ErrUnsupportedGrantType = errors.New("unsupported_grant_type")
	// ErrInvalidGrant 无效的grant
	ErrInvalidGrant = errors.New("invalid_grant")
	// ErrInvalidScope 无效scope
	ErrInvalidScope = errors.New("invalid_scope")
	// ErrTemporarilyUnavailable 暂时不可用
	ErrTemporarilyUnavailable = errors.New("temporarily_unavailable")
	// ErrServerError 服务器错误
	ErrServerError = errors.New("server_error")
	// ErrInvalidClient 无效的客户
	ErrInvalidClient = errors.New("invalid_client")
	// ErrExpiredToken 过期的令牌
	ErrExpiredToken = errors.New("expired_token")
	// ErrAuthorizationPending 授权待定
	// https://tools.ietf.org/html/rfc8628#section-3.5
	ErrAuthorizationPending = errors.New("authorization_pending")
	// ErrSlowDown 轮询太频繁
	// https://tools.ietf.org/html/rfc8628#section-3.5
	ErrSlowDown = errors.New("slow_down")
	// ErrUnsupportedTokenType 不支持的令牌类型
	// https://tools.ietf.org/html/rfc7009#section-4.1.1
	ErrUnsupportedTokenType = errors.New("unsupported_token_type")
)

var (
	// ErrVerifyClientFuncNil ...
	ErrVerifyClientFuncNil = errors.New("OAuth2 Server VerifyClient Is Nil")
	// ErrVerifyClientIDFuncNil ...
	ErrVerifyClientIDFuncNil = errors.New("OAuth2 Server VerifyClientID Is Nil")
	// ErrVerifyPasswordFuncNil ...
	ErrVerifyPasswordFuncNil = errors.New("OAuth2 Server VerifyPassword Is Nil")
	// ErrVerifyRedirectURIFuncNil ...
	ErrVerifyRedirectURIFuncNil = errors.New("OAuth2 Server VerifyRedirectURI Is Nil")
	// ErrGenerateCodeFuncNil ...
	ErrGenerateCodeFuncNil = errors.New("OAuth2 Server GenerateCode Is Nil")
	// ErrVerifyCodeFuncNil ...
	ErrVerifyCodeFuncNil = errors.New("OAuth2 Server VerifyCode Is Nil")
	// ErrVerifyScopeFuncNil ...
	ErrVerifyScopeFuncNil = errors.New("OAuth2 Server VerifyScope Is Nil")
	// ErrGenerateAccessTokenFuncNil ...
	ErrGenerateAccessTokenFuncNil = errors.New("OAuth2 Server GenerateAccessTokenFunc Is Nil")
	// ErrGenerateDeviceAuthorizationFuncNil ...
	ErrGenerateDeviceAuthorizationFuncNil = errors.New("OAuth2 Server GenerateDeviceAuthorizationFunc Is Nil")
	// ErrVerifyDeviceCodeFuncNil ...
	ErrVerifyDeviceCodeFuncNil = errors.New("OAuth2 Server ErrVerifyDeviceCodeFunc Is Nil")
	// ErrRefreshAccessTokenFuncNil ...
	ErrRefreshAccessTokenFuncNil = errors.New("OAuth2 Server ErrRefreshAccessTokenFuncNil Is Nil")
	// ErrParseAccessTokenFuncNil ...
	ErrParseAccessTokenFuncNil = errors.New("OAuth2 Server ParseAccessTokenFunc Is Nil")
	// ErrVerifyIntrospectionTokenFuncNil ...
	ErrVerifyIntrospectionTokenFuncNil = errors.New("OAuth2 Server VerifyIntrospectionToken Is Nil")
	// ErrTokenRevocationFuncNil ...
	ErrTokenRevocationFuncNil = errors.New("OAuth2 Server TokenRevocation Is Nil")
	// ErrInvalidAccessToken 无效的访问令牌
	ErrInvalidAccessToken = errors.New("invalid_access_token")
	// ErrInvalidRedirectURI 无效的RedirectURI
	ErrInvalidRedirectURI = errors.New("invalid_redirect_uri")
	// ErrStateValueDidNotMatch ...
	ErrStateValueDidNotMatch = errors.New("state value did not match")
	// ErrMissingAccessToken ...
	ErrMissingAccessToken = errors.New("missing access token")
)

var (
	// Errors ...
	Errors = map[string]error{
		ErrVerifyClientFuncNil.Error():   ErrVerifyClientFuncNil,
		ErrInvalidAccessToken.Error():    ErrInvalidAccessToken,
		ErrStateValueDidNotMatch.Error(): ErrStateValueDidNotMatch,
		ErrMissingAccessToken.Error():    ErrMissingAccessToken,

		ErrInvalidRequest.Error():          ErrInvalidRequest,
		ErrUnauthorizedClient.Error():      ErrUnauthorizedClient,
		ErrAccessDenied.Error():            ErrAccessDenied,
		ErrUnsupportedResponseType.Error(): ErrUnsupportedResponseType,
		ErrUnsupportedGrantType.Error():    ErrUnsupportedGrantType,
		ErrInvalidGrant.Error():            ErrInvalidGrant,
		ErrInvalidScope.Error():            ErrInvalidScope,
		ErrTemporarilyUnavailable.Error():  ErrTemporarilyUnavailable,
		ErrServerError.Error():             ErrServerError,
		ErrInvalidClient.Error():           ErrInvalidClient,
		ErrExpiredToken.Error():            ErrExpiredToken,
		ErrAuthorizationPending.Error():    ErrAuthorizationPending,
		ErrSlowDown.Error():                ErrSlowDown,
		ErrUnsupportedTokenType.Error():    ErrUnsupportedTokenType,
	}
	// ErrStatusCodes ...
	ErrStatusCodes = map[error]int{
		ErrInvalidRequest:          http.StatusBadRequest,           // 400
		ErrUnauthorizedClient:      http.StatusUnauthorized,         // 401
		ErrAccessDenied:            http.StatusForbidden,            // 403
		ErrUnsupportedResponseType: http.StatusUnauthorized,         // 401
		ErrInvalidScope:            http.StatusBadRequest,           // 400
		ErrServerError:             http.StatusInternalServerError,  // 400
		ErrTemporarilyUnavailable:  http.StatusServiceUnavailable,   // 503
		ErrInvalidClient:           http.StatusUnauthorized,         // 401
		ErrInvalidGrant:            http.StatusUnauthorized,         // 401
		ErrUnsupportedGrantType:    http.StatusUnauthorized,         // 401
		ErrExpiredToken:            http.StatusUnauthorized,         // 401
		ErrAuthorizationPending:    http.StatusPreconditionRequired, // 428
		ErrSlowDown:                http.StatusForbidden,            // 403 https://tools.ietf.org/html/rfc6749#section-5.2
		ErrUnsupportedTokenType:    http.StatusServiceUnavailable,   // 503 https://tools.ietf.org/html/rfc7009#section-2.2.1
	}
)
