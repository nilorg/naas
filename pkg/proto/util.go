package proto

import (
	"encoding/json"

	"github.com/dapr/go-sdk/service/common"
	"google.golang.org/protobuf/proto"
)

// EncodeValue 对值进行编码
func EncodeValue(src interface{}) (out *common.Content, err error) {
	out = new(common.Content)
	defer func() {
		if err != nil {
			out = nil
		}
	}()
	if msg, ok := src.(proto.Message); ok {
		out.ContentType = "application/x-protobuf"
		out.Data, err = proto.Marshal(msg)
	} else if s, ok := src.(string); ok {
		out.ContentType = "text/plain"
		out.Data = []byte(s)
	} else {
		out.ContentType = "application/json"
		out.Data, err = json.Marshal(src)
	}

	return
}

// DecodeValue 对值进行解码
func DecodeValue(in *common.InvocationEvent, dist interface{}) (err error) {
	if in.ContentType == "application/json" {
		err = json.Unmarshal(in.Data, dist)
	} else if in.ContentType == "application/x-protobuf" {
		err = proto.Unmarshal(in.Data, dist.(proto.Message))
	} else if in.ContentType == "text/plain" {
		dist = string(in.Data)
	}
	return
}
