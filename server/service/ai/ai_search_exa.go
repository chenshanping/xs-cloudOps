package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultExaMCPEndpoint       = "https://mcp.exa.ai/mcp?tools=web_search_exa"
	exaMCPTimeout               = 15 * time.Second
	exaSearchMaxResults         = 5
	exaEvidenceSnippetMaxLength = 360
	exaEvidenceMaxLength        = 2400
)

type ExaMCPClient struct {
	Endpoint   string
	HTTPClient *http.Client
	APIKey     string
}

type exaInitializeResult struct {
	ProtocolVersion string `json:"protocolVersion"`
}

type exaToolCallResult struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	IsError bool `json:"isError,omitempty"`
}

type exaJSONRPCResponse[T any] struct {
	Result T `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (s *AIService) getSearcher() AISearcher {
	if s != nil && s.searcher != nil {
		return s.searcher
	}
	return &ExaMCPClient{
		Endpoint: defaultExaMCPEndpoint,
		HTTPClient: &http.Client{
			Timeout: exaMCPTimeout,
		},
	}
}

func (c *ExaMCPClient) Search(ctx context.Context, query string) (*AISearchGrounding, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, errors.New("搜索问题不能为空")
	}

	endpoint := strings.TrimSpace(c.Endpoint)
	if endpoint == "" {
		endpoint = defaultExaMCPEndpoint
	}
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: exaMCPTimeout}
	}

	sessionID, err := c.initialize(ctx, client, endpoint)
	if err != nil {
		return nil, err
	}
	if err := c.initialized(ctx, client, endpoint, sessionID); err != nil {
		return nil, err
	}

	searchResult, err := c.callSearchTool(ctx, client, endpoint, sessionID, query)
	if err != nil {
		return nil, err
	}

	return normalizeExaSearchGrounding(query, searchResult), nil
}

func (c *ExaMCPClient) initialize(ctx context.Context, client *http.Client, endpoint string) (string, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2025-03-26",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "go-base-ai-chat",
				"version": "1.0.0",
			},
		},
	}

	var result exaInitializeResult
	responseSessionID, err := c.postForResult(ctx, client, endpoint, "", payload, &result)
	if err != nil {
		return "", err
	}
	return responseSessionID, nil
}

func (c *ExaMCPClient) initialized(ctx context.Context, client *http.Client, endpoint, sessionID string) error {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "notifications/initialized",
	}

	return c.postNotification(ctx, client, endpoint, sessionID, payload)
}

func (c *ExaMCPClient) callSearchTool(ctx context.Context, client *http.Client, endpoint, sessionID, query string) (string, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "web_search_exa",
			"arguments": map[string]interface{}{
				"query":      query,
				"numResults": exaSearchMaxResults,
			},
		},
	}

	var result exaToolCallResult
	if _, err := c.postForResult(ctx, client, endpoint, sessionID, payload, &result); err != nil {
		return "", err
	}
	if result.IsError {
		return "", errors.New("Exa 搜索返回错误结果")
	}

	var parts []string
	for _, item := range result.Content {
		if item.Type != "text" || strings.TrimSpace(item.Text) == "" {
			continue
		}
		parts = append(parts, strings.TrimSpace(item.Text))
	}
	if len(parts) == 0 {
		return "", errors.New("Exa 搜索未返回可用内容")
	}
	return strings.Join(parts, "\n\n"), nil
}

func (c *ExaMCPClient) postForResult(ctx context.Context, client *http.Client, endpoint, sessionID string, payload map[string]interface{}, target interface{}) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if sessionID != "" {
		req.Header.Set("Mcp-Session-Id", sessionID)
	}
	if strings.TrimSpace(c.APIKey) != "" {
		req.Header.Set("x-api-key", strings.TrimSpace(c.APIKey))
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Exa 搜索请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return "", readExaHTTPError(resp)
	}

	responseSessionID := resp.Header.Get("Mcp-Session-Id")
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.Contains(contentType, "text/event-stream") {
		if err := decodeExaSSEBody(resp.Body, target); err != nil {
			return "", err
		}
		return responseSessionID, nil
	}
	if strings.Contains(contentType, "application/json") {
		if err := decodeExaJSONBody(resp.Body, target); err != nil {
			return "", err
		}
		return responseSessionID, nil
	}

	return "", errors.New("Exa 搜索返回了不支持的响应格式")
}

func (c *ExaMCPClient) postNotification(ctx context.Context, client *http.Client, endpoint, sessionID string, payload map[string]interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if sessionID != "" {
		req.Header.Set("Mcp-Session-Id", sessionID)
	}
	if strings.TrimSpace(c.APIKey) != "" {
		req.Header.Set("x-api-key", strings.TrimSpace(c.APIKey))
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Exa 初始化通知失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return readExaHTTPError(resp)
	}
	return nil
}

func decodeExaSSEBody(reader io.Reader, target interface{}) error {
	scanner := bufio.NewScanner(reader)
	var dataLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" {
			continue
		}
		dataLines = append(dataLines, data)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取 Exa SSE 响应失败: %w", err)
	}
	if len(dataLines) == 0 {
		return errors.New("Exa 搜索未返回有效 SSE 数据")
	}

	return unmarshalExaJSONRPC(strings.Join(dataLines, "\n"), target)
}

func decodeExaJSONBody(reader io.Reader, target interface{}) error {
	body, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("读取 Exa JSON 响应失败: %w", err)
	}
	return unmarshalExaJSONRPC(string(body), target)
}

func unmarshalExaJSONRPC(data string, target interface{}) error {
	switch typed := target.(type) {
	case *exaInitializeResult:
		var envelope exaJSONRPCResponse[exaInitializeResult]
		if err := json.Unmarshal([]byte(data), &envelope); err != nil {
			return fmt.Errorf("解析 Exa 响应失败: %w", err)
		}
		if envelope.Error != nil {
			return errors.New(strings.TrimSpace(envelope.Error.Message))
		}
		*typed = envelope.Result
		return nil
	case *exaToolCallResult:
		var envelope exaJSONRPCResponse[exaToolCallResult]
		if err := json.Unmarshal([]byte(data), &envelope); err != nil {
			return fmt.Errorf("解析 Exa 响应失败: %w", err)
		}
		if envelope.Error != nil {
			return errors.New(strings.TrimSpace(envelope.Error.Message))
		}
		*typed = envelope.Result
		return nil
	default:
		return errors.New("不支持的 Exa 响应类型")
	}
}

func readExaHTTPError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	message := strings.TrimSpace(string(body))
	if message == "" {
		message = http.StatusText(resp.StatusCode)
	}
	return fmt.Errorf("Exa 搜索失败(%d): %s", resp.StatusCode, message)
}

func normalizeExaSearchGrounding(query, raw string) *AISearchGrounding {
	grounding := &AISearchGrounding{
		Query: strings.TrimSpace(query),
	}

	blocks := strings.Split(raw, "\n---")
	for _, block := range blocks {
		source := parseExaSourceBlock(block)
		if source.URL == "" && source.Title == "" && source.Snippet == "" {
			continue
		}
		grounding.Sources = append(grounding.Sources, source)
		if len(grounding.Sources) >= exaSearchMaxResults {
			break
		}
	}

	if len(grounding.Sources) == 0 {
		grounding.Evidence = trimToMax(strings.TrimSpace(raw), exaEvidenceMaxLength)
		return grounding
	}

	var evidenceParts []string
	for index, source := range grounding.Sources {
		lines := []string{
			fmt.Sprintf("来源 %d: %s", index+1, fallbackString(source.Title, source.URL)),
			fmt.Sprintf("- URL: %s", source.URL),
		}
		if source.Published != "" {
			lines = append(lines, fmt.Sprintf("- 发布时间: %s", source.Published))
		}
		if source.Snippet != "" {
			lines = append(lines, fmt.Sprintf("- 摘要: %s", source.Snippet))
		}
		evidenceParts = append(evidenceParts, strings.Join(lines, "\n"))
	}
	grounding.Evidence = trimToMax(strings.Join(evidenceParts, "\n\n"), exaEvidenceMaxLength)
	return grounding
}

func parseExaSourceBlock(block string) AISearchSource {
	var source AISearchSource
	lines := strings.Split(block, "\n")
	var snippetLines []string
	inHighlights := false
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || line == "[...]" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "Title:"):
			inHighlights = false
			source.Title = strings.TrimSpace(strings.TrimPrefix(line, "Title:"))
		case strings.HasPrefix(line, "URL:"):
			inHighlights = false
			source.URL = strings.TrimSpace(strings.TrimPrefix(line, "URL:"))
		case strings.HasPrefix(line, "Published:"):
			inHighlights = false
			published := strings.TrimSpace(strings.TrimPrefix(line, "Published:"))
			if !strings.EqualFold(published, "N/A") {
				source.Published = published
			}
		case strings.HasPrefix(line, "Author:"):
			inHighlights = false
		case strings.HasPrefix(line, "Highlights:"):
			inHighlights = true
			rest := strings.TrimSpace(strings.TrimPrefix(line, "Highlights:"))
			if rest != "" && rest != "[...]" {
				snippetLines = append(snippetLines, rest)
			}
		default:
			if inHighlights {
				snippetLines = append(snippetLines, line)
			}
		}
	}

	source.Snippet = trimToMax(strings.Join(snippetLines, " "), exaEvidenceSnippetMaxLength)
	return source
}

func fallbackString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func trimToMax(text string, max int) string {
	text = strings.TrimSpace(text)
	if max <= 0 || len(text) <= max {
		return text
	}
	if max <= 3 {
		return text[:max]
	}
	return strings.TrimSpace(text[:max-3]) + "..."
}
