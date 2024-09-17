package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"maps"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

func createSign(secret string, body map[string]any, urlParams map[string]string) string {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	signParams := map[string]string{
		"access_key_id":   APP_KEY,
		"timestamp":       time.Now().Format("2006-01-02T15:04:05"),
		"signature_nonce": uuid.NewString(),
		"request_body":    string(bodyJson),
	}

	for k, v := range urlParams {
		signParams[k] = v
	}

	signature, err := GenerateSignature(signParams, secret)
	if err != nil {
		panic(err)
	}

	delete(signParams, "request_body")
	signParams["signature"] = url.QueryEscape(signature)

	return urlFormat(signParams)
}

// GenerateSignature 用于生成签名
func GenerateSignature(parameters map[string]string, accessKeySecret string) (string, error) {
	// 计算证书签名
	stringToSign := urlFormat(parameters)

	// 进行base64 encode
	secret := accessKeySecret + "&"
	h := hmac.New(sha1.New, []byte(secret))
	if _, err := h.Write([]byte(stringToSign)); err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return strings.TrimRight(signature, "\n"), nil
}

// urlFormat 将 map 参数转换为 url 参数格式
func urlFormat(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	var buf strings.Builder
	keys := slices.Collect(maps.Keys(params))
	slices.Sort(keys)
	for _, k := range keys {
		v := params[k]
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
	}
	return buf.String()
}
