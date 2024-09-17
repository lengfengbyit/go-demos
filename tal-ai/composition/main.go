package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

// 手写识别

const (
	API_URL = "http://openai.100tal.com/aiimage/cn-composition"
	IMG_URL = "../imgs/pic_fangge.jpg"
)

var (
	APP_KEY    = os.Getenv("TAL_OPENAI_KEY")
	APP_SECRET = os.Getenv("TAL_OPENAI_SECRET")
)

func main() {
	if APP_KEY == "" || APP_SECRET == "" {
		log.Fatal("APP_KEY or APP_SECRET is empty")
	}
	params := map[string]any{
		//"image_url":           IMG_URL,
		"image_base64":        getBase64Image(),
		"details":             false,
		"paragraph_detection": true,
	}
	urlParams := createSign(APP_SECRET, params, nil)

	body := httpPost(API_URL+"?"+urlParams, params)
	_ = os.WriteFile("result.json", body, 0644)
}

func getBase64Image() string {
	content, err := os.ReadFile(IMG_URL)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(content)
}

func httpPost(url string, params map[string]any) []byte {
	paramsJson, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(paramsJson))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body
}

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
	fmt.Println("signature:", signParams["signature"])

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
