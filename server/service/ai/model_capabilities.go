package ai

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"

	appconfig "server/config"
	modelresponse "server/model/response"
)

const (
	searchStrategyNone    = "none"
	searchStrategyBuiltin = "builtin"
	searchStrategyTool    = "tool"
)

var capabilityTagOrder = []string{
	"reasoning",
	"vision",
	"search",
	"free",
	"embedding",
	"rerank",
	"tool",
}

var capabilityTagAliases = map[string]string{
	"reasoning":   "reasoning",
	"inference":   "reasoning",
	"chat":        "reasoning",
	"llm":         "reasoning",
	"text":        "reasoning",
	"推理":          "reasoning",
	"vision":      "vision",
	"visual":      "vision",
	"image":       "vision",
	"multimodal":  "vision",
	"multi_modal": "vision",
	"视觉":          "vision",
	"search":      "search",
	"websearch":   "search",
	"web_search":  "search",
	"online":      "search",
	"联网":          "search",
	"free":        "free",
	"免费":          "free",
	"embedding":   "embedding",
	"embeddings":  "embedding",
	"embed":       "embedding",
	"嵌入":          "embedding",
	"rerank":      "rerank",
	"reranker":    "rerank",
	"重排":          "rerank",
	"tool":        "tool",
	"tools":       "tool",
	"function":    "tool",
	"工具":          "tool",
}

func normalizeAdminModel(input appconfig.AIModel) appconfig.AIModel {
	model := input
	model.ID = strings.TrimSpace(model.ID)
	model.Name = strings.TrimSpace(model.Name)
	model.Description = strings.TrimSpace(model.Description)
	model.SearchStrategy = normalizeSearchStrategy(model.SearchStrategy)
	model.Tags = normalizeCapabilityTags(model.Tags)
	if model.Name == "" {
		model.Name = model.ID
	}

	if !model.IsThinking && capabilityTagSet(model.Tags)["reasoning"] {
		model.IsThinking = true
	}
	if !model.SupportVision && capabilityTagSet(model.Tags)["vision"] {
		model.SupportVision = true
	}
	if !model.SupportTools && capabilityTagSet(model.Tags)["tool"] {
		model.SupportTools = true
	}
	if model.SearchStrategy == searchStrategyTool {
		model.SupportTools = true
	}
	if model.SearchStrategy == searchStrategyNone && capabilityTagSet(model.Tags)["search"] {
		model.SearchStrategy = searchStrategyBuiltin
	}
	if !model.SupportEmbedding && capabilityTagSet(model.Tags)["embedding"] {
		model.SupportEmbedding = true
	}
	if !model.SupportRerank && capabilityTagSet(model.Tags)["rerank"] {
		model.SupportRerank = true
	}
	if !model.IsFree && capabilityTagSet(model.Tags)["free"] {
		model.IsFree = true
	}

	if model.Temperature != nil && (math.IsNaN(*model.Temperature) || math.IsInf(*model.Temperature, 0)) {
		model.Temperature = nil
	}
	if model.ContextWindow != nil && *model.ContextWindow < 0 {
		model.ContextWindow = nil
	}

	model.Tags = buildCapabilityTags(
		model.IsThinking,
		model.SupportVision,
		model.SearchStrategy,
		model.IsFree,
		model.SupportEmbedding,
		model.SupportRerank,
		model.SupportTools,
	)
	return model
}

func normalizeCapabilityTags(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}

	normalized := make([]string, 0, len(tags))
	seen := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		key := canonicalCapabilityTag(tag)
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		normalized = append(normalized, key)
	}
	return normalized
}

func canonicalCapabilityTag(tag string) string {
	normalized := strings.TrimSpace(strings.ToLower(tag))
	if normalized == "" {
		return ""
	}
	normalized = strings.ReplaceAll(normalized, "-", "_")
	normalized = strings.ReplaceAll(normalized, " ", "_")
	if mapped, ok := capabilityTagAliases[normalized]; ok {
		return mapped
	}
	return ""
}

func capabilityTagSet(tags []string) map[string]bool {
	set := make(map[string]bool, len(tags))
	for _, tag := range tags {
		set[tag] = true
	}
	return set
}

func normalizeSearchStrategy(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "", searchStrategyNone:
		return searchStrategyNone
	case searchStrategyBuiltin:
		return searchStrategyBuiltin
	case searchStrategyTool:
		return searchStrategyTool
	default:
		return searchStrategyNone
	}
}

func buildCapabilityTags(
	isThinking bool,
	supportVision bool,
	searchStrategy string,
	isFree bool,
	supportEmbedding bool,
	supportRerank bool,
	supportTools bool,
) []string {
	tags := make([]string, 0, len(capabilityTagOrder))
	if isThinking {
		tags = append(tags, "reasoning")
	}
	if supportVision {
		tags = append(tags, "vision")
	}
	if normalizeSearchStrategy(searchStrategy) == searchStrategyBuiltin {
		tags = append(tags, "search")
	}
	if isFree {
		tags = append(tags, "free")
	}
	if supportEmbedding {
		tags = append(tags, "embedding")
	}
	if supportRerank {
		tags = append(tags, "rerank")
	}
	if supportTools {
		tags = append(tags, "tool")
	}
	return tags
}

func normalizeRemoteProviderModelItem(raw map[string]any, capabilityContext providerCapabilityContext) modelresponse.AIProviderModelItem {
	tags := normalizeCapabilityTags(extractStringSlice(raw, "tags", "capabilities"))
	tagSet := capabilityTagSet(tags)
	searchStrategy := normalizeSearchStrategy(extractString(raw, "search_strategy", "searchStrategy"))

	modelID := strings.TrimSpace(extractString(raw, "id"))
	displayName := strings.TrimSpace(extractString(raw, "name", "display_name", "title", "label"))
	description := strings.TrimSpace(extractString(raw, "description"))
	lowerModelID := lowerBaseModelName(modelID)
	lowerDisplayName := lowerBaseModelName(displayName)
	textBasis := strings.ToLower(strings.Join([]string{
		lowerModelID,
		lowerDisplayName,
		description,
		strings.Join(tags, " "),
	}, " "))

	supportEmbedding := boolFromPayload(raw, []string{"support_embedding", "supports_embedding", "embedding"}, tagSet["embedding"] || looksLikeCherryEmbeddingModel(lowerModelID) || looksLikeCherryEmbeddingModel(lowerDisplayName) || looksLikeEmbeddingModel(textBasis))
	supportRerank := boolFromPayload(raw, []string{"support_rerank", "supports_rerank", "rerank"}, tagSet["rerank"] || looksLikeCherryRerankModel(lowerModelID) || looksLikeCherryRerankModel(lowerDisplayName) || looksLikeRerankModel(textBasis))
	isThinking := boolFromPayload(raw, []string{"is_thinking", "supports_reasoning", "reasoning"}, tagSet["reasoning"] || (!supportEmbedding && !supportRerank && (looksLikeCherryReasoningModel(lowerModelID) || looksLikeCherryReasoningModel(lowerDisplayName) || looksLikeReasoningModel(textBasis))))
	supportVision := boolFromPayload(raw, []string{"support_vision", "supports_vision", "vision"}, tagSet["vision"] || (!supportEmbedding && !supportRerank && (supportsCherryVisionForProvider(capabilityContext, lowerModelID, lowerDisplayName) || looksLikeVisionModel(textBasis))))
	supportTools := boolFromPayload(raw, []string{"support_tools", "supports_tools", "supports_function_call", "function_call", "tool_call"}, tagSet["tool"] || (!supportEmbedding && !supportRerank && (looksLikeCherryFunctionCallingModel(capabilityContext, lowerModelID, strings.ToLower(displayName)) || looksLikeCherryFunctionCallingModel(capabilityContext, lowerDisplayName, strings.ToLower(displayName)) || looksLikeToolCapableModel(textBasis))))
	isFree := boolFromPayload(raw, []string{"is_free", "free"}, tagSet["free"] || strings.Contains(textBasis, "free") || strings.Contains(textBasis, "trial"))
	if searchStrategy == searchStrategyTool {
		supportTools = true
	}
	if searchStrategy == searchStrategyNone {
		switch {
		case supportEmbedding || supportRerank:
			searchStrategy = searchStrategyNone
		case tagSet["search"]:
			searchStrategy = searchStrategyBuiltin
		case inferCherrySearchStrategy(capabilityContext, lowerModelID) == searchStrategyBuiltin:
			searchStrategy = searchStrategyBuiltin
		case inferCherrySearchStrategy(capabilityContext, lowerDisplayName) == searchStrategyBuiltin:
			searchStrategy = searchStrategyBuiltin
		case strings.Contains(textBasis, "websearch") || strings.Contains(textBasis, "web_search") || strings.Contains(textBasis, "联网"):
			searchStrategy = searchStrategyBuiltin
		}
	}

	temperature := extractFloatPointer(raw, "temperature")
	contextWindow := extractIntPointer(raw, "context_window", "contextWindow", "context_length", "max_context", "max_tokens")
	item := modelresponse.AIProviderModelItem{
		ID:               modelID,
		Name:             displayName,
		Description:      description,
		Object:           strings.TrimSpace(extractString(raw, "object")),
		Created:          extractInt64(raw, "created"),
		OwnedBy:          strings.TrimSpace(extractString(raw, "owned_by", "ownedBy")),
		IsThinking:       isThinking,
		SupportVision:    supportVision,
		SupportTools:     supportTools,
		SearchStrategy:   searchStrategy,
		SupportEmbedding: supportEmbedding,
		SupportRerank:    supportRerank,
		IsFree:           isFree,
		Temperature:      temperature,
		ContextWindow:    contextWindow,
	}
	if item.Name == "" {
		item.Name = item.ID
	}
	applyProviderCapabilityOverrides(&item, capabilityContext, lowerModelID)
	item.Tags = buildCapabilityTags(
		item.IsThinking,
		item.SupportVision,
		item.SearchStrategy,
		item.IsFree,
		item.SupportEmbedding,
		item.SupportRerank,
		item.SupportTools,
	)
	return item
}

func applyProviderCapabilityOverrides(item *modelresponse.AIProviderModelItem, capabilityContext providerCapabilityContext, lowerModelID string) {
	if item == nil {
		return
	}

	switch capabilityContext.ProviderKind {
	case providerKindXiaomiMimo:
		applyXiaomiMimoCapabilityOverrides(item, lowerModelID)
	}
}

func applyXiaomiMimoCapabilityOverrides(item *modelresponse.AIProviderModelItem, lowerModelID string) {
	switch {
	case isXiaomiMimoTTSModel(lowerModelID):
		item.IsThinking = false
		item.SupportVision = false
		item.SupportTools = false
		item.SupportEmbedding = false
		item.SupportRerank = false
		item.SearchStrategy = searchStrategyNone
		item.IsFree = true
	case isXiaomiMimoReasoningModel(lowerModelID):
		item.IsThinking = true
		item.SupportVision = false
	case isXiaomiMimoVisionModel(lowerModelID):
		item.SupportVision = true
	}
}

func isXiaomiMimoReasoningModel(modelID string) bool {
	switch modelID {
	case "mimo-v2-pro", "mimo-v2.5-pro":
		return true
	default:
		return false
	}
}

func isXiaomiMimoVisionModel(modelID string) bool {
	switch modelID {
	case "mimo-v2-omni", "mimo-v2.5":
		return true
	default:
		return false
	}
}

func isXiaomiMimoTTSModel(modelID string) bool {
	return strings.HasPrefix(modelID, "mimo-v2-tts") || strings.HasPrefix(modelID, "mimo-v2.5-tts")
}

func boolFromPayload(raw map[string]any, keys []string, fallback bool) bool {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok {
			continue
		}
		normalized, exists := normalizeBoolValue(value)
		if exists {
			return normalized
		}
	}
	return fallback
}

func normalizeBoolValue(value any) (bool, bool) {
	switch typed := value.(type) {
	case bool:
		return typed, true
	case float64:
		return typed != 0, true
	case int:
		return typed != 0, true
	case int64:
		return typed != 0, true
	case json.Number:
		n, err := typed.Int64()
		if err == nil {
			return n != 0, true
		}
	case string:
		switch strings.TrimSpace(strings.ToLower(typed)) {
		case "1", "true", "yes", "y", "on":
			return true, true
		case "0", "false", "no", "n", "off":
			return false, true
		}
	}
	return false, false
}

func extractString(raw map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok || value == nil {
			continue
		}
		switch typed := value.(type) {
		case string:
			if strings.TrimSpace(typed) != "" {
				return typed
			}
		case json.Number:
			return typed.String()
		case float64:
			if typed != 0 {
				return strconv.FormatFloat(typed, 'f', -1, 64)
			}
		}
	}
	return ""
}

func extractStringSlice(raw map[string]any, keys ...string) []string {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok || value == nil {
			continue
		}
		switch typed := value.(type) {
		case []string:
			return typed
		case []any:
			items := make([]string, 0, len(typed))
			for _, item := range typed {
				switch cast := item.(type) {
				case string:
					items = append(items, cast)
				case json.Number:
					items = append(items, cast.String())
				}
			}
			return items
		}
	}
	return nil
}

func extractInt64(raw map[string]any, key string) int64 {
	value, ok := raw[key]
	if !ok || value == nil {
		return 0
	}
	switch typed := value.(type) {
	case float64:
		return int64(typed)
	case int64:
		return typed
	case int:
		return int64(typed)
	case json.Number:
		n, err := typed.Int64()
		if err == nil {
			return n
		}
	case string:
		n, err := strconv.ParseInt(strings.TrimSpace(typed), 10, 64)
		if err == nil {
			return n
		}
	}
	return 0
}

func extractFloatPointer(raw map[string]any, keys ...string) *float64 {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok || value == nil {
			continue
		}
		switch typed := value.(type) {
		case float64:
			if !math.IsNaN(typed) && !math.IsInf(typed, 0) {
				return &typed
			}
		case json.Number:
			f, err := typed.Float64()
			if err == nil && !math.IsNaN(f) && !math.IsInf(f, 0) {
				return &f
			}
		case string:
			f, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
			if err == nil && !math.IsNaN(f) && !math.IsInf(f, 0) {
				return &f
			}
		}
	}
	return nil
}

func extractIntPointer(raw map[string]any, keys ...string) *int {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok || value == nil {
			continue
		}
		switch typed := value.(type) {
		case float64:
			n := int(typed)
			if n >= 0 {
				return &n
			}
		case int:
			if typed >= 0 {
				return &typed
			}
		case int64:
			if typed >= 0 {
				n := int(typed)
				return &n
			}
		case json.Number:
			n, err := typed.Int64()
			if err == nil && n >= 0 {
				value := int(n)
				return &value
			}
		case string:
			n, err := strconv.Atoi(strings.TrimSpace(typed))
			if err == nil && n >= 0 {
				return &n
			}
		}
	}
	return nil
}

func looksLikeEmbeddingModel(text string) bool {
	return containsAny(text, "embedding", "embeddings", "embed", "text-embedding", "bge", "m3e")
}

func looksLikeRerankModel(text string) bool {
	return containsAny(text, "rerank", "reranker")
}

func looksLikeReasoningModel(text string) bool {
	return containsAny(text, "reasoning", "thinking", "think", "deepthink", "r1", "o1", "o3", "qwq")
}

func looksLikeVisionModel(text string) bool {
	if containsAny(text, "embedding", "rerank") {
		return false
	}
	return containsAny(text, "vision", "visual", "omni", "multi-modal", "multimodal", "image", "img", "vl", "vlm")
}

func looksLikeToolCapableModel(text string) bool {
	if containsAny(text, "audio", "speech", "tts", "asr", "moderation", "embedding", "rerank") {
		return false
	}
	return containsAny(text,
		"gpt", "qwen", "glm", "deepseek", "claude", "gemini", "kimi",
		"moonshot", "doubao", "hunyuan", "internlm", "yi", "mistral", "llama", "chat")
}

func containsAny(text string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}
