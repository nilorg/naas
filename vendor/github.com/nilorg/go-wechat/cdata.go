package wechat

import "encoding/xml"

// https://segmentfault.com/q/1010000008231680

// CDATA ...
type CDATA string

// MarshalXML ...
func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		string `xml:",cdata"`
	}{string(c)}, start)
}

// String ...
func (c CDATA) String() string {
	return string(c)
}
