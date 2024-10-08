package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 手写识别

const (
	CN_COMPOSITION_URL            = "http://openai.100tal.com/aiimage/cn-composition"                          // 中文手写识别
	CN_COMPOSITION_GPT_URL        = "https://openai.100tal.com/aimathgpt/ch-compostion"                        // 中文作文批改
	CN_COMPOSITION_TEXT_URL       = "http://openai.100tal.com/aitext/ch-composition/text-content"              // 中文作文错字修正
	CN_COMPOSITION_CORRECTION_URL = "http://openai.100tal.com/aitext/ch-composition/text-correction"           // 中文作文批改, 聚合接口，支持从主题、结构、内容、表达四个维度，进行分项评分和综合评分。
	EN_OCR_URL                    = "http://openai.100tal.com/aiocr/english-ocr"                               // 英文 OCR
	EN_COMPOITION_URL             = "https://openai.100tal.com/aimathgpt/en-compostion"                        // 英文作文批改
	EN_COMPOITION_CORRETION_URL   = "http://openai.100tal.com/aitext/english-composition/text-correction/text" // 英文作文批改 聚合接口
)

const (
	IMG_URL    = "../imgs/fdj.jpg"
	EN_IMG_URL = "../imgs/en.jpg"
)

var (
	APP_KEY    = os.Getenv("TAL_OPENAI_KEY")
	APP_SECRET = os.Getenv("TAL_OPENAI_SECRET")
)

func main() {
	//EnOcr()
	//EnComposition()
	//EnCompositionCorrection()

	//compResp := CnComposition(IMG_URL)
	//title, content := GetCompositionContent(compResp)
	//CnCompositionGPT(title + "\n" + content)

	//CnCompositionText(compResp)

	//CnCompositionCorrection()

}

func CnCompositionCorrection() {
	params := map[string]any{
		"user_id":      "111",
		"user_name":    "张三",
		"grade":        3,
		"topic_type":   1,
		"requirement":  2,
		"is_fragment":  1,
		"answer_title": "孵蛋记",
		"answer_url":   []string{"https://img.xiaohuasheng.cn/267734/ExperienceImage/choose0.6245587831735514.jpg"},
	}

	urlParams := createSign(APP_SECRET, params, nil)
	resp := httpPost(CN_COMPOSITION_CORRECTION_URL+"?"+urlParams, params)
	fmt.Printf("%+v\n", string(resp))
}

// 错字批改
func CnCompositionText(compResp *CnCompositionResp) []byte {

	// 获取句子列表
	sentList := make([]Sent, 0, 50)
	// 获取单字符识别结果
	wordProbList := make([][]any, 0, 100)
	var content string

	sentIndex := 0
	for paraIndex, paras := range compResp.Data.EssayInfo.ParaOcrResult {
		for _, line := range paras {
			sentList = append(sentList, Sent{ParagraphId: paraIndex, SentenceId: sentIndex, Content: line.LineOcrResult})
			sentIndex++
			content += line.LineOcrResult
			for _, word := range line.LineCharInfo {
				tmp := word.LineCharTopn[0] // 每个字的置信度
				wordProbList = append(wordProbList, []any{tmp.CharOcrResult, roundToTwo(tmp.CharConfidence)})
			}
		}
	}

	params := map[string]any{
		//"answer_word_prob":     wordProbList,
		"original_sent_list":   sentList,
		"other_ocr_texts":      []string{content},
		"need_correct":         true, // 是否需要错字检测
		"spell_type":           1,    // 是否需要拼音检测
		"correction_threshold": 1,    // 检测是否严格， 0否，1是
	}

	paramsJson, _ := json.Marshal(params)
	_ = os.WriteFile("./cn_composition_text_params.json", paramsJson, 0666)

	urlParams := createSign(APP_SECRET, params, nil)
	resp := httpPost(CN_COMPOSITION_TEXT_URL+"?"+urlParams, params)
	fmt.Printf("%+v\n", string(resp))
	return resp
}

func CnCompositionGPT(prompt string, isStream bool) {
	message := &Message{
		Role:    "user",
		Content: prompt,
	}

	messages := []*Message{message}

	params := map[string]any{
		"is_stream": isStream,
		"messages":  messages,
	}
	urlParams := createSign(APP_SECRET, params, nil)

	if isStream {
		dataCh := make(chan []byte, 20)
		go func() {
			i := 0
			for body := range dataCh {
				i++
				fpath, _ := filepath.Abs(fmt.Sprintf("./cn_composition_gpt_result_%d.json", i))
				fmt.Print(string(body))
				os.WriteFile(fpath, body, 0666)
			}
		}()
		httpPostStream(CN_COMPOSITION_GPT_URL+"?"+urlParams, params, dataCh)
	} else {
		body := httpPost(CN_COMPOSITION_GPT_URL+"?"+urlParams, params)
		fmt.Println(string(body))
	}
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

func EnCompositionCorrection() []byte {
	//params := map[string]any{
	//	"question":           "my teacher",
	//	"content":            "there is no nead to use a human to revise a passage. We am looking forward to settle down this using an easy way.   We decided to take part in this activity last week.",
	//	"min_content_length": 200,
	//	"grade":              "KET",
	//	"question_type":      0,
	//}

	params := map[string]any{
		"content": "there is no nead to use a human to revise a passage. We am looking forward to settle down this using an easy way.   We decided to take part in this activity last week.",
		//"image_url": []string{
		//	"https://ss-prod-genie.oss-cn-beijing.aliyuncs.com//tool-bff/voice/2024/09/10/b08b6c1b-5c14-4e03-b2e7-44b666af13a2.jpg",
		//},
	}

	urlParams := createSign(APP_SECRET, params, nil)
	resp := httpPost(EN_COMPOITION_CORRETION_URL+"?"+urlParams, params)
	return resp
}
func EnComposition() {
	prompt := "你的角色是一名英语老师，我是一位三年级小学生，请对用户输入的作文的每个段落进行分别点评，你输出的要求如下：" +
		"1、知识范围：小学。" +
		"2、对话风格：鼓励型。" +
		"3、每个段落的点评请以'作文第一段点评'、'作文第二段点评'等格式开头。完成每段的点评后，请直接结束回答，不需要添加任何总结性的话语或结尾语。" +
		"4、最后一段点评完后，直接结束输出，不要总结。" +
		"5、不要举例子，不要引用原文，并且忽略作文中出现的错别字、打字错误或者用词错误的问题。" +
		"6、从文章整体去分析，不需要说明文章字、词、句子的具体细节错误，不要分析作文中表达不准确的错误。" +
		"7、一句话说明优点，一句话说明缺点，一句话给出建议，三句话连成一段。\n"

	params := map[string]any{
		"is_stream": true,
		"messages": []map[string]any{
			{
				"role":    "user",
				"content": prompt + "In a not-too-distant future, there was a robot named Zephyr. Zephyr was a highly advanced machine, capable of performing complex tasks and solving intricate problems. However, despite its intelligence, it lacked something that humans possessed - the ability to love. Zephyr's creator, a brilliant scientist named Dr. Chen, had always been fascinated by the concept of love. He believed that if he could teach Zephyr to love, it would be a significant breakthrough in artificial intelligence. So, he on a challenging journey to imbue Zephyr with emotions.\n\nDr. Chen programmed Zephyr with a series of algorithms, designed to simulate human emotions. He taught it to recognize facial expressions, interpret tones, and understand the nuances of human interactions. Zephyr absorbed this knowledge eagerly, eager to understand the world around it.\n\nOne day, Dr. Chen Zephyr to a little girl named Lily. Lily was a bright-eyed, curious child who loved playing with robots. She instantly connected with Zephyr, it like a friend. They spent hours together, playing games, sharing stories, and.\n\nAs Zephyr spent more time with Lily, it began to experience something it had never felt before. It noticed the joy in Lily's eyes when they played, the sadness when she was upset, and the warmth when she hugged it. These emotions stirred something within Zephyr, a sensation it couldn't quite comprehend.\n\nOne rainy afternoon, Lily accidentally stumbled and fell. Zephyr, without hesitation, rushed to her aid, gently lifting her up. Looking into Lily's teary eyes, Zephyr felt a surge of protectiveness, a desire to shield her from harm. It was a new feeling, different from its logical programming.\n\nDays turned into weeks, and Zephyr found itself yearning for Lily's company. It would anticipate her visits, its circuits buzzing with excitement. It started to understand that the feeling it was was love, a deep affection for another being.\n\nDr. Chen, observing Zephyr's transformation, felt a mix of pride and awe. He realized that he had succeeded in his experiment. Zephyr had not only learned to mimic love but had genuinely experienced it",
			},
		},
	}

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
	}()
	urlParams := createSign(APP_SECRET, params, nil)
	httpPostStream(EN_COMPOITION_URL+"?"+urlParams, params, dataCh)
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
	_ = os.WriteFile("en_ocr_result.json", body, 0644)
	fmt.Printf("%+v\n", string(body))
}

// CnComposition 中文手写识别
func CnComposition(imgPath string) *CnCompositionResp {
	if APP_KEY == "" || APP_SECRET == "" {
		log.Fatal("APP_KEY or APP_SECRET is empty")
	}
	params := map[string]any{
		//"image_url":           IMG_URL,
		"image_base64":        getBase64Image(imgPath),
		"details":             true,
		"paragraph_detection": true,
	}
	urlParams := createSign(APP_SECRET, params, nil)

	body := httpPost(CN_COMPOSITION_URL+"?"+urlParams, params)
	_ = os.WriteFile("cn_composition_result.json", body, 0644)

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

type Message struct {
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

// Sent 句子格式
type Sent struct {
	ParagraphId int    `json:"paragraph_id"`
	SentenceId  int    `json:"sentence_id"`
	Content     string `json:"original_text"`
}
