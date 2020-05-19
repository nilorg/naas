package oauth2

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/tools/key"
	"github.com/nilorg/oauth2"
	"github.com/nilorg/pkg/logger"
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
	var err error
	var client *model.OAuth2Client
	client, err = service.OAuth2.GetClient(clientID)
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
	var clientInfo *model.OAuth2ClientInfo
	clientInfo, err = service.OAuth2.GetClientInfo(clientID)
	// query scope checked scopes list.
	scope := ctx.Query("scope")
	scopes := make([]map[string]interface{}, 0)
	for _, value := range SourceScope {
		scopes = append(scopes, map[string]interface{}{
			"text":    value,
			"checked": value == scope,
		})
	}
	ctx.HTML(http.StatusOK, "authorize.tmpl", gin.H{
		"error":       err,
		"client_info": clientInfo,
		"scopes":      scopes,
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
