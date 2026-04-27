package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	appconfig "server/config"
	"server/global"
)

func (s *AIService) defaultModelID() string {
	aiConfig := Config.GetAIConfig()
	if provider := aiConfig.GetDefaultProvider(); provider != nil && len(provider.Models) > 0 {
		return provider.Models[0].ID
	}
	return ""
}

func (s *AIService) resolveProvider(modelName string) (*appconfig.AIProvider, error) {
	aiConfig := Config.GetAIConfig()
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

// 调用API（非流式）
func (s *AIService) callAPI(modelName string, messages []ChatMessage, stream bool, enableSearch bool, enableThinking bool) ([]byte, error) {
	provider, err := s.resolveProvider(modelName)
	if err != nil {
		return nil, err
	}

	reqBody := ChatCompletionRequest{
		Model:          modelName,
		Messages:       messages,
		Stream:         stream,
		EnableSearch:   enableSearch,
		EnableThinking: enableThinking,
	}

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

// 调用流式API
func (s *AIService) callStreamAPI(modelName string, messages []ChatMessage, enableSearch bool, enableThinking bool) (io.ReadCloser, error) {
	provider, err := s.resolveProvider(modelName)
	if err != nil {
		return nil, err
	}

	reqBody := ChatCompletionRequest{
		Model:          modelName,
		Messages:       messages,
		Stream:         true,
		EnableSearch:   enableSearch,
		EnableThinking: enableThinking,
	}

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

// 测试AI配置
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
