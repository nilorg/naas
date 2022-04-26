package http

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Post send post request.
func Post(url, contentType string, args map[string]string) (result string, err error) {
	if contentType == "" {
		contentType = "application/x-www-form-urlencoded"
	}
	resp, err := http.Post(url, contentType, strings.NewReader(argsEncode(args)))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = string(body)
	return
}

// Get send get request.
func Get(url string, args map[string]string) (result string, err error) {
	if args != nil {
		url += "?" + argsEncode(args)
	}
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = string(body)
	return
}

func argsEncode(params map[string]string) string {
	args := url.Values{}
	for k, v := range params {
		args.Set(k, v)
	}
	return args.Encode()
}
