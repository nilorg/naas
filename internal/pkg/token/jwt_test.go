package token

import (
	"os"
	"testing"

	"github.com/nilorg/naas/internal/module/global"
	"github.com/nilorg/oauth2"
	"github.com/nilorg/pkg/logger"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	logger.Init()
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	configFilename := "configs/config.yaml"
	if v := os.Getenv("NAAS_CONFIG"); v != "" {
		configFilename = v
	}
	viper.SetConfigFile(configFilename)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		logger.Fatalf("Fatal error config file: %s ", err)
	}
	global.Init()
	m.Run()
}

func TestJwtToken(t *testing.T) {
	logger.Debugln("TestJwtToken....")
	var (
		err           error
		tokenResponse *oauth2.TokenResponse
		claims        *oauth2.JwtClaims
	)
	gat := NewGenerateAccessToken(global.JwtPrivateKey, true)
	tokenResponse, err = gat("naas", "1001", "openid profile", "1111", nil)
	if err != nil {
		t.Error(err)
		return
	}
	pat := NewParseAccessToken(global.JwtPrivateKey)
	claims, err = pat(tokenResponse.AccessToken)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(claims)
}
