package ai

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"server/model"
	"server/service/storagesvc"
)

func (s *AIService) resolveFileStorage(file model.SysFile) (*model.StorageProfile, error) {
	if strings.TrimSpace(file.StorageType) != "" {
		return storagesvc.Default.GetStorageByType(model.StorageType(file.StorageType))
	}
	return storagesvc.Default.GetDefaultStorage()
}

// 读取文件内容
func (s *AIService) readFileContent(file model.SysFile) (string, error) {
	const maxSize = 100 * 1024 // 100KB

	var data []byte
	var err error

	storage, storageErr := s.resolveFileStorage(file)
	if storageErr == nil && storage != nil && storage.Type == model.StorageTypeLocal {
		var config model.LocalConfig
		if jsonErr := json.Unmarshal([]byte(storage.Config), &config); jsonErr != nil {
			return "", fmt.Errorf("解析存储配置失败: %v", jsonErr)
		}
		fullPath := filepath.Join(config.BasePath, file.Path)
		data, err = os.ReadFile(fullPath)
	} else {
		data, err = s.httpGetFileContent(file.URL)
	}

	if err != nil {
		return "", err
	}

	if len(data) > maxSize {
		return string(data[:maxSize]) + "\n...(文件内容已截断，超过100KB限制)", nil
	}
	return string(data), nil
}

// 本地图片转base64
func (s *AIService) localFileToBase64(file model.SysFile) (string, error) {
	const maxImageSize = 5 * 1024 * 1024 // 5MB

	storage, err := s.resolveFileStorage(file)
	if err != nil {
		return "", err
	}

	var config model.LocalConfig
	if err := json.Unmarshal([]byte(storage.Config), &config); err != nil {
		return "", err
	}
	fullPath := filepath.Join(config.BasePath, file.Path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	if len(data) > maxImageSize {
		return "", fmt.Errorf("图片超过5MB限制")
	}

	mimeType := file.MimeType
	if mimeType == "" {
		ext := strings.ToLower(filepath.Ext(file.Name))
		switch ext {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".png":
			mimeType = "image/png"
		case ".gif":
			mimeType = "image/gif"
		case ".webp":
			mimeType = "image/webp"
		default:
			mimeType = "image/png"
		}
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, b64), nil
}

// HTTP获取文件内容
func (s *AIService) httpGetFileContent(url string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("下载文件失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载文件失败: HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取文件内容失败: %v", err)
	}
	return data, nil
}
