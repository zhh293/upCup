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
	res := ai.Run(req.Text, setting.SettingVar.AIPort)
	if res == "err" {
		return c.JSON(fiber.Map{"code": 1, "message": "ai error"})
	}
	resMap := make(map[string]int)
	resMap["中性"] = 0
	resMap["网络交易及兼职诈骗"] = 1
	resMap["虚假金融及投资诈骗"] = 2
	resMap["身份冒充及威胁诈骗"] = 3
	_, e := database.RedisClient.Get(redisContext, telephone+":"+strconv.Itoa(resMap[res])).Result()
	if e != nil {
		database.RedisClient.Set(redisContext, telephone+":"+strconv.Itoa(resMap[res]), int64(1), 0)
	} else {
		database.RedisClient.Incr(redisContext, telephone+":"+strconv.Itoa(resMap[res]))
	}
	return c.JSON(fiber.Map{"code": 0, "message": "", "type_id": resMap[res], "type": res})
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
	apiKey := os.Getenv("BAIDU_API_KEY")
	secretKey := os.Getenv("BAIDU_SECRET_KEY")
	if apiKey == "" || secretKey == "" {
		return "", fmt.Errorf("未配置百度语音识别密钥，请设置环境变量 BAIDU_API_KEY 和 BAIDU_SECRET_KEY")
	}

	audioData, err := os.ReadFile(audioPath)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", apiKey)
	values.Set("client_secret", secretKey)
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

	reqPayload := map[string]interface{}{
		"format":   "pcm",
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

	sttReq, err := http.NewRequest("POST", "https://vop.baidu.com/server_api", bytes.NewReader(reqBytes))
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
		return c.JSON(fiber.Map{"code": "A40005", "msg": "音频转写失败", "data": nil})
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

	aiResponse := ai.Run(prompt, setting.SettingVar.AIPort)
	if aiResponse == "err" {
		return c.JSON(fiber.Map{"code": "A40005", "msg": "音频分析失败", "data": nil})
	}

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
	if err := json.Unmarshal([]byte(aiResponse), &result); err != nil {
		return c.JSON(fiber.Map{"code": "A40005", "msg": "音频分析失败", "data": nil})
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
