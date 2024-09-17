package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// 手写识别

const (
	CN_COMPOSITION_URL = "http://openai.100tal.com/aiimage/cn-composition"
	EN_OCR_URL         = "http://openai.100tal.com/aiocr/english-ocr"
	IMG_URL            = "../imgs/pic_fangge.jpg"
	EN_IMG_URL         = "../imgs/en.jpg"
)

var (
	APP_KEY    = os.Getenv("TAL_OPENAI_KEY")
	APP_SECRET = os.Getenv("TAL_OPENAI_SECRET")
)

func main() {
	EnOcr()
}

// EnOcr 英文 OCR
func EnOcr() {
	params := map[string]any{
		//"image_url": IMG_URL,
		"image_base64": getBase64Image(EN_IMG_URL),
		"rotate":       0, // 是否旋转
		"enhanced":     0, // 是否多次识别
	}
	urlParams := createSign(APP_SECRET, params, nil)
	body := httpPost(EN_OCR_URL+"?"+urlParams, params)
	fmt.Printf("%+v\n", string(body))
}

// CnComposition 中文手写识别
func CnComposition() {
	if APP_KEY == "" || APP_SECRET == "" {
		log.Fatal("APP_KEY or APP_SECRET is empty")
	}
	params := map[string]any{
		//"image_url":           IMG_URL,
		"image_base64":        getBase64Image(IMG_URL),
		"details":             false,
		"paragraph_detection": true,
	}
	urlParams := createSign(APP_SECRET, params, nil)

	body := httpPost(CN_COMPOSITION_URL+"?"+urlParams, params)
	_ = os.WriteFile("result.json", body, 0644)
}

func getBase64Image(imgPath string) string {
	content, err := os.ReadFile(imgPath)
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
