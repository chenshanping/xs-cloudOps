package ai

import (
	"errors"
	"strings"
	"time"

	"server/global"
	"server/model"
)

const adminTimeLayout = "2006-01-02 15:04:05"

// AdminListAIUsers 管理员视角分页查询有 AI 对话记录的用户列表（page 分页，按最近活动倒序）
func (s *AIService) AdminListAIUsers(input AdminAIUserListInput) ([]AdminAIUserItem, int64, error) {
	input = input.Normalize()

	// 子查询：聚合每个 user 的对话数 + 最后活动时间
	sub := global.DB.Table("ai_conversations").
		Select("user_id, COUNT(*) AS conversation_count, MAX(updated_at) AS last_active_at").
		Where("deleted_at IS NULL").
		Group("user_id")

	db := global.DB.Table("(?) AS agg", sub).
		Joins("JOIN sys_user u ON u.id = agg.user_id AND u.deleted_at IS NULL")

	if keyword := strings.TrimSpace(input.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("u.username LIKE ? OR u.nickname LIKE ?", like, like)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	type row struct {
		ID                uint
		Username          string
		Nickname          string
		ConversationCount int64
		LastActiveAt      *time.Time
	}

	var rows []row
	if err := db.
		Select("u.id AS id, u.username AS username, u.nickname AS nickname, agg.conversation_count AS conversation_count, agg.last_active_at AS last_active_at").
		Order("agg.last_active_at DESC, u.id DESC").
		Offset(input.Offset()).
		Limit(input.PageSize).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]AdminAIUserItem, 0, len(rows))
	for _, r := range rows {
		item := AdminAIUserItem{
			ID:                r.ID,
			Username:          r.Username,
			Nickname:          r.Nickname,
			ConversationCount: r.ConversationCount,
		}
		if r.LastActiveAt != nil {
			item.LastActiveAt = r.LastActiveAt.Format(adminTimeLayout)
		}
		items = append(items, item)
	}

	return items, total, nil
}

// AdminListConversations 管理员视角分页查询所有对话（cursor 分页）
// 排序：id DESC（最近创建在前）
// 返回：列表、next_cursor（0 表示无更多）、has_more
func (s *AIService) AdminListConversations(input AdminConversationListInput) ([]AdminConversationItem, uint, bool, error) {
	input.CursorInput = input.CursorInput.Normalize()

	db := global.DB.Model(&model.AIConversation{})
	if input.UserID > 0 {
		db = db.Where("user_id = ?", input.UserID)
	}
	if keyword := strings.TrimSpace(input.Keyword); keyword != "" {
		db = db.Where("title LIKE ?", "%"+keyword+"%")
	}
	if input.Cursor > 0 {
		db = db.Where("id < ?", input.Cursor)
	}

	// 多取一条用于判断 has_more
	var conversations []model.AIConversation
	if err := db.Order("id DESC").Limit(input.Limit + 1).Find(&conversations).Error; err != nil {
		return nil, 0, false, err
	}

	hasMore := len(conversations) > input.Limit
	if hasMore {
		conversations = conversations[:input.Limit]
	}

	if len(conversations) == 0 {
		return []AdminConversationItem{}, 0, false, nil
	}

	// 收集 user_id 和 conversation_id，批量查询附加信息
	userIDs := make([]uint, 0, len(conversations))
	conversationIDs := make([]uint, 0, len(conversations))
	seenUser := make(map[uint]struct{}, len(conversations))
	for _, c := range conversations {
		conversationIDs = append(conversationIDs, c.ID)
		if _, ok := seenUser[c.UserID]; !ok {
			userIDs = append(userIDs, c.UserID)
			seenUser[c.UserID] = struct{}{}
		}
	}

	// 批量查询用户信息
	type userBrief struct {
		ID       uint
		Username string
		Nickname string
	}
	users := make([]userBrief, 0, len(userIDs))
	if len(userIDs) > 0 {
		if err := global.DB.Model(&model.SysUser{}).
			Select("id", "username", "nickname").
			Where("id IN ?", userIDs).
			Find(&users).Error; err != nil {
			return nil, 0, false, err
		}
	}
	userMap := make(map[uint]userBrief, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	// 批量统计每个对话的消息数
	type messageCountRow struct {
		ConversationID uint
		Cnt            int64
	}
	counts := make([]messageCountRow, 0, len(conversationIDs))
	if err := global.DB.Model(&model.AIMessage{}).
		Select("conversation_id, COUNT(*) AS cnt").
		Where("conversation_id IN ?", conversationIDs).
		Group("conversation_id").
		Find(&counts).Error; err != nil {
		return nil, 0, false, err
	}
	countMap := make(map[uint]int64, len(counts))
	for _, c := range counts {
		countMap[c.ConversationID] = c.Cnt
	}

	items := make([]AdminConversationItem, 0, len(conversations))
	for _, c := range conversations {
		u := userMap[c.UserID]
		item := AdminConversationItem{
			ID:           c.ID,
			UserID:       c.UserID,
			Username:     u.Username,
			Nickname:     u.Nickname,
			Title:        c.Title,
			Model:        c.Model,
			MessageCount: countMap[c.ID],
			CreatedAt:    c.CreatedAt.Format(adminTimeLayout),
			UpdatedAt:    c.UpdatedAt.Format(adminTimeLayout),
		}
		if c.ContextClearedAt != nil {
			item.ContextClearedAt = c.ContextClearedAt.Format(adminTimeLayout)
		}
		items = append(items, item)
	}

	nextCursor := uint(0)
	if hasMore && len(items) > 0 {
		nextCursor = items[len(items)-1].ID
	}

	return items, nextCursor, hasMore, nil
}

// AdminListMessages 管理员视角分页查询某对话的消息（cursor 分页）
// 排序：id ASC（按时间正序展示）
// 返回：列表、next_cursor（0 表示无更多）、has_more
func (s *AIService) AdminListMessages(conversationID uint, input CursorInput) ([]model.AIMessage, uint, bool, error) {
	if conversationID == 0 {
		return nil, 0, false, errors.New("对话不存在")
	}

	var conversation model.AIConversation
	if err := global.DB.First(&conversation, conversationID).Error; err != nil {
		return nil, 0, false, errors.New("对话不存在")
	}

	input = input.Normalize()
	db := global.DB.Model(&model.AIMessage{}).Where("conversation_id = ?", conversationID)
	if input.Cursor > 0 {
		db = db.Where("id > ?", input.Cursor)
	}

	var messages []model.AIMessage
	if err := db.Order("id ASC").Limit(input.Limit + 1).Find(&messages).Error; err != nil {
		return nil, 0, false, err
	}

	hasMore := len(messages) > input.Limit
	if hasMore {
		messages = messages[:input.Limit]
	}

	nextCursor := uint(0)
	if hasMore && len(messages) > 0 {
		nextCursor = messages[len(messages)-1].ID
	}

	return messages, nextCursor, hasMore, nil
}

// AdminDeleteConversation 管理员删除任意对话（含消息）
func (s *AIService) AdminDeleteConversation(conversationID uint) error {
	var conversation model.AIConversation
	if err := global.DB.First(&conversation, conversationID).Error; err != nil {
		return errors.New("对话不存在")
	}

	if err := global.DB.Where("conversation_id = ?", conversationID).Delete(&model.AIMessage{}).Error; err != nil {
		return err
	}

	return global.DB.Delete(&conversation).Error
}
