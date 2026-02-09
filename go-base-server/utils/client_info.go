package utils

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// ClientInfo 客户端信息
type ClientInfo struct {
	Browser  string
	OS       string
	Location string
}

// ParseUserAgent 解析 User-Agent 获取浏览器和操作系统
func ParseUserAgent(ua string) (browser, os string) {
	ua = strings.ToLower(ua)

	// 解析浏览器
	switch {
	case strings.Contains(ua, "edg"):
		browser = "Edge"
	case strings.Contains(ua, "chrome") && !strings.Contains(ua, "edg"):
		browser = "Chrome"
	case strings.Contains(ua, "firefox"):
		browser = "Firefox"
	case strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome"):
		browser = "Safari"
	case strings.Contains(ua, "opera") || strings.Contains(ua, "opr"):
		browser = "Opera"
	case strings.Contains(ua, "msie") || strings.Contains(ua, "trident"):
		browser = "IE"
	default:
		browser = "Unknown"
	}

	// 解析操作系统
	switch {
	case strings.Contains(ua, "windows"):
		os = "Windows"
	case strings.Contains(ua, "mac os"):
		os = "MacOS"
	case strings.Contains(ua, "linux"):
		if strings.Contains(ua, "android") {
			os = "Android"
		} else {
			os = "Linux"
		}
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
		os = "iOS"
	case strings.Contains(ua, "android"):
		os = "Android"
	default:
		os = "Unknown"
	}

	return browser, os
}

// ipAPIResponse ip-api.com 响应结构
type ipAPIResponse struct {
	Status     string `json:"status"`
	Country    string `json:"country"`
	RegionName string `json:"regionName"`
	City       string `json:"city"`
}

// GetIPLocation 通过 IP 获取地理位置
func GetIPLocation(ip string) string {
	// 本地/内网IP直接返回
	if ip == "" || ip == "127.0.0.1" || ip == "::1" || strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "172.") {
		return "本地"
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://ip-api.com/json/" + ip + "?lang=zh-CN")
	if err != nil {
		return "未知"
	}
	defer resp.Body.Close()

	var result ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "未知"
	}

	if result.Status != "success" {
		return "未知"
	}

	// 组合地址
	location := result.Country
	if result.RegionName != "" && result.RegionName != result.Country {
		location += " " + result.RegionName
	}
	if result.City != "" && result.City != result.RegionName {
		location += " " + result.City
	}

	return location
}

// GetClientInfo 获取完整客户端信息
func GetClientInfo(ip, userAgent string) ClientInfo {
	browser, os := ParseUserAgent(userAgent)
	location := GetIPLocation(ip)
	return ClientInfo{
		Browser:  browser,
		OS:       os,
		Location: location,
	}
}
