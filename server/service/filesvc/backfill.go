package filesvc

import (
	"encoding/json"
	"strings"

	"gorm.io/gorm"

	"server/global"
	"server/model"
)

func (s *ReferenceService) BackfillFileReferences() error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.backfillUserAvatarRefs(tx); err != nil {
			return err
		}
		if err := s.backfillConfigImageRefs(tx); err != nil {
			return err
		}
		if err := s.backfillAIMessageRefs(tx); err != nil {
			return err
		}
		return nil
	})
}

func (s *ReferenceService) backfillUserAvatarRefs(tx *gorm.DB) error {
	var users []model.SysUser
	if err := tx.Select("id", "avatar_file_id").Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		refs := []FileRef{}
		if user.AvatarFileID > 0 {
			refs = append(refs, FileRef{FileID: user.AvatarFileID, Field: "avatar"})
		}
		if err := s.ReplaceRefs(tx, "sys_user", user.ID, refs); err != nil {
			return err
		}
	}
	return nil
}

func (s *ReferenceService) backfillConfigImageRefs(tx *gorm.DB) error {
	var configs []model.SysConfig
	if err := tx.Select("id", "`key`", "value").
		Where("`key` IN ?", imageFileReferenceKeys).
		Find(&configs).Error; err != nil {
		return err
	}

	for _, config := range configs {
		fileID, ok := parseConfigFileID(config.Value)
		refs := []FileRef{}
		if ok {
			refs = append(refs, FileRef{FileID: fileID, Field: config.Key})
		}
		if err := s.ReplaceRefs(tx, "sys_config", config.ID, refs); err != nil {
			return err
		}
	}
	return nil
}

func (s *ReferenceService) backfillAIMessageRefs(tx *gorm.DB) error {
	var messages []model.AIMessage
	if err := tx.Select("id", "file_ids").Find(&messages).Error; err != nil {
		return err
	}

	for _, message := range messages {
		fileIDs := decodeFileIDs(message.FileIDs)
		refs := make([]FileRef, 0, len(fileIDs))
		for _, fileID := range fileIDs {
			refs = append(refs, FileRef{FileID: fileID, Field: "attachment"})
		}
		if err := s.ReplaceRefs(tx, "ai_message", message.ID, refs); err != nil {
			return err
		}
	}
	return nil
}

func decodeFileIDs(raw string) []uint {
	if raw == "" {
		return nil
	}

	var fileIDs []uint
	if err := json.Unmarshal([]byte(raw), &fileIDs); err != nil {
		return nil
	}
	return normalizeIDs(fileIDs)
}

func parseConfigFileID(rawValue string) (uint, bool) {
	rawValue = strings.TrimSpace(rawValue)
	if rawValue == "" || rawValue == "0" {
		return 0, false
	}

	var fileID uint
	for _, ch := range rawValue {
		if ch < '0' || ch > '9' {
			return 0, false
		}
		fileID = fileID*10 + uint(ch-'0')
	}
	if fileID == 0 {
		return 0, false
	}
	return fileID, true
}

var imageFileReferenceKeys = []string{
	"sys_logo_file_id",
	"register_logo_file_id",
	"login_bg_image_file_id",
	"slider_captcha_bg_file_id",
}
