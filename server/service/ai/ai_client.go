package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	appconfig "server/config"
	"server/global"
	modelresponse "server/model/response"
)

const (
	providerModelsFetchTimeout  = 15 * time.Second
	providerModelsResponseLimit = 1 << 20
)

type AIProviderModelFetchError struct {
	Code    int
	Message string
}

func (e *AIProviderModelFetchError) Error() string {
	return e.Message
}

type openAIModelsResponse struct {
	Data []modelresponse.AIProviderModelItem `json:"data"`
}

func (s *AIService) defaultModelID() (string, error) {
	aiConfig, err := s.GetAdminConfig()
	if err != nil {
		return "", err
	}
	if provider := aiConfig.GetDefaultProvider(); provider != nil && len(provider.Models) > 0 {
		return provider.Models[0].ID, nil
	}
	return "", nil
}

func (s *AIService) resolveProvider(modelName string) (*appconfig.AIProvider, error) {
	aiConfig, err := s.GetAdminConfig()
	if err != nil {
		return nil, err
	}
	provider := aiConfig.GetProviderByModel(modelName)
	if provider == nil {
		provider = aiConfig.GetDefaultProvider()
	}
	if provider == nil {
		return nil, errors.New("未配置AI平台")
	}
	if provider.APIKey == "" {
		return nil, errors.New("未配置API Key")
	}
	return provider, nil
}

func (s *AIService) callAPI(modelName string, messages []ChatMessage, stream bool, enableSearch bool, enableThinking bool) ([]byte, error) {
	provider, err := s.resolveProvider(modelName)
	if err != nil {
		return nil, err
	}

	reqBody := ChatCompletionRequest{
		Model: modelName,
	}
	reqBody = buildChatCompletionRequest(modelName, messages, stream, enableSearch, enableThinking)

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", provider.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+provider.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败: %s", string(body))
	}

	return body, nil
}

func (s *AIService) callStreamAPI(modelName string, messages []ChatMessage, enableSearch bool, enableThinking bool) (io.ReadCloser, error) {
	provider, err := s.resolveProvider(modelName)
	if err != nil {
		return nil, err
	}

	reqBody := ChatCompletionRequest{
		Model: modelName,
	}
	reqBody = buildChatCompletionRequest(modelName, messages, true, enableSearch, enableThinking)

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", provider.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+provider.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API请求失败: %s", string(body))
	}

	return resp.Body, nil
}

func (s *AIService) TestConfig(apiKey, baseURL, model string) error {
	if apiKey == "" {
		return errors.New("未配置API Key")
	}
	if baseURL == "" {
		return errors.New("未配置Base URL")
	}
	if model == "" {
		return errors.New("未选择模型")
	}

	reqBody := ChatCompletionRequest{
		Model: model,
		Messages: []ChatMessage{
			{Role: "user", Content: "你好，请回复OK"},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	global.Log.Info("AI 请求日志:" + baseURL + "/chat/completions" + "\n参数为:" + string(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API请求失败(%d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}
	if len(chatResp.Choices) == 0 {
		return errors.New("API返回空响应")
	}

	return nil
}

func buildChatCompletionRequest(modelName string, messages []ChatMessage, stream bool, enableSearch bool, enableThinking bool) ChatCompletionRequest {
	req := ChatCompletionRequest{
		Model:          modelName,
		Messages:       cloneChatMessages(messages),
		Stream:         stream,
		EnableThinking: enableThinking,
	}

	if !enableSearch {
		return req
	}

	req.Messages = append([]ChatMessage{{
		Role: "system",
		Content: fmt.Sprintf(
			"当前日期是 %s。你已启用联网搜索，并且系统已经提供了本轮检索得到的来源摘要。对于实时信息、新闻、价格、天气、政策、版本发布，以及名单、参赛阵容、赛程、比分、排名、获奖名单等时效性事实问题，必须优先依据这些检索来源回答；如果来源与既有知识冲突，以来源为准。优先采用官网、协会、赛事主办方、权威媒体等来源；如果无法从权威来源确认答案，必须明确说明暂未确认，不能猜测、不能编造、不能用历史经验补全。直接回答问题，不要提及“知识截止日期”、不要说自己无法访问实时信息。",
			time.Now().Format("2006-01-02"),
		),
	}}, req.Messages...)

	return req
}

func cloneChatMessages(messages []ChatMessage) []ChatMessage {
	if len(messages) == 0 {
		return nil
	}
	cloned := make([]ChatMessage, len(messages))
	copy(cloned, messages)
	return cloned
}

func (s *AIService) FetchProviderModels(apiKey, baseURL string) ([]modelresponse.AIProviderModelItem, error) {
	client := &http.Client{Timeout: providerModelsFetchTimeout}
	return s.fetchProviderModelsWithClient(client, apiKey, baseURL)
}

func (s *AIService) fetchProviderModelsWithClient(client *http.Client, apiKey, baseURL string) ([]modelresponse.AIProviderModelItem, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, &AIProviderModelFetchError{Code: 400, Message: "请先填写 API Key"}
	}
	if strings.TrimSpace(baseURL) == "" {
		return nil, &AIProviderModelFetchError{Code: 400, Message: "请先填写 Base URL"}
	}

	candidates, err := buildProviderModelsCandidates(baseURL)
	if err != nil {
		return nil, &AIProviderModelFetchError{Code: 400, Message: "Base URL 格式不正确"}
	}

	var lastErr error
	for index, candidate := range candidates {
		models, fetchErr := fetchProviderModelsFromURL(client, apiKey, candidate)
		if fetchErr == nil {
			return models, nil
		}

		lastErr = fetchErr
		fetchBizErr, ok := fetchErr.(*AIProviderModelFetchError)
		shouldFallback := index < len(candidates)-1 && ok && fetchBizErr.Code == http.StatusNotFound
		if !shouldFallback {
			return nil, fetchErr
		}
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, &AIProviderModelFetchError{Code: 502, Message: "拉取模型列表失败，请稍后重试"}
}

func buildProviderModelsCandidates(baseURL string) ([]string, error) {
	normalized, err := normalizeProviderBaseURL(baseURL)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(strings.ToLower(normalized.Path), "/v1") {
		return []string{appendURLPath(normalized, "models")}, nil
	}
	return []string{
		appendURLPath(normalized, "v1", "models"),
		appendURLPath(normalized, "models"),
	}, nil
}

func normalizeProviderBaseURL(baseURL string) (*url.URL, error) {
	trimmed := strings.TrimSpace(strings.TrimRight(baseURL, "/"))
	if trimmed == "" {
		return nil, errors.New("empty base url")
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, errors.New("invalid base url")
	}
	if parsed.Path == "" {
		parsed.Path = "/"
	}
	return parsed, nil
}

func appendURLPath(base *url.URL, segments ...string) string {
	cloned := *base
	parts := []string{cloned.Path}
	parts = append(parts, segments...)
	cloned.Path = path.Join(parts...)
	if !strings.HasPrefix(cloned.Path, "/") {
		cloned.Path = "/" + cloned.Path
	}
	return cloned.String()
}

func fetchProviderModelsFromURL(client *http.Client, apiKey, endpoint string) ([]modelresponse.AIProviderModelItem, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, &AIProviderModelFetchError{Code: 400, Message: "模型列表地址无效"}
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if isTimeoutError(err) {
			return nil, &AIProviderModelFetchError{Code: 504, Message: "拉取模型列表超时，请稍后重试"}
		}
		return nil, &AIProviderModelFetchError{Code: 502, Message: "拉取模型列表失败，请检查 Base URL 是否可用"}
	}
	defer resp.Body.Close()

	reader := io.LimitReader(resp.Body, providerModelsResponseLimit)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(reader)
		return nil, sanitizeProviderModelsError(resp.StatusCode, string(body))
	}

	var payload openAIModelsResponse
	if err := json.NewDecoder(reader).Decode(&payload); err != nil {
		return nil, &AIProviderModelFetchError{Code: 502, Message: "平台模型列表响应格式无效"}
	}

	models := make([]modelresponse.AIProviderModelItem, 0, len(payload.Data))
	for _, item := range payload.Data {
		if strings.TrimSpace(item.ID) == "" {
			continue
		}
		models = append(models, modelresponse.AIProviderModelItem{
			ID:      strings.TrimSpace(item.ID),
			Object:  strings.TrimSpace(item.Object),
			Created: item.Created,
			OwnedBy: strings.TrimSpace(item.OwnedBy),
		})
	}
	return models, nil
}

func sanitizeProviderModelsError(statusCode int, _ string) error {
	switch statusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return &AIProviderModelFetchError{Code: statusCode, Message: "拉取模型列表失败，请检查 API Key 或平台权限"}
	case http.StatusNotFound:
		return &AIProviderModelFetchError{Code: statusCode, Message: "未找到模型列表接口，请检查 Base URL"}
	case http.StatusTooManyRequests:
		return &AIProviderModelFetchError{Code: 502, Message: "上游平台请求过于频繁，请稍后重试"}
	default:
		if statusCode >= 500 {
			return &AIProviderModelFetchError{Code: 502, Message: "上游平台暂时不可用，请稍后重试"}
		}
		return &AIProviderModelFetchError{Code: 502, Message: fmt.Sprintf("拉取模型列表失败，上游返回异常(%d)", statusCode)}
	}
}

func isTimeoutError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
