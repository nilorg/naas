package wechat

import (
	"errors"
)

// MetadataAccessTokenKey Metadata wehcat AccessToken key.
const MetadataAccessTokenKey = "wechat-access-token"

// MetadataJsAPITicketKey Metadata wehcat JsAPITicketKey key.
const MetadataJsAPITicketKey = "wechat-js-api-ticket"

var (
	// ErrMetadataNotFoundClient 元数据不存在客户端AccessToken错误
	ErrMetadataNotFoundClient = errors.New("Metadata不存在客户端AccessToken")
)

type metadataClient struct {
	accessToken string
	jsAPITicket string
}

func (m *metadataClient) GetAccessToken() string {
	return m.accessToken
}

func (m *metadataClient) GetJsAPITicket() string {
	return m.jsAPITicket
}

// FromMetadata 从元数据中获取微信客户端
func FromMetadata(metadata map[string]string) (Clienter, error) {
	accessToken, ok := metadata[MetadataAccessTokenKey]
	if !ok {
		return nil, ErrMetadataNotFoundClient
	}
	jsAPITicket, _ := metadata[MetadataJsAPITicketKey]
	return &metadataClient{
		accessToken: accessToken,
		jsAPITicket: jsAPITicket,
	}, nil
}
