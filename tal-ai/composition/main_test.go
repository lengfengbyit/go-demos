package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const (
	AI_SUBJECt_CLF_URL = "http://openai.100tal.com/ai-subject-clf"        // 判断输入文本所属学科
	COM_EDUCATION_RUL  = "https://openai.100tal.com/aiimage/comeducation" // 教育通用 OCR
)

const (
	ROOT = "/Users/tal/GolandProjects/go-demos/tal-ai/composition/"
	// 语文作文孵蛋记
	FDJ_IMG_URL = "https://img.xiaohuasheng.cn/267734/ExperienceImage/choose0.6245587831735514.jpg"
	FA_IMG_URL  = "../imgs/zw.png"
)

type ParamMap map[string]any

func TestAiSubjectClf(t *testing.T) {
	enTxt := "there is no need to use a human to revise a passage"

	params := ParamMap{
		"params": ParamMap{
			"text": enTxt,
		},
	}

	urlParams := createSign(APP_SECRET, params, nil)
	resp := httpPost(AI_SUBJECt_CLF_URL+"?"+urlParams, params)
	fmt.Println(string(resp))
}

// 教育通用 OCR, 用于检测是否存在手写字体
// 返回文字、文字的位置，没有段落信息；
func TestComEducation(t *testing.T) {
	params := ParamMap{
		//"image_url":   FDJ_IMG_URL,
		"image_base64": getBase64Image(FA_IMG_URL),
		"function":     2,         // 全部类型检测+仅输出手写结果
		"subject":      "liberat", // 文科，不识别公式
		"textInTable":  true,      // 是否打印图片中表格里的文字
		"textInImage":  true,      // 是否打印图片中的文字
	}

	urlParams := createSign(APP_SECRET, params, nil)
	resp := httpPost(COM_EDUCATION_RUL+"?"+urlParams, params)
	fmt.Println(string(resp))
	saveFile("./com_education_result.json", resp)
}

// 中文作文手写识别接口，识别出文字位置、段落信息
func TestCnComposition(t *testing.T) {
	resp := CnComposition("../imgs/pic_fangge.jpg")
	t.Logf("%+v", resp)
}

// 大模型：语文作文批改，好词好句、评分、总评等
func TestChCompositionGpt(t *testing.T) {
	compResp := CnComposition("../imgs/zw.png") // 中文 OCR 获取文字内容
	title, content := GetCompositionContent(compResp)
	// 流式返回
	CnCompositionGPT(title + "\n" + content)
}

func saveFile(filePath string, content []byte) {
	_ = os.WriteFile(filepath.Join(ROOT, filePath), content, 0644)
}
