package ai

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func (s *AIService) prepareSearchGrounding(ctx context.Context, messages []ChatMessage, query string, enableSearch bool) ([]ChatMessage, string, error) {
	if !enableSearch {
		return cloneChatMessages(messages), "", nil
	}

	grounding, err := s.getSearcher().Search(ctx, query)
	if err != nil {
		return nil, "", fmt.Errorf("联网搜索失败: %w", err)
	}

	groundedMessages := cloneChatMessages(messages)
	groundedMessages = append([]ChatMessage{{
		Role:    "system",
		Content: buildSearchGroundingSystemMessage(grounding),
	}}, groundedMessages...)

	return groundedMessages, buildSourcesMarkdown(grounding.Sources), nil
}

func buildSearchGroundingSystemMessage(grounding *AISearchGrounding) string {
	if grounding == nil {
		return ""
	}
	return fmt.Sprintf(
		"当前日期是 %s。系统已通过 Exa 联网搜索为本轮问题检索到以下来源摘要。你回答这类时效性事实问题时，必须优先依据这些来源；如果这些来源仍不足以确认答案，必须明确说明暂未确认，不能猜测、不能编造，也不要提及“知识截止日期”。\n\n检索问题：%s\n\n检索证据：\n%s",
		time.Now().Format("2006-01-02"),
		strings.TrimSpace(grounding.Query),
		strings.TrimSpace(grounding.Evidence),
	)
}

func buildSourcesMarkdown(sources []AISearchSource) string {
	if len(sources) == 0 {
		return ""
	}

	var lines []string
	lines = append(lines, "来源：")
	for index, source := range sources {
		title := fallbackString(source.Title, source.URL)
		line := fmt.Sprintf("%d. [%s](%s)", index+1, title, source.URL)
		if strings.TrimSpace(source.Published) != "" {
			line += fmt.Sprintf("（%s）", strings.TrimSpace(source.Published))
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func appendSourcesToContent(content, sourcesMarkdown string) string {
	content = strings.TrimSpace(content)
	sourcesMarkdown = strings.TrimSpace(sourcesMarkdown)
	if sourcesMarkdown == "" {
		return content
	}
	if content == "" {
		return sourcesMarkdown
	}
	return content + "\n\n---\n" + sourcesMarkdown
}
