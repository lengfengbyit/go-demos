package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
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
	panicOnError(err)
	resp, err := http.Post(url, "application/json", bytes.NewReader(paramsJson))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body
}

// httpPostStream 流式读取响应体
func httpPostStream(url string, params map[string]any, dataCh chan<- []byte) {
	client := &http.Client{}

	paramsJson, err := json.Marshal(params)
	panicOnError(err)

	req, err := http.NewRequest("POST", url, bytes.NewReader(paramsJson))
	panicOnError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	panicOnError(err)
	defer resp.Body.Close()

	defer close(dataCh)

	// 逐块读取响应体
	bodyReader := resp.Body
	buf := make([]byte, 1024)
	for {
		n, err := bodyReader.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		dataCh <- buf[:n]
	}
}
