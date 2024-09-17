package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestSign(t *testing.T) {

	var (
		access_key_id     = "4300865906099200"
		access_key_secret = "7a810a4245534cdab787bd82d0a63ca9"
	)

	body := map[string]any{
		"key3": "value1",
		"key4": "value2",
	}

	params := map[string]string{
		"key1":            "value1",
		"key2":            "value2",
		"signature_nonce": "fd16bc90-08a5-4034-a06b-aa7004f9d0c5",
		"timestamp":       "2020-04-14T11:11:30",
		"access_key_id":   access_key_id,
	}

	_ = createSign(access_key_secret, body, params)
}

func TestHmac(t *testing.T) {

	accessKeySecret := "7a810a4245534cdab787bd82d0a63ca9"
	//stringToSign := `access_key_id=4300865906099200&key1=value1&key2=value2&request_body={"key3":"value3","key4":"value4"}&signature_nonce=fd16bc90-08a5-4034-a06b-aa7004f9d0c5&timestamp=2020-04-14T11:11:30`

	// 计算证书签名
	//stringToSign := "access_key_id=4300865906099200&key1=value1&key2=value2&request_body={\"key3\":\"value1\",\"key4\":\"value2\"}&signature_nonce=fd16bc90-08a5-4034-a06b-aa7004f9d0c5&timestamp=2020-04-14T11:11:30"
	stringToSign := `access_key_id=4300865906099200&key1=value1&key2=value2&request_body={"key3":"value1","key4":"value2"}&signature_nonce=fd16bc90-08a5-4034-a06b-aa7004f9d0c5&timestamp=2020-04-14T11:11:30`
	//stringToSign := `access_key_id=4300865906099200&key1=value1&key2=value2&request_body={"key3":"value3","key4":"value4"}&signature_nonce=fd16bc90-08a5-4034-a06b-aa7004f9d0c5&timestamp=2020-04-14T11:11:30`

	// 进行base64 encode
	secret := accessKeySecret + "&"
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	fmt.Println(signature)

}
