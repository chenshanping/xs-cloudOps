package ai

import (
	"regexp"
	"strings"
)

const (
	providerKindOpenAICompatible = "openai-compatible"
	providerKindOpenAI           = "openai"
	providerKindAzureOpenAI      = "azure-openai"
	providerKindGemini           = "gemini"
	providerKindVertex           = "vertex"
	providerKindAnthropic        = "anthropic"
	providerKindDashScope        = "dashscope"
	providerKindPerplexity       = "perplexity"
	providerKindOpenRouter       = "openrouter"
	providerKindAIHubMix         = "aihubmix"
	providerKindGrok             = "grok"
	providerKindHunyuan          = "hunyuan"
	providerKindZhipu            = "zhipu"
	providerKindDoubao           = "doubao"
	providerKindXiaomiMimo       = "xiaomi-mimo"
)

type providerCapabilityContext struct {
	ProviderKind string
}

var (
	cherryEmbeddingRegex             = regexp.MustCompile(`(?:^text-|embed|bge-|e5-|llm2vec|retrieval|uae-|gte-|jina-clip|jina-embeddings|voyage-)`)
	cherryRerankRegex                = regexp.MustCompile(`(?:rerank|re-rank|re-ranker|re-ranking|retrieval|retriever)`)
	cherryReasoningRegex             = regexp.MustCompile(`(?:^o\d+(?:-[\w-]+)?$|\b(?:reasoning|reasoner|thinking|think)\b|-[rR]\d+|\bqwq(?:-[\w-]+)?\b|\bhunyuan-t1(?:-[\w-]+)?\b|\bglm-zero-preview\b|\bgrok-(?:3-mini|4|4-fast)(?:-[\w-]+)?\b)`)
	cherryVisionRegex                = regexp.MustCompile(`(?:llava|moondream|minicpm|gemini-1\.5|gemini-2\.0|gemini-2\.5|gemini-3(?:\.\d)?-(?:flash|pro)(?:-preview)?|gemini-(?:flash|pro|flash-lite)-latest|gemini-exp|claude-3|claude-haiku-4|claude-sonnet-4|claude-opus-4|vision|glm-4(?:\.\d+)?v(?:-[\w-]+)?|qwen-vl|qwen2-vl|qwen2.5-vl|qwen3-vl|qwen3\.[5-9](?:-[\w-]+)?|qwen2.5-omni|qwen3-omni(?:-[\w-]+)?|qvq|internvl2|grok-vision-beta|grok-4(?:-[\w-]+)?|pixtral|gpt-4(?:-[\w-]+)|gpt-4\.1(?:-[\w-]+)?|gpt-4o(?:-[\w-]+)?|gpt-4\.5(?:-[\w-]+)?|gpt-5(?:-[\w-]+)?|chatgpt-4o(?:-[\w-]+)?|o1(?:-[\w-]+)?|o3(?:-[\w-]+)?|o4(?:-[\w-]+)?|deepseek-vl(?:[\w-]+)?|kimi-k2\.[56](?:-[\w-]+)?|kimi-latest|gemma-?[3-4](?:[-.\w]+)?|doubao-seed-1[.-][68](?:-[\w-]+)?|doubao-seed-2[.-]0(?:-[\w-]+)?|doubao-seed-code(?:-[\w-]+)?|kimi-thinking-preview|gemma3(?:[-:\w]+)?|kimi-vl-a3b-thinking(?:-[\w-]+)?|llama-guard-4(?:-[\w-]+)?|llama-4(?:-[\w-]+)?|step-1o(?:.*vision)?|step-1v(?:-[\w-]+)?|qwen-omni(?:-[\w-]+)?|mistral-large-(?:2512|latest)|mistral-medium-(?:2508|latest)|mistral-small|mimo-v2\.5$|mimo-v2-omni(?:-[\w-]+)?|glm-5v-turbo)`)
	cherryFunctionCallingRegex       = regexp.MustCompile(`(?:gpt-4o|gpt-4o-mini|gpt-4|gpt-4\.5|gpt-oss(?:-[\w-]+)?|gpt-5(?:-[0-9-]+)?|o(?:1|3|4)(?:-[\w-]+)?|claude|qwen|qwen3|hunyuan|deepseek|glm-4(?:-[\w-]+)?|glm-4\.5(?:-[\w-]+)?|glm-4\.7(?:-[\w-]+)?|glm-5(?:-[\w-]+)?|learnlm(?:-[\w-]+)?|gemini(?:-[\w-]+)?|gemma-?4(?:[-.\w]+)?|grok-3(?:-[\w-]+)?|grok-4(?:-[\w-]+)?|doubao-seed-1[.-][68](?:-[\w-]+)?|doubao-seed-2[.-]0(?:-[\w-]+)?|doubao-seed-code(?:-[\w-]+)?|kimi-k2(?:-[\w-]+)?|ling-\w+(?:-[\w-]+)?|ring-\w+(?:-[\w-]+)?|minimax-m2(?:\.\d+)?(?:-[\w-]+)?|mimo-v2(?:\.5)?(?:-(?:pro|flash|omni))?$|glm-5v-turbo)`)
	cherryClaudeWebSearchRegex       = regexp.MustCompile(`\b(?:claude-3(-|\.)(7|5)-sonnet(?:-[\w-]+)|claude-3(-|\.)5-haiku(?:-[\w-]+)|claude-(haiku|sonnet|opus)-4(?:-[\w-]+)?)\b`)
	cherryGeminiSearchRegex          = regexp.MustCompile(`gemini-(?:2.*(?:-latest)?|3(?:\.\d+)?-(?:flash|pro)(?:-(?:image-)?preview)?|flash-latest|pro-latest|flash-lite-latest)(?:-[\w-]+)*$`)
	cherryQwenMaxReasoningRegex      = regexp.MustCompile(`^(?:qwen3-max|qwen-max-latest)(?:-|$)`)
	cherryQwenPlusReasoningRegex     = regexp.MustCompile(`^qwen(?:3\.[5-9])?-plus(?:-|$)`)
	cherryQwenFlashReasoningRegex    = regexp.MustCompile(`^qwen(?:3\.[5-9])?-flash(?:-|$)`)
	cherryQwenTurboReasoningRegex    = regexp.MustCompile(`^qwen(?:3\.[5-9])?-turbo(?:-|$)`)
	cherryQwen35SeriesReasoningRegex = regexp.MustCompile(`^qwen3\.[5-9]`)
	cherryQwen3OpenReasoningRegex    = regexp.MustCompile(`^qwen3-\d`)
	cherryKimiReasoningRegex         = regexp.MustCompile(`^kimi-k2-thinking(?:-turbo)?$|^kimi-k(?:2\.[5-9]\d*|[3-9]\d*)(?:[.-]\w+)*$`)
)

var cherryPerplexitySearchModels = map[string]struct{}{
	"sonar":               {},
	"sonar-pro":           {},
	"sonar-reasoning":     {},
	"sonar-reasoning-pro": {},
	"sonar-deep-research": {},
}

func buildProviderCapabilityContext(providerName, baseURL string) providerCapabilityContext {
	source := strings.ToLower(strings.TrimSpace(providerName + " " + baseURL))

	switch {
	case containsAny(source, "openrouter"):
		return providerCapabilityContext{ProviderKind: providerKindOpenRouter}
	case containsAny(source, "perplexity", "sonar"):
		return providerCapabilityContext{ProviderKind: providerKindPerplexity}
	case containsAny(source, "aihubmix"):
		return providerCapabilityContext{ProviderKind: providerKindAIHubMix}
	case containsAny(source, "dashscope", "aliyuncs", "百炼"):
		return providerCapabilityContext{ProviderKind: providerKindDashScope}
	case containsAny(source, "anthropic", "claude"):
		return providerCapabilityContext{ProviderKind: providerKindAnthropic}
	case containsAny(source, "vertex", "aiplatform", "googleapis.com"):
		return providerCapabilityContext{ProviderKind: providerKindVertex}
	case containsAny(source, "generativelanguage", "gemini", "google"):
		return providerCapabilityContext{ProviderKind: providerKindGemini}
	case containsAny(source, "openai.azure.com", "/openai/"):
		return providerCapabilityContext{ProviderKind: providerKindAzureOpenAI}
	case containsAny(source, "api.openai.com", "openai", "chatgpt"):
		return providerCapabilityContext{ProviderKind: providerKindOpenAI}
	case containsAny(source, "x.ai", "grok"):
		return providerCapabilityContext{ProviderKind: providerKindGrok}
	case containsAny(source, "hunyuan", "腾讯混元", "tencentcloudapi"):
		return providerCapabilityContext{ProviderKind: providerKindHunyuan}
	case containsAny(source, "bigmodel", "zhipu", "chatglm", "智谱"):
		return providerCapabilityContext{ProviderKind: providerKindZhipu}
	case containsAny(source, "doubao", "volces", "火山"):
		return providerCapabilityContext{ProviderKind: providerKindDoubao}
	case containsAny(source, "xiaomimimo", "xiaomi mimo", "mimo api", "mimo开放平台", "mimo 开放平台"):
		return providerCapabilityContext{ProviderKind: providerKindXiaomiMimo}
	default:
		return providerCapabilityContext{ProviderKind: providerKindOpenAICompatible}
	}
}

func lowerBaseModelName(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if index := strings.LastIndex(value, "/"); index >= 0 && index+1 < len(value) {
		value = value[index+1:]
	}
	return value
}

func looksLikeCherryEmbeddingModel(modelID string) bool {
	return cherryEmbeddingRegex.MatchString(modelID)
}

func looksLikeCherryRerankModel(modelID string) bool {
	return cherryRerankRegex.MatchString(modelID)
}

func looksLikeCherryReasoningModel(modelID string) bool {
	if strings.Contains(modelID, "-non-reasoning") {
		return false
	}
	return cherryReasoningRegex.MatchString(modelID) ||
		isSupportedThinkingTokenQwenModel(modelID) ||
		isQwenReasoningModel(modelID) ||
		isKimiReasoningModel(modelID) ||
		strings.Contains(modelID, "deepseek-v3.2-speciale")
}

func looksLikeCherryVisionModel(modelID string) bool {
	if isCherryVisionExcluded(modelID) {
		return false
	}
	return cherryVisionRegex.MatchString(modelID)
}

func supportsCherryVisionForProvider(ctx providerCapabilityContext, modelID, modelName string) bool {
	if isProviderVisionExcluded(ctx, modelID, modelName) {
		return false
	}
	return looksLikeCherryVisionModel(modelID) || looksLikeCherryVisionModel(modelName)
}

func looksLikeCherryFunctionCallingModel(ctx providerCapabilityContext, modelID, modelName string) bool {
	if strings.Contains(modelID, "deepseek-v3.1") || strings.Contains(modelID, "deepseek-v3_1") {
		switch ctx.ProviderKind {
		case providerKindDashScope, providerKindDoubao:
			return false
		}
		return true
	}
	if isCherryFunctionCallingExcluded(modelID) {
		return false
	}
	if ctx.ProviderKind == providerKindDoubao || strings.Contains(modelID, "doubao") {
		return cherryFunctionCallingRegex.MatchString(modelName) || cherryFunctionCallingRegex.MatchString(modelID)
	}
	return cherryFunctionCallingRegex.MatchString(modelID)
}

func inferCherrySearchStrategy(ctx providerCapabilityContext, modelID string) string {
	if looksLikeCherryEmbeddingModel(modelID) || looksLikeCherryRerankModel(modelID) || looksLikeCherryTextToImageModel(modelID) {
		return searchStrategyNone
	}

	switch ctx.ProviderKind {
	case providerKindAnthropic:
		if cherryClaudeWebSearchRegex.MatchString(modelID) {
			return searchStrategyBuiltin
		}
	case providerKindOpenAI, providerKindAzureOpenAI:
		if isCherryOpenAIWebSearchModel(modelID) {
			return searchStrategyBuiltin
		}
	case providerKindPerplexity:
		if _, ok := cherryPerplexitySearchModels[modelID]; ok {
			return searchStrategyBuiltin
		}
	case providerKindAIHubMix:
		if (!strings.HasSuffix(modelID, "-search") && looksLikeCherryGeminiSearchModel(modelID)) || isCherryOpenAIWebSearchModel(modelID) {
			return searchStrategyBuiltin
		}
	case providerKindOpenAICompatible:
		if looksLikeCherryGeminiSearchModel(modelID) || isCherryOpenAIWebSearchModel(modelID) {
			return searchStrategyBuiltin
		}
	case providerKindGemini, providerKindVertex:
		if looksLikeCherryGeminiSearchModel(modelID) {
			return searchStrategyBuiltin
		}
	case providerKindHunyuan:
		if modelID != "hunyuan-lite" {
			return searchStrategyBuiltin
		}
	case providerKindZhipu:
		return searchStrategyNone
	case providerKindDashScope:
		for _, prefix := range []string{"qwen-turbo", "qwen-max", "qwen-plus", "qwq", "qwen-flash", "qwen3-max"} {
			if strings.HasPrefix(modelID, prefix) {
				return searchStrategyBuiltin
			}
		}
	case providerKindOpenRouter:
		return searchStrategyBuiltin
	case providerKindGrok:
		return searchStrategyBuiltin
	}

	return searchStrategyNone
}

func isCherryOpenAIWebSearchModel(modelID string) bool {
	return strings.Contains(modelID, "gpt-4o-search-preview") ||
		strings.Contains(modelID, "gpt-4o-mini-search-preview") ||
		(strings.Contains(modelID, "gpt-4.1") && !strings.Contains(modelID, "gpt-4.1-nano")) ||
		(strings.Contains(modelID, "gpt-4o") && !strings.Contains(modelID, "gpt-4o-image")) ||
		strings.Contains(modelID, "o3") ||
		strings.Contains(modelID, "o4") ||
		(strings.Contains(modelID, "gpt-5") && !strings.Contains(modelID, "chat"))
}

func looksLikeCherryTextToImageModel(modelID string) bool {
	return containsAny(modelID,
		"dall-e", "gpt-image", "grok-2-image", "imagen", "flux", "stable-diffusion", "stabilityai", "sdxl",
		"cogview", "qwen-image", "janus", "midjourney", "seedream", "kandinsky", "hunyuanimage")
}

func isQwenReasoningModel(modelID string) bool {
	if strings.HasPrefix(modelID, "qwen3") && strings.Contains(modelID, "thinking") {
		return true
	}
	return isSupportedThinkingTokenQwenModel(modelID) || strings.Contains(modelID, "qwq") || strings.Contains(modelID, "qvq")
}

func isSupportedThinkingTokenQwenModel(modelID string) bool {
	if modelID == "" {
		return false
	}
	if containsAny(modelID, "coder", "asr", "tts", "reranker", "embedding", "instruct", "thinking") {
		return false
	}
	if cherryQwen35SeriesReasoningRegex.MatchString(modelID) {
		return true
	}
	if strings.Contains(modelID, "qwen3-max-2025-09-23") {
		return false
	}
	return cherryQwenMaxReasoningRegex.MatchString(modelID) ||
		cherryQwenPlusReasoningRegex.MatchString(modelID) ||
		cherryQwenFlashReasoningRegex.MatchString(modelID) ||
		cherryQwenTurboReasoningRegex.MatchString(modelID) ||
		cherryQwen3OpenReasoningRegex.MatchString(modelID)
}

func isKimiReasoningModel(modelID string) bool {
	return cherryKimiReasoningRegex.MatchString(modelID)
}

func isProviderVisionExcluded(ctx providerCapabilityContext, modelID, modelName string) bool {
	if ctx.ProviderKind == providerKindDashScope && (strings.Contains(modelID, "kimi-k2.6") || strings.Contains(modelName, "kimi-k2.6")) {
		return true
	}
	return false
}

func looksLikeCherryGeminiSearchModel(modelID string) bool {
	if strings.Contains(modelID, "-image-preview") {
		return false
	}
	return cherryGeminiSearchRegex.MatchString(modelID)
}

func isCherryVisionExcluded(modelID string) bool {
	return containsAny(modelID,
		"gpt-4-turbo-preview",
		"gpt-4-32k",
		"o1-mini",
		"o3-mini",
		"o1-preview",
		"aidc-ai/marco-o1",
	)
}

func isCherryFunctionCallingExcluded(modelID string) bool {
	return containsAny(modelID,
		"aqa",
		"imagen",
		"o1-mini",
		"o1-preview",
		"aidc-ai/marco-o1",
		"qwen-mt",
		"gpt-5-chat",
		"glm-4.5v",
		"gemini-2.5-flash-image",
		"gemini-2.0-flash-preview-image-generation",
		"gemini-3-pro-image",
		"deepseek-v3.2-speciale",
	)
}
