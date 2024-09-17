package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// 手写识别

const (
	CN_COMPOSITION_URL        = "http://openai.100tal.com/aiimage/cn-composition"   // 中文手写识别
	CN_COMPOSITION_REVISE_URL = "https://openai.100tal.com/aimathgpt/ch-compostion" // 中文作文批改
	EN_OCR_URL                = "http://openai.100tal.com/aiocr/english-ocr"        // 英文 OCR
	IMG_URL                   = "../imgs/pic_fangge.jpg"
	EN_IMG_URL                = "../imgs/en.jpg"
)

var (
	APP_KEY    = os.Getenv("TAL_OPENAI_KEY")
	APP_SECRET = os.Getenv("TAL_OPENAI_SECRET")
)

func main() {
	//EnOcr()

	compResp := CnComposition()
	title, content := GetCompositionContent(compResp)
	CnCompositionRevise(title + "\n" + content)

}

func CnCompositionRevise(content string) {
	prompt := "你的角色是一名语文老师，我是一位三年级小学生，请对用户输入的作文的每个段落进行分别点评，你输出的要求如下：" +
		"1、知识范围：小学。" +
		"2、对话风格：鼓励型。" +
		"3、每个段落的点评请以'作文第一段点评'、'作文第二段点评'等格式开头。完成每段的点评后，请直接结束回答，不需要添加任何总结性的话语或结尾语。" +
		"4、最后一段点评完后，直接结束输出，不要总结。" +
		"5、不要举例子，不要引用原文，并且忽略作文中出现的错别字、打字错误或者用词错误的问题。" +
		"6、从文章整体去分析，不需要说明文章字、词、句子的具体细节错误，不要分析作文中表达不准确的错误。" +
		"7、一句话说明优点，一句话说明缺点，一句话给出建议，三句话连成一段。"
	message := &CnCompositionReviseMessage{
		Role:    "user",
		Content: prompt + content,
	}

	messages := []*CnCompositionReviseMessage{message}

	params := map[string]any{
		"is_stream": true,
		"messages":  messages,
	}
	urlParams := createSign(APP_SECRET, params, nil)

	done := make(chan struct{})
	dataCh := make(chan []byte, 10)
	go func() {
		for body := range dataCh {
			var data CnCompositionReviseResp
			_ = json.Unmarshal(body, &data)
			if data.Code != 20000 {
				continue
			}
			fmt.Print(data.Data.Result)
		}
		done <- struct{}{}
	}()

	httpPostStream(CN_COMPOSITION_REVISE_URL+"?"+urlParams, params, dataCh)
	<-done
}

// GetCompositionContent 获取作文内容
func GetCompositionContent(resp *CnCompositionResp) (string, string) {
	var content strings.Builder
	for _, items := range resp.Data.EssayInfo.ParaOcrResult {
		for _, item := range items {
			content.WriteString(item.LineOcrResult)
		}
	}

	title := resp.Data.TitleInfo.TitleOcrResult

	return title, content.String()
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
func CnComposition() *CnCompositionResp {
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

	var data CnCompositionResp
	_ = json.Unmarshal(body, &data)
	return &data
}

type CnCompositionResp struct {
	Code int `json:"code"`
	Data struct {
		EssayInfo struct {
			ParaOcrResult [][]struct {
				LineCharInfo []struct {
					CharLocation []struct {
						X int `json:"x"`
						Y int `json:"y"`
					} `json:"char_location"`
					LineCharTopn []struct {
						CharConfidence float64 `json:"char_confidence"`
						CharOcrResult  string  `json:"char_ocr_result"`
					} `json:"line_char_topn"`
				} `json:"line_char_info"`
				LineOcrResult string `json:"line_ocr_result"`
				ParaType      int    `json:"para_type"`
			} `json:"para_ocr_result"`
		} `json:"essay_info"`
		TitleInfo struct {
			TitleCharInfo []struct {
				CharLocation []struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"char_location"`
				CharOcrTopn []struct {
					CharConfidence float64 `json:"char_confidence"`
					CharOcrResult  string  `json:"char_ocr_result"`
				} `json:"char_ocr_topn"`
			} `json:"title_char_info"`
			TitleOcrResult string `json:"title_ocr_result"`
		} `json:"title_info"`
	} `json:"data"`
	Msg       string `json:"msg"`
	RequestId string `json:"requestId"`
}

type CnCompositionReviseMessage struct {
	Role    string `json:"role"` // user,assistant
	Content string `json:"content"`
}

type CnCompositionReviseResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		IsEnd  int    `json:"is_end"`
		Result string `json:"result"`
		Mod    string `json:"mod"`
	} `json:"data"`
	RequestId string `json:"request_id"`
}
