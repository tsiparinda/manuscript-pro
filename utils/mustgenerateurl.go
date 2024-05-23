package utils

import (
	"github.com/vedicsociety/platform/http/handling"
	"net/url"
)

func MustGenerateUrl(gen handling.URLGenerator, target interface{}, data ...interface{}) string {
	url, err := gen.GenerateUrl(target, data...)
	if err != nil {
		panic(err)
	}
	return url
}

func GenerateVerifiedSignup(gen handling.URLGenerator, code string) string {
	url := "/auth/verifiedsignup/signup/" + code
	return url
}

func GenerateMsgUrl(gen handling.URLGenerator, message string) string {
	url := "/auth/message/msg"
	url = AddQueryParams(url, "message", message)
	return url
}

func AddQueryParams(rawURL string, key string, value string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	params := parsed.Query()
	params.Add(key, value)
	parsed.RawQuery = params.Encode()

	return parsed.String()
}
