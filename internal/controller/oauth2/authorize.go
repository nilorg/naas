package oauth2

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/tools"
	"github.com/nilorg/naas/pkg/tools/key"
	"github.com/nilorg/oauth2"
	"github.com/nilorg/pkg/logger"
	"github.com/nilorg/sdk/convert"
	sdkStrings "github.com/nilorg/sdk/strings"
)

// AuthorizePage 授权页面
func AuthorizePage(ctx *gin.Context) {
	errMsg := GetErrorMessage(ctx)
	if errMsg != "" {
		ctx.HTML(http.StatusOK, "authorize.tmpl", gin.H{
			"error": errMsg,
		})
		return
	}

	clientID := ctx.Query("client_id")
	var (
		err        error
		client     *model.OAuth2Client
		clientInfo *model.OAuth2ClientInfo
	)
	client, err = service.OAuth2.GetClient(convert.ToUint64(clientID))
	if err != nil {
		err = SetErrorMessage(ctx, err.Error())
		if err != nil {
			logger.Errorln(err)
		}
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	// uri := *ctx.Request.URL
	// query := uri.Query()
	query := ctx.Request.URL.Query()
	queryRedirectURI := query.Get(oauth2.RedirectURIKey)
	if queryRedirectURI == "" {
		// query.Set(oauth2.RedirectURIKey, client.RedirectURI)
		// uri.RawQuery = query.Encode()
		queryRedirectURI = client.RedirectURI
	}
	// 判断重定向URL存在数据库中的前缀
	if !strings.HasPrefix(queryRedirectURI, client.RedirectURI) && queryRedirectURI != client.RedirectURI {
		err = SetErrorMessage(ctx, "重定向域名不符合后台配置规范")
		if err != nil {
			logger.Errorln(err)
		}
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	clientInfo, err = service.OAuth2.GetClientInfo(convert.ToUint64(clientID))
	if err != nil {
		ctx.HTML(http.StatusOK, "authorize.tmpl", gin.H{
			"error": err.Error(),
		})
		return
	}
	session := sessions.Default(ctx)
	currentAccount := session.Get(key.SessionAccount)
	cu := currentAccount.(*model.SessionAccount)
	// query scope checked scopes list.
	scope := ctx.Query("scope")
	scopeSplit := sdkStrings.Split(scope, " ")
	scopes := make([]map[string]interface{}, 0)
	var scopeInfos []*service.OAuth2ClientScopeInfo
	scopeInfos, err = service.OAuth2.GetClientAllScopeInfo(clientInfo.ClientID)
	if err != nil {
		ctx.HTML(http.StatusOK, "authorize.tmpl", gin.H{
			"error": err.Error(),
		})
		return
	}
	for _, v := range scopeInfos {
		scopes = append(scopes, map[string]interface{}{
			"info":    v,
			"checked": tools.InStringSplit(v.Code, scopeSplit),
		})
	}
	logBackInURI, _ := url.Parse("/oauth2/login")
	logBackInURIQuery := url.Values{}
	logBackInURIQuery.Set("client_id", clientID)
	logBackInURIQuery.Set("login_redirect_uri", ctx.Request.RequestURI)
	logBackInURI.RawQuery = logBackInURIQuery.Encode()
	ctx.HTML(http.StatusOK, "authorize.tmpl", gin.H{
		"client_info":  clientInfo,
		"scopes":       scopes,
		"current_user": cu,
		"log_back_in":  logBackInURI.String(),
	})
	return
}

func formScopeValue(r *http.Request) string {
	if r.Form == nil {
		var defaultMaxMemory int64 = 32 << 20 // 32 MB
		r.ParseMultipartForm(defaultMaxMemory)
	}
	if vs := r.Form["scope"]; len(vs) > 0 {
		return strings.Join(r.Form["scope"], " ")
	}
	return ""
}

// Authorize authorize post
func Authorize(ctx *gin.Context) {
	session := sessions.Default(ctx)
	currentAccount := session.Get(key.SessionAccount)
	cu := currentAccount.(*model.SessionAccount)
	scope := formScopeValue(ctx.Request)
	ctx.Request.Form.Set("scope", scope)
	rctx := oauth2.NewOpenIDContext(ctx.Request.Context(), fmt.Sprint(cu.UserID))
	req := ctx.Request.WithContext(rctx)
	// 模拟请求客户端
	oauth2Server.HandleAuthorize(ctx.Writer, req)
}
