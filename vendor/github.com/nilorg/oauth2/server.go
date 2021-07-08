package oauth2

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Server OAuth2Server
type Server struct {
	VerifyClient                VerifyClientFunc
	VerifyClientID              VerifyClientIDFunc
	VerifyScope                 VerifyScopeFunc
	VerifyPassword              VerifyPasswordFunc
	VerifyRedirectURI           VerifyRedirectURIFunc
	GenerateCode                GenerateCodeFunc
	VerifyCode                  VerifyCodeFunc
	GenerateAccessToken         GenerateAccessTokenFunc
	GenerateDeviceAuthorization GenerateDeviceAuthorizationFunc
	VerifyDeviceCode            VerifyDeviceCodeFunc
	RefreshAccessToken          RefreshAccessTokenFunc
	ParseAccessToken            ParseAccessTokenFunc
	VerifyIntrospectionToken    VerifyIntrospectionTokenFunc
	TokenRevocation             TokenRevocationFunc
	opts                        ServerOptions
}

// NewServer 创建服务器
func NewServer(opts ...ServerOption) *Server {
	options := newServerOptions(opts...)
	return &Server{
		opts: options,
	}
}

// Init 初始化
func (srv *Server) Init(opts ...ServerOption) {
	for _, o := range opts {
		o(&srv.opts)
	}

	if srv.VerifyClient == nil {
		panic(ErrVerifyClientFuncNil)
	}
	if srv.VerifyClientID == nil {
		panic(ErrVerifyClientIDFuncNil)
	}
	if srv.VerifyPassword == nil {
		panic(ErrVerifyPasswordFuncNil)
	}
	if srv.VerifyRedirectURI == nil {
		panic(ErrVerifyRedirectURIFuncNil)
	}
	if srv.GenerateCode == nil {
		panic(ErrGenerateCodeFuncNil)
	}
	if srv.VerifyCode == nil {
		panic(ErrVerifyCodeFuncNil)
	}
	if srv.VerifyScope == nil {
		panic(ErrVerifyScopeFuncNil)
	}
	if srv.GenerateAccessToken == nil {
		panic(ErrGenerateAccessTokenFuncNil)
	}
	if srv.RefreshAccessToken == nil {
		panic(ErrRefreshAccessTokenFuncNil)
	}
	if srv.ParseAccessToken == nil {
		panic(ErrParseAccessTokenFuncNil)
	}

	if srv.opts.DeviceAuthorizationEndpointEnabled {
		if srv.GenerateDeviceAuthorization == nil {
			panic(ErrGenerateDeviceAuthorizationFuncNil)
		}
		if srv.VerifyDeviceCode == nil {
			panic(ErrVerifyDeviceCodeFuncNil)
		}
	}
	if srv.opts.IntrospectEndpointEnabled {
		if srv.VerifyIntrospectionToken == nil {
			panic(ErrVerifyIntrospectionTokenFuncNil)
		}
	}
	if srv.opts.TokenRevocationEnabled {
		if srv.TokenRevocation == nil {
			panic(ErrTokenRevocationFuncNil)
		}
	}
}

// HandleAuthorize 处理Authorize
func (srv *Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	// 判断参数
	responseType := r.FormValue(ResponseTypeKey)
	clientID := r.FormValue(ClientIDKey)
	scope := r.FormValue(ScopeKey)
	state := r.FormValue(StateKey)
	redirectURIStr := r.FormValue(RedirectURIKey)
	redirectURI, err := url.Parse(redirectURIStr)
	if err != nil {
		WriterError(w, ErrInvalidRequest)
		return
	}
	if responseType == "" || clientID == "" {
		RedirectError(w, r, redirectURI, ErrInvalidRequest)
		return
	}

	err = srv.VerifyRedirectURI(clientID, redirectURI.String())
	if err != nil {
		RedirectError(w, r, redirectURI, err)
		return
	}

	if err = srv.VerifyScope(StringSplit(scope, " "), clientID); err != nil {
		// ErrInvalidScope
		RedirectError(w, r, redirectURI, err)
		return
	}
	var openID string
	openID, err = OpenIDFromContext(r.Context())
	if err != nil {
		RedirectError(w, r, redirectURI, ErrServerError)
		return
	}
	switch responseType {
	case CodeKey:
		var code string
		code, err = srv.authorizeAuthorizationCode(clientID, redirectURIStr, scope, openID)
		if err != nil {
			RedirectError(w, r, redirectURI, err)
		} else {
			RedirectSuccess(w, r, redirectURI, code)
		}
	case TokenKey:
		var token *TokenResponse
		token, err = srv.authorizeImplicit(clientID, scope, openID)
		if err != nil {
			RedirectError(w, r, redirectURI, err)
		} else {
			http.Redirect(w, r, fmt.Sprintf("%s#access_token=%s&state=%s&token_type=%s&expires_in=%d", redirectURIStr, token.AccessToken, state, token.TokenType, token.ExpiresIn), http.StatusFound)
		}
	default:
		RedirectError(w, r, redirectURI, ErrUnsupportedResponseType)
	}
}

// HandleDeviceAuthorization 处理DeviceAuthorization
// https://tools.ietf.org/html/rfc8628#section-3.1
func (srv *Server) HandleDeviceAuthorization(w http.ResponseWriter, r *http.Request) {
	// 判断参数
	clientID := r.FormValue(ClientIDKey)
	scope := r.FormValue(ScopeKey)
	if clientID == "" {
		WriterError(w, ErrInvalidRequest)
		return
	}
	resp, err := srv.authorizeDeviceCode(clientID, scope)
	if err != nil {
		WriterError(w, err)
	} else {
		WriterJSON(w, resp)
	}
}

// HandleTokenIntrospection 处理内省端点
// https://tools.ietf.org/html/rfc7662#section-2.1
func (srv *Server) HandleTokenIntrospection(w http.ResponseWriter, r *http.Request) {
	var reqClientBasic *ClientBasic
	var err error
	reqClientBasic, err = RequestClientBasic(r)
	if err != nil {
		WriterError(w, err)
		return
	}
	err = srv.VerifyClient(reqClientBasic)
	if err != nil {
		WriterError(w, err)
		return
	}

	token := r.FormValue(TokenKey)
	tokenTypeHint := r.FormValue(TokenTypeHintKey)
	if tokenTypeHint != "" && tokenTypeHint != AccessTokenKey && tokenTypeHint != RefreshTokenKey {
		WriterError(w, ErrUnsupportedTokenType)
		return
	}
	var resp *IntrospectionResponse
	resp, err = srv.VerifyIntrospectionToken(token, reqClientBasic.ID, tokenTypeHint)
	if err != nil {
		WriterError(w, err)
	} else {
		WriterJSON(w, resp)
	}
}

// HandleTokenRevocation 处理Token销毁
// https://tools.ietf.org/html/rfc7009
func (srv *Server) HandleTokenRevocation(w http.ResponseWriter, r *http.Request) {
	var reqClientBasic *ClientBasic
	var err error
	reqClientBasic, err = RequestClientBasic(r)
	if err != nil {
		WriterError(w, err)
		return
	}
	err = srv.VerifyClient(reqClientBasic)
	if err != nil {
		WriterError(w, err)
		return
	}

	// 判断参数
	clientID := r.FormValue(ClientIDKey)
	if clientID == "" {
		WriterError(w, ErrInvalidRequest)
		return
	}
	if reqClientBasic.ID != clientID {
		WriterError(w, ErrInvalidRequest)
		return
	}
	token := r.FormValue(TokenKey)
	tokenTypeHint := r.FormValue(TokenTypeHintKey)
	if tokenTypeHint != "" && tokenTypeHint != AccessTokenKey && tokenTypeHint != RefreshTokenKey {
		WriterError(w, ErrUnsupportedTokenType)
		return
	}
	srv.TokenRevocation(token, clientID, tokenTypeHint)
	w.WriteHeader(http.StatusOK)
}

// HandleToken 处理Token
func (srv *Server) HandleToken(w http.ResponseWriter, r *http.Request) {

	grantType := r.PostFormValue(GrantTypeKey)
	if grantType == "" {
		WriterError(w, ErrInvalidRequest)
		return
	}

	var reqClientBasic *ClientBasic
	var err error
	var clientID string
	// explain: https://tools.ietf.org/html/rfc8628#section-3.4 {
	if grantType != DeviceCodeKey && grantType != UrnIetfParamsOAuthGrantTypeDeviceCodeKey {
		reqClientBasic, err = RequestClientBasic(r)
		if err != nil {
			WriterError(w, err)
			return
		}
		err = srv.VerifyClient(reqClientBasic)
		if err != nil {
			WriterError(w, err)
			return
		}
		clientID = reqClientBasic.ID
	} else {
		clientID = r.PostFormValue(ClientIDKey)
		err = srv.VerifyClientID(clientID)
		if err != nil {
			WriterError(w, err)
			return
		}
	}

	scope := r.PostFormValue(ScopeKey)
	if err = srv.VerifyScope(StringSplit(scope, " "), clientID); err != nil {
		// ErrInvalidScope
		WriterError(w, err)
		return
	}

	if grantType == RefreshTokenKey {
		refreshToken := r.PostFormValue(RefreshTokenKey)
		model, err := srv.RefreshAccessToken(reqClientBasic.ID, refreshToken)
		if err != nil {
			WriterError(w, err)
		} else {
			WriterJSON(w, model)
		}
	} else if grantType == AuthorizationCodeKey {
		code := r.PostFormValue(CodeKey)
		redirectURIStr := r.PostFormValue(RedirectURIKey)
		clientID := r.PostFormValue(ClientIDKey)
		if code == "" || redirectURIStr == "" || clientID == "" {
			WriterError(w, ErrInvalidRequest)
			return
		}
		var model *TokenResponse
		model, err = srv.tokenAuthorizationCode(reqClientBasic, clientID, code, redirectURIStr)
		if err != nil {
			WriterError(w, err)
		} else {
			WriterJSON(w, model)
		}
	} else if grantType == PasswordKey {
		username := r.PostFormValue(UsernameKey)
		password := r.PostFormValue(PasswordKey)
		if username == "" || password == "" {
			WriterError(w, ErrInvalidRequest)
			return
		}
		var model *TokenResponse
		model, err := srv.tokenResourceOwnerPasswordCredentials(reqClientBasic, username, password, scope)
		if err != nil {
			WriterError(w, err)
		} else {
			WriterJSON(w, model)
		}
	} else if grantType == ClientCredentialsKey {
		model, err := srv.tokenClientCredentials(reqClientBasic, scope)
		if err != nil {
			WriterError(w, err)
		} else {
			WriterJSON(w, model)
		}
	} else if grantType == UrnIetfParamsOAuthGrantTypeDeviceCodeKey || grantType == DeviceCodeKey { // https://tools.ietf.org/html/rfc8628#section-3.4
		deviceCode := r.PostFormValue(DeviceCodeKey)
		clientID := r.PostFormValue(ClientIDKey)
		model, err := srv.tokenDeviceCode(clientID, deviceCode)
		if err != nil {
			WriterError(w, err)
		} else {
			WriterJSON(w, model)
		}
	} else {
		if srv.opts.CustomGrantTypeEnabled {
			custom, ok := srv.opts.CustomGrantTypeAuthentication[grantType]
			if ok {
				model, err := srv.generateCustomGrantTypeAccessToken(reqClientBasic, scope, r, custom)
				if err != nil {
					WriterError(w, err)
				} else {
					WriterJSON(w, model)
				}
				return
			}
		}
		WriterError(w, ErrUnsupportedGrantType)
	}
}

// 授权码（authorization-code）
func (srv *Server) authorizeAuthorizationCode(clientID, redirectURI, scope, openID string) (code string, err error) {
	return srv.GenerateCode(clientID, openID, redirectURI, StringSplit(scope, " "))
}

func (srv *Server) tokenAuthorizationCode(client *ClientBasic, clientID, code, redirectURI string) (token *TokenResponse, err error) {
	if client.ID != clientID {
		err = ErrInvalidClient
		return
	}
	var value *CodeValue
	value, err = srv.VerifyCode(code, client.ID, redirectURI)
	if err != nil {
		return
	}
	scope := strings.Join(value.Scope, " ")
	token, err = srv.GenerateAccessToken(srv.opts.Issuer, client.ID, scope, value.OpenID, value)
	return
}

// 隐藏式（implicit）
func (srv *Server) authorizeImplicit(clientID, scope, openID string) (token *TokenResponse, err error) {
	token, err = srv.GenerateAccessToken(srv.opts.Issuer, clientID, scope, openID, nil)
	return
}

// 设备模式（Device Code）
func (srv *Server) authorizeDeviceCode(clientID, scope string) (resp *DeviceAuthorizationResponse, err error) {
	resp, err = srv.GenerateDeviceAuthorization(srv.opts.Issuer, srv.opts.DeviceVerificationURI, clientID, StringSplit(scope, " "))
	return
}

// 密码式（password）
func (srv *Server) tokenResourceOwnerPasswordCredentials(client *ClientBasic, username, password, scope string) (token *TokenResponse, err error) {
	var openID string
	openID, err = srv.VerifyPassword(username, password)
	if err != nil {
		return
	}
	token, err = srv.GenerateAccessToken(srv.opts.Issuer, client.ID, scope, openID, nil)
	return
}

// generateCustomGrantTypeAccessToken 生成自定义GrantType Token
func (srv *Server) generateCustomGrantTypeAccessToken(client *ClientBasic, scope string, req *http.Request, custom CustomGrantTypeAuthenticationFunc) (token *TokenResponse, err error) {
	var openID string
	openID, err = custom(client, req)
	if err != nil {
		return
	}
	token, err = srv.GenerateAccessToken(srv.opts.Issuer, client.ID, scope, openID, nil)
	return
}

// 客户端凭证（client credentials）
func (srv *Server) tokenClientCredentials(client *ClientBasic, scope string) (token *TokenResponse, err error) {
	token, err = srv.GenerateAccessToken(srv.opts.Issuer, client.ID, scope, "", nil)
	return
}

// 设备模式（Device Code）
func (srv *Server) tokenDeviceCode(clientID, deviceCode string) (token *TokenResponse, err error) {
	var value *DeviceCodeValue
	value, err = srv.VerifyDeviceCode(deviceCode, clientID)
	if err != nil {
		return
	}
	scope := strings.Join(value.Scope, " ")
	token, err = srv.GenerateAccessToken(srv.opts.Issuer, clientID, scope, value.OpenID, nil)
	return
}
