package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

const (
	AI_SUBJECt_CLF_URL     = "http://openai.100tal.com/ai-subject-clf"                             // 判断输入文本所属学科
	COM_EDUCATION_RUL      = "https://openai.100tal.com/aiimage/comeducation"                      // 教育通用 OCR
	EN_TEXT_CORRECTION_URL = "http://openai.100tal.com/aitext/english-composition/text-correction" // 英语作文批改
)

const (
	ROOT = "/Users/tal/GolandProjects/go-demos/tal-ai/composition/"
	// 语文作文孵蛋记
	FDJ_IMG_URL   = "https://img.xiaohuasheng.cn/267734/ExperienceImage/choose0.6245587831735514.jpg"
	FA_IMG_URL    = "../imgs/zw.png"
	EN_ZW_IMG_URL = "../imgs/en_zw.jpg"
)

type ParamMap map[string]any

func TestEnTextCorrection(m *testing.T) {
	//params := ParamMap{
	//	"content": "there is no nead to use a human to revise a passage. We am looking forward to settle down this using an easy way.   We decided to take part in this activity last week.",
	//}
	//img := "https://ss-prod-genie.oss-cn-beijing.aliyuncs.com//tool-bff/voice/2024/09/03/0015c49d-ec49-4cb6-b45a-69823a8a4fcc"
	//img := "https://ss-prod-genie.oss-cn-beijing.aliyuncs.com//tool-bff/voice/2024/09/03/16572e15-7a37-48be-906e-a01fcc64b348"
	img := "https://ss-prod-genie.oss-cn-beijing.aliyuncs.com//tool-bff/voice/2024/09/03/0165f229-7e26-4f4b-adfd-a8016c6ebe95"
	params := ParamMap{
		"image_url": []string{img},
	}
	urlParams := createSign(APP_SECRET, params, nil)
	resp := httpPost(EN_TEXT_CORRECTION_URL+"?"+urlParams, params)
	fmt.Println(string(resp))
	saveFile("./en_text_correction.json", resp)
}

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
		"image_url": "https://ss-prod-genie.oss-cn-beijing.aliyuncs.com//tool-bff/voice/2024/09/10/9d4d5181-aa85-4c4f-8c26-a80a875ad222.jpg",
		//"image_base64": getBase64Image(FA_IMG_URL),
		"function":    2,         // 全部类型检测+仅输出手写结果
		"subject":     "liberat", // 文科，不识别公式
		"textInTable": false,     // 是否打印图片中表格里的文字
		"textInImage": false,     // 是否打印图片中的文字
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
	//scorePrompt = fmt.Sprintf(scorePrompt, "三", title+"\n"+content)
	//CnCompositionGPT(scorePrompt)

	// 字词句批改
	prompt := fmt.Sprintf(correctionPrompt, title+"\n"+content)
	CnCompositionGPT(prompt, false)
}

// 中文作文错字批改
func TestCnCompositionText(t *testing.T) {
	resp := CnComposition("../imgs/pic_fangge.jpg")
	body := CnCompositionText(resp)
	saveFile("cn_composition_text_result.json", body)
}

func TestEnCompositionCorrection(t *testing.T) {
	resp := EnCompositionCorrection()
	saveFile("en_composition_correction_result.json", resp)
}

func TestEnOcr(t *testing.T) {
	EnOcr()
}

func stringToMd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
func TestMd5(t *testing.T) {
	fmt.Println(stringToMd5("hello"))

	fmt.Println(utf8.RuneCountInString("hello 你好"))
}

func saveFile(filePath string, content []byte) {
	_ = os.WriteFile(filepath.Join(ROOT, filePath), content, 0644)
}

func TestSlice(t *testing.T) {
	sl := []int{1, 2, 3, 4, 5}
	for i, num := range sl {
		fmt.Print(num, " ")
		if num == 3 {
			sl = append(sl[:i], sl[i+1:]...)
		}
	}

	fmt.Println(sl)
}

func TestDemo(t *testing.T) {
	str := "hi \n\n"
	str2 := strings.Trim(str, "\n")
	fmt.Println(len(str2))
}
