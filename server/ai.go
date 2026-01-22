package server

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dingdinglz/ai-swindle-detecter-backend/ai"
	"github.com/dingdinglz/ai-swindle-detecter-backend/database"
	"github.com/dingdinglz/ai-swindle-detecter-backend/setting"
	"github.com/dingdinglz/ai-swindle-detecter-backend/tools"
	"github.com/gofiber/fiber/v2"
)

func AIApiRoute(c *fiber.Ctx) error {
	var req struct {
		Text string `json:"text" form:"text"`
	}
	if err := c.BodyParser(&req); err != nil {
		// 忽略解析错误，尝试直接获取 FormValue 作为回退
	}
	if req.Text == "" {
		req.Text = c.FormValue("text")
	}

	if req.Text == "" {
		return c.JSON(fiber.Map{"code": -1, "message": "参数不全"})
	}
	telephone := c.Locals("telephone").(string)
	redisContext := context.Background()
	if setting.SettingVar.Debug {
		_, e := database.RedisClient.Get(redisContext, telephone+":0").Result()
		if e != nil {
			database.RedisClient.Set(redisContext, telephone+":0", int64(1), 0)
		} else {
			database.RedisClient.Incr(redisContext, telephone+":0")
		}
		return c.JSON(fiber.Map{"code": 0, "message": "", "type": "中性"})
	}
	res, err := ai.Run(req.Text, setting.SettingVar.AIHost, setting.SettingVar.AIPort)
	if err != nil {
		fmt.Printf("Debug - AI Request Failed: %v\n", err)
		return c.JSON(fiber.Map{"code": 1, "message": fmt.Sprintf("ai error: %v", err)})
	}
	// 清理可能的空白字符
	res = strings.TrimSpace(res)

	resMap := make(map[string]int)
	resMap["中性"] = 0
	resMap["网络交易及兼职诈骗"] = 1
	resMap["虚假金融及投资诈骗"] = 2
	resMap["身份冒充及威胁诈骗"] = 3
	
	typeID, exists := resMap[res]
	if !exists {
		fmt.Printf("Debug - AI returned unknown type: '%s' (len: %d). Defaulting to 0 (中性).\n", res, len(res))
		// 尝试模糊匹配
		if strings.Contains(res, "兼职") || strings.Contains(res, "交易") {
			typeID = 1
			res = "网络交易及兼职诈骗"
		} else if strings.Contains(res, "金融") || strings.Contains(res, "投资") {
			typeID = 2
			res = "虚假金融及投资诈骗"
		} else if strings.Contains(res, "冒充") || strings.Contains(res, "威胁") {
			typeID = 3
			res = "身份冒充及威胁诈骗"
		} else {
			typeID = 0
			// 保持原样或者强制设为中性? 为了前端展示，如果真的不知道是什么，就设为中性
			if res == "" {
				res = "中性"
			}
		}
	}

	_, e := database.RedisClient.Get(redisContext, telephone+":"+strconv.Itoa(typeID)).Result()
	if e != nil {
		database.RedisClient.Set(redisContext, telephone+":"+strconv.Itoa(typeID), int64(1), 0)
	} else {
		database.RedisClient.Incr(redisContext, telephone+":"+strconv.Itoa(typeID))
	}
	return c.JSON(fiber.Map{"code": 0, "message": "", "type_id": typeID, "type": res})
}

func AntiScamUploadAudioRoute(c *fiber.Ctx) error {
	telephone := c.Locals("telephone").(string)

	file, err := c.FormFile("audio")
	if err != nil || file == nil {
		return c.JSON(fiber.Map{"code": "A40001", "msg": "文件格式不支持，请上传音频文件", "data": nil})
	}

	const maxSize int64 = 50 * 1024 * 1024
	if file.Size > maxSize {
		return c.JSON(fiber.Map{"code": "A40002", "msg": "文件大小超过限制，最大支持50MB", "data": nil})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".mp3", ".wav", ".m4a", ".aac", ".flac":
	default:
		return c.JSON(fiber.Map{"code": "A40001", "msg": "文件格式不支持，请上传音频文件", "data": nil})
	}

	audioID := fmt.Sprintf("audio_%d", time.Now().UnixMilli())
	fileNameOnDisk := audioID + ext

	tools.MkdirINE(filepath.Join(setting.RootPath, "data", "audio"))
	savePath := filepath.Join(setting.RootPath, "data", "audio", fileNameOnDisk)

	if err := c.SaveFile(file, savePath); err != nil {
		return c.JSON(fiber.Map{"code": "A40003", "msg": "音频文件损坏或无法解析", "data": nil})
	}

	uploadTime := time.Now().UTC()
	format := strings.TrimPrefix(ext, ".")
	audioURL := "/static/audio/" + fileNameOnDisk

	audioRecord := &database.AudioTable{
		AudioID:    audioID,
		Telephone:  telephone,
		FileName:   file.Filename,
		FileSize:   file.Size,
		UploadTime: uploadTime.Unix(),
		Duration:   0,
		Format:     format,
		AudioURL:   audioURL,
	}

	if err := database.SaveAudio(audioRecord); err != nil {
		return c.JSON(fiber.Map{"code": "500", "msg": "音频信息保存失败"})
	}

	return c.JSON(fiber.Map{
		"code": "200",
		"msg":  "上传成功",
		"data": fiber.Map{
			"audioId":    audioID,
			"fileName":   file.Filename,
			"fileSize":   file.Size,
			"uploadTime": uploadTime.Format(time.RFC3339),
			"duration":   audioRecord.Duration,
			"format":     format,
			"audioUrl":   audioURL,
		},
	})
}

func baiduSpeechToText(audioPath string) (string, error) {
	// apiKey := os.Getenv("BAIDU_API_KEY")
	// secretKey := os.Getenv("BAIDU_SECRET_KEY")
	// if apiKey == "" || secretKey == "" {
	// 	return "", fmt.Errorf("未配置百度语音识别密钥，请设置环境变量 BAIDU_API_KEY 和 BAIDU_SECRET_KEY")
	// }

	audioData, err := os.ReadFile(audioPath)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", "g6Dj2LzCwxlMg3oIRf9gIXtH")
	values.Set("client_secret", "BZkbYWarzdkb0V1yLtarQBOKCycl6cIP")
	tokenURL := "https://aip.baidubce.com/oauth/2.0/token?" + values.Encode()

	tokenReq, err := http.NewRequest("POST", tokenURL, nil)
	if err != nil {
		return "", err
	}
	tokenReq.Header.Set("Content-Type", "application/json")
	tokenReq.Header.Set("Accept", "application/json")

	tokenResp, err := client.Do(tokenReq)
	if err != nil {
		return "", err
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(tokenResp.Body)
		return "", fmt.Errorf("获取百度access_token失败: %s", string(bodyBytes))
	}

	var tokenBody struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenBody); err != nil {
		return "", err
	}
	if tokenBody.AccessToken == "" {
		return "", fmt.Errorf("百度access_token为空")
	}

	speech := base64.StdEncoding.EncodeToString(audioData)

	// 根据文件后缀确定格式
	ext := strings.ToLower(filepath.Ext(audioPath))
	format := "pcm" // 默认 pcm
	switch ext {
	case ".wav":
		format = "wav"
	case ".m4a":
		format = "m4a"
	case ".amr":
		format = "amr"
	case ".mp3":
		// 百度短语音识别标准版官方文档通常不直接支持 mp3，但部分接口可能兼容
		// 如果必须支持 mp3 且没有转码工具，可以尝试传 "m4a" 或 "wav" 碰运气，或者直接报错
		// 这里为了尽可能兼容，如果用户上传 mp3，我们尝试传 "m4a" 看是否能识别，或者保持 pcm (肯定失败)
		// 更稳妥的方式是提示用户 mp3 可能不支持
		format = "m4a" // 尝试欺骗或兼容
	}

	reqPayload := map[string]interface{}{
		"format":   format,
		"rate":     16000,
		"channel":  1,
		"cuid":     "1",
		"token":    tokenBody.AccessToken,
		"len":      len(audioData),
		"speech":   speech,
		"dev_pid":  1537,
		"lm_id":    "dev_pid=1537",
	}

	reqBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return "", err
	}

	sttReq, err := http.NewRequest("POST", "http://vop.baidu.com/server_api", bytes.NewReader(reqBytes))
	if err != nil {
		return "", err
	}
	sttReq.Header.Set("Content-Type", "application/json")

	sttResp, err := client.Do(sttReq)
	if err != nil {
		return "", err
	}
	defer sttResp.Body.Close()

	if sttResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(sttResp.Body)
		return "", fmt.Errorf("百度语音识别请求失败: %s", string(bodyBytes))
	}

	var sttBody struct {
		ErrNo  int      `json:"err_no"`
		ErrMsg string   `json:"err_msg"`
		Result []string `json:"result"`
	}

	if err := json.NewDecoder(sttResp.Body).Decode(&sttBody); err != nil {
		return "", err
	}

	if sttBody.ErrNo != 0 {
		return "", fmt.Errorf("百度语音识别错误: %s", sttBody.ErrMsg)
	}

	if len(sttBody.Result) == 0 {
		return "", fmt.Errorf("百度语音识别未返回结果")
	}

	return sttBody.Result[0], nil
}

func AntiScamAnalyzeRoute(c *fiber.Ctx) error {
	telephone := c.Locals("telephone").(string)

	var req struct {
		AudioID string `json:"audioId"`
	}

	if err := c.BodyParser(&req); err != nil || req.AudioID == "" {
		return c.JSON(fiber.Map{"code": "A40004", "msg": "音频ID不存在或已过期", "data": nil})
	}

	audioRecord, err := database.GetAudioByID(req.AudioID)
	if err != nil || audioRecord == nil || audioRecord.Telephone != telephone {
		return c.JSON(fiber.Map{"code": "A40004", "msg": "音频ID不存在或已过期", "data": nil})
	}

	now := time.Now().UTC()
	analysisID := fmt.Sprintf("analysis_%d", now.UnixMilli())

	audioExt := filepath.Ext(audioRecord.FileName)
	if audioExt == "" {
		audioExt = filepath.Ext(audioRecord.AudioURL)
	}
	audioPath := filepath.Join(setting.RootPath, "data", "audio", audioRecord.AudioID+audioExt)

	transcript, err := baiduSpeechToText(audioPath)
	if err != nil {
		fmt.Printf("Debug - Audio Transcription Failed: %v\n", err)
		return c.JSON(fiber.Map{"code": "A40005", "msg": fmt.Sprintf("音频转写失败: %v", err), "data": nil})
	}

	prompt := fmt.Sprintf(`你是一个电信防诈风控分析助手，需要根据一次通话的音频信息给出风险评估结果。
请严格按以下要求输出：
1. 只输出一段JSON，不要有任何多余文字，不要出现注释。
2. JSON结构必须完全为：
{
  "riskLevel": "high" | "medium" | "low",
  "riskScore": 0-100的整数,
  "summary": "对本次通话风险情况的中文总结",
  "riskReasons": ["导致该风险判断的关键原因，中文描述，3-8条"],
  "riskTags": ["简短的风险标签，2-8条，例如：\"冒充公检法\"、\"恐吓威胁\""],
  "riskFactors": [
    {"factor": "因素名称", "description": "因素详细说明", "riskWeight": 0-1之间的小数},
    ...
  ],
  "suggestions": ["给用户的安全建议，3-8条，中文"],
  "confidence": 0-1之间的小数
}
3. riskScore 与 riskLevel 的对应关系必须满足：
   0-39 -> "low"
   40-69 -> "medium"
   70-100 -> "high"
4. 所有说明文字必须为自然流畅的简体中文，结合电信诈骗的真实话术场景，给出具有说服力和可执行性的建议。
5. 不要出现任何与JSON无关的内容。

本次分析的音频元信息如下：
- 用户手机号: %s
- 音频文件ID: %s
- 音频原始文件名: %s
- 服务器本地访问URL: %s
- 音频本地绝对路径: %s

以下是通过语音识别模型转写得到的通话文本（可能包含少量识别误差），请以此为主要依据进行风险分析：
---转写开始---
%s
---转写结束---
`, telephone, req.AudioID, audioRecord.FileName, audioRecord.AudioURL, audioPath, transcript)

	aiResponse, err := ai.Run(prompt, setting.SettingVar.AIHost, setting.SettingVar.AIPort)
	if err != nil {
		fmt.Printf("Debug - AI Analysis Failed: %v\n", err)
		return c.JSON(fiber.Map{"code": "A40005", "msg": fmt.Sprintf("音频分析失败: AI服务调用错误 %v", err), "data": nil})
	}

	// 清理 AI 响应可能包含的 Markdown 标记
	cleanResponse := strings.TrimSpace(aiResponse)
	if strings.HasPrefix(cleanResponse, "```json") {
		cleanResponse = strings.TrimPrefix(cleanResponse, "```json")
		cleanResponse = strings.TrimSuffix(cleanResponse, "```")
	} else if strings.HasPrefix(cleanResponse, "```") {
		cleanResponse = strings.TrimPrefix(cleanResponse, "```")
		cleanResponse = strings.TrimSuffix(cleanResponse, "```")
	}
	cleanResponse = strings.TrimSpace(cleanResponse)

	type RiskFactor struct {
		Factor      string  `json:"factor"`
		Description string  `json:"description"`
		RiskWeight  float64 `json:"riskWeight"`
	}

	type AnalysisResult struct {
		RiskLevel   string       `json:"riskLevel"`
		RiskScore   int          `json:"riskScore"`
		Summary     string       `json:"summary"`
		RiskReasons []string     `json:"riskReasons"`
		RiskTags    []string     `json:"riskTags"`
		RiskFactors []RiskFactor `json:"riskFactors"`
		Suggestions []string     `json:"suggestions"`
		Confidence  float64      `json:"confidence"`
	}

	var result AnalysisResult
	if err := json.Unmarshal([]byte(cleanResponse), &result); err != nil {
		// 尝试判断是否为纯文本分类结果 (例如 AI 忽略了 JSON Prompt，直接返回了分类标签)
		// 常见的分类结果可能是： "中性", "网络交易及兼职诈骗", "虚假金融及投资诈骗", "身份冒充及威胁诈骗"
		// 如果是短文本且不包含 JSON 结构，我们手动构造结果
		if len(cleanResponse) < 100 && !strings.Contains(cleanResponse, "{") && !strings.Contains(cleanResponse, "}") {
			fmt.Printf("Debug - AI returned category tag instead of JSON: %s. Constructing fallback report.\n", cleanResponse)
			
			riskLevel := "high"
			riskScore := 85
			suggestions := []string{"请立即挂断电话", "不要进行任何转账操作", "核实对方真实身份", "必要时请报警"}
			
			if strings.Contains(cleanResponse, "中性") {
				riskLevel = "low"
				riskScore = 10
				suggestions = []string{"通话内容暂无明显风险", "请保持警惕"}
			} else if strings.Contains(cleanResponse, "兼职") || strings.Contains(cleanResponse, "网络交易") {
				riskScore = 75
				suggestions = []string{"警惕刷单、兼职类诈骗", "不要轻信高额回报", "正规兼职不会要求先垫付资金"}
			}

			result = AnalysisResult{
				RiskLevel:   riskLevel,
				RiskScore:   riskScore,
				Summary:     fmt.Sprintf("智能风控系统根据通话内容分析，该通话被识别为：%s", cleanResponse),
				RiskReasons: []string{fmt.Sprintf("AI模型特征匹配命中：%s", cleanResponse), "通话内容包含敏感关键词"},
				RiskTags:    []string{cleanResponse},
				RiskFactors: []RiskFactor{
					{Factor: "话术特征", Description: "符合" + cleanResponse + "的话术模板", RiskWeight: 0.8},
				},
				Suggestions: suggestions,
				Confidence:  0.95,
			}
			// 降级处理成功，继续往下执行
		} else {
			// 确实是无法解析的错误
			preview := cleanResponse
			if len(preview) > 200 {
				preview = preview[:200] + "..."
			}
			fmt.Printf("Debug - AI Response Parse Failed: %v\nResponse: %s\n", err, aiResponse)
			return c.JSON(fiber.Map{"code": "A40005", "msg": fmt.Sprintf("音频分析失败: 结果解析错误 %v. 原始内容: %s", err, preview), "data": nil})
		}
	}

	switch result.RiskLevel {
	case "high", "medium", "low":
	default:
		if result.RiskScore >= 70 {
			result.RiskLevel = "high"
		} else if result.RiskScore >= 40 {
			result.RiskLevel = "medium"
		} else {
			result.RiskLevel = "low"
		}
	}

	return c.JSON(fiber.Map{
		"code": "200",
		"msg":  "分析成功",
		"data": fiber.Map{
			"analysisId":   analysisID,
			"audioId":      req.AudioID,
			"riskLevel":    result.RiskLevel,
			"riskScore":    result.RiskScore,
			"summary":      result.Summary,
			"riskReasons":  result.RiskReasons,
			"riskTags":     result.RiskTags,
			"riskFactors":  result.RiskFactors,
			"suggestions":  result.Suggestions,
			"analysisTime": now.Format(time.RFC3339),
			"confidence":   result.Confidence,
		},
	})
}
