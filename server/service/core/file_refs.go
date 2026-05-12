package core

import (
	"server/global"
	"server/model"
)

func FillUserAvatarURL(user *model.SysUser) error {
	if user == nil {
		return nil
	}
	if user.AvatarFileID == 0 {
		user.AvatarFileURL = ""
		return nil
	}

	urlMap, err := getFileURLMap([]uint{user.AvatarFileID})
	if err != nil {
		return err
	}
	user.AvatarFileURL = urlMap[user.AvatarFileID]
	return nil
}

func FillUserAvatarURLs(users []model.SysUser) error {
	if len(users) == 0 {
		return nil
	}

	fileIDs := make([]uint, 0, len(users))
	seen := make(map[uint]struct{}, len(users))
	for _, user := range users {
		if user.AvatarFileID == 0 {
			continue
		}
		if _, exists := seen[user.AvatarFileID]; exists {
			continue
		}
		seen[user.AvatarFileID] = struct{}{}
		fileIDs = append(fileIDs, user.AvatarFileID)
	}

	urlMap, err := getFileURLMap(fileIDs)
	if err != nil {
		return err
	}
	for i := range users {
		users[i].AvatarFileURL = urlMap[users[i].AvatarFileID]
	}
	return nil
}

func getFileURLMap(fileIDs []uint) (map[uint]string, error) {
	result := make(map[uint]string, len(fileIDs))
	if len(fileIDs) == 0 {
		return result, nil
	}

	var files []model.SysFile
	if err := global.DB.Select("id", "url").Where("id IN ? AND status = ?", fileIDs, 1).Find(&files).Error; err != nil {
		return nil, err
	}
	for _, file := range files {
		result[file.ID] = file.URL
	}
	return result, nil
}
