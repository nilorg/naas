package oauth2

// ServerOptions server可选参数列表
type ServerOptions struct {
	Log                                Logger
	Issuer                             string
	DeviceAuthorizationEndpointEnabled bool   // https://tools.ietf.org/html/rfc8628
	DeviceVerificationURI              string // https://tools.ietf.org/html/rfc8628#section-3.2
	IntrospectEndpointEnabled          bool   // https://tools.ietf.org/html/rfc7662
	TokenRevocationEnabled             bool   // https://tools.ietf.org/html/rfc7009
	CustomGrantTypeEnabled             bool   // 自定义身份验证
	CustomGrantTypeAuthentication      map[string]CustomGrantTypeAuthenticationFunc
}

// ServerOption 为可选参数赋值的函数
type ServerOption func(*ServerOptions)

// ServerLogger ...
func ServerLogger(log Logger) ServerOption {
	return func(o *ServerOptions) {
		o.Log = log
	}
}

// ServerIssuer ...
func ServerIssuer(issuer string) ServerOption {
	return func(o *ServerOptions) {
		o.Issuer = issuer
	}
}

// ServerDeviceAuthorizationEndpointEnabled ...
func ServerDeviceAuthorizationEndpointEnabled(deviceAuthorizationEndpointEnabled bool) ServerOption {
	return func(o *ServerOptions) {
		o.DeviceAuthorizationEndpointEnabled = deviceAuthorizationEndpointEnabled
	}
}

// ServerDeviceVerificationURI ...
func ServerDeviceVerificationURI(deviceVerificationURI string) ServerOption {
	return func(o *ServerOptions) {
		o.DeviceVerificationURI = deviceVerificationURI
	}
}

// ServerIntrospectEndpointEnabled ...
func ServerIntrospectEndpointEnabled(introspectEndpointEnabled bool) ServerOption {
	return func(o *ServerOptions) {
		o.IntrospectEndpointEnabled = introspectEndpointEnabled
	}
}

// ServerTokenRevocationEnabled ...
func ServerTokenRevocationEnabled(tokenRevocationEnabled bool) ServerOption {
	return func(o *ServerOptions) {
		o.TokenRevocationEnabled = tokenRevocationEnabled
	}
}

// ServerCustomGrantTypeEnabled ...
func ServerCustomGrantTypeEnabled(customGrantTypeEnabled bool) ServerOption {
	return func(o *ServerOptions) {
		o.CustomGrantTypeEnabled = customGrantTypeEnabled
	}
}

// ServerCustomGrantTypeAuthentication ...
func ServerCustomGrantTypeAuthentication(customGrantTypeAuthentication map[string]CustomGrantTypeAuthenticationFunc) ServerOption {
	return func(o *ServerOptions) {
		o.CustomGrantTypeAuthentication = customGrantTypeAuthentication
	}
}

// newServerOptions 创建server可选参数
func newServerOptions(opts ...ServerOption) ServerOptions {
	opt := ServerOptions{
		Log:                                &DefaultLogger{},
		Issuer:                             DefaultJwtIssuer,
		DeviceAuthorizationEndpointEnabled: false,
		DeviceVerificationURI:              "/device",
		IntrospectEndpointEnabled:          false,
		TokenRevocationEnabled:             false,
		CustomGrantTypeEnabled:             false,
		CustomGrantTypeAuthentication:      make(map[string]CustomGrantTypeAuthenticationFunc),
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
