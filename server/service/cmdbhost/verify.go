package cmdbhost

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"server/global"
	"server/model"

	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

func (s *Service) VerifyHost(id uint) error {
	db := global.DB
	if db == nil {
		return errors.New("数据库未初始化")
	}
	var host model.CmdbHost
	if err := db.First(&host, id).Error; err != nil {
		return errors.New("主机不存在")
	}
	var credential model.CmdbSshCredential
	if err := db.First(&credential, host.CredentialID).Error; err != nil {
		return s.markVerifyFailed(db, &host, "SSH凭据不存在")
	}

	clientConfig, err := buildSshClientConfig(credential)
	if err != nil {
		return s.markVerifyFailed(db, &host, err.Error())
	}
	address := net.JoinHostPort(host.SshHost, strconv.Itoa(host.SshPort))
	client, err := ssh.Dial("tcp", address, clientConfig)
	if err != nil {
		return s.markVerifyFailed(db, &host, "SSH连接失败")
	}
	defer client.Close()

	hostname, _ := runRemoteCommand(client, "hostname")
	kernelFamily, _ := runRemoteCommand(client, "uname -s")
	osName, _ := runRemoteCommand(client, `sh -c '. /etc/os-release 2>/dev/null && printf "%s" "${NAME:-$ID}"'`)
	platform, _ := runRemoteCommand(client, `sh -c '. /etc/os-release 2>/dev/null && printf "%s" "${ID:-$NAME}"'`)
	platformVersion, _ := runRemoteCommand(client, `sh -c '. /etc/os-release 2>/dev/null && printf "%s" "${PRETTY_NAME:-${VERSION:-$VERSION_ID}}"'`)
	kernelVersion, _ := runRemoteCommand(client, "uname -r")
	arch, _ := runRemoteCommand(client, "uname -m")
	cpuCoresRaw, _ := runRemoteCommand(client, "nproc")
	memRaw, _ := runRemoteCommand(client, "awk '/MemTotal/ {print $2}' /proc/meminfo")
	legacyRelease, _ := runRemoteCommand(client, "sh -c 'cat /etc/redhat-release 2>/dev/null || head -n 1 /etc/issue 2>/dev/null'")
	privateIPRaw, _ := runRemoteCommand(client, "hostname -I 2>/dev/null || ip -4 -o addr show scope global 2>/dev/null | awk '{print $4}' | cut -d/ -f1 | tr '\\n' ' ' || ifconfig 2>/dev/null | awk '/inet /{print $2}' | tr '\\n' ' '")
	publicIPRaw, _ := runRemoteCommand(client, "curl -fsS -m 3 ifconfig.me 2>/dev/null || curl -fsS -m 3 ipinfo.io/ip 2>/dev/null || curl -fsS -m 3 ip.sb 2>/dev/null")

	resolvedOS, resolvedPlatform, resolvedPlatformVersion := resolvePlatformInfo(
		sanitizeCommandOutput(kernelFamily),
		sanitizeCommandOutput(osName),
		sanitizeCommandOutput(platform),
		sanitizeCommandOutput(platformVersion),
		sanitizeCommandOutput(legacyRelease),
	)

	now := time.Now()
	host.VerifyStatus = model.CmdbHostVerifyStatusSuccess
	host.VerifyMessage = "校验成功"
	host.LastVerifiedAt = &now
	host.Hostname = sanitizeCommandOutput(hostname)
	host.OS = resolvedOS
	host.Platform = resolvedPlatform
	host.PlatformVersion = resolvedPlatformVersion
	host.KernelVersion = sanitizeCommandOutput(kernelVersion)
	host.Architecture = sanitizeCommandOutput(arch)
	if cores, parseErr := strconv.Atoi(sanitizeCommandOutput(cpuCoresRaw)); parseErr == nil {
		host.CpuCores = cores
	}
	if memKB, parseErr := strconv.ParseInt(sanitizeCommandOutput(memRaw), 10, 64); parseErr == nil {
		host.MemoryMB = memKB / 1024
	}
	if strings.TrimSpace(host.PrivateIP) == "" {
		if detected := pickPrivateIPv4(privateIPRaw); detected != "" {
			host.PrivateIP = detected
		}
	}
	if strings.TrimSpace(host.PublicIP) == "" {
		if detected := pickPublicIPv4(sanitizeCommandOutput(publicIPRaw)); detected != "" {
			host.PublicIP = detected
		}
	}
	return db.Save(&host).Error
}

// pickPrivateIPv4 从 `hostname -I` 的输出中挑选第一个非回环、非链路本地的 IPv4 地址。
func pickPrivateIPv4(raw string) string {
	for _, field := range strings.Fields(raw) {
		ip := net.ParseIP(strings.TrimSpace(field))
		if ip == nil {
			continue
		}
		ipv4 := ip.To4()
		if ipv4 == nil {
			continue
		}
		if ipv4.IsLoopback() || ipv4.IsLinkLocalUnicast() || ipv4.IsUnspecified() {
			continue
		}
		return ipv4.String()
	}
	return ""
}

// pickPublicIPv4 校验远端 curl 返回的内容是否为合法 IPv4，避免把错误信息写入字段。
func pickPublicIPv4(raw string) string {
	ip := net.ParseIP(strings.TrimSpace(raw))
	if ip == nil {
		return ""
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		return ""
	}
	if ipv4.IsLoopback() || ipv4.IsPrivate() || ipv4.IsLinkLocalUnicast() || ipv4.IsUnspecified() {
		return ""
	}
	return ipv4.String()
}

func (s *Service) markVerifyFailed(db *gorm.DB, host *model.CmdbHost, message string) error {
	now := time.Now()
	host.VerifyStatus = model.CmdbHostVerifyStatusFailed
	host.VerifyMessage = message
	host.LastVerifiedAt = &now
	return db.Save(host).Error
}

func buildSshClientConfig(credential model.CmdbSshCredential) (*ssh.ClientConfig, error) {
	config := &ssh.ClientConfig{
		User:            credential.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	switch credential.AuthType {
	case model.CmdbCredentialAuthTypePassword:
		config.Auth = []ssh.AuthMethod{ssh.Password(credential.Password)}
	case model.CmdbCredentialAuthTypePrivateKey:
		var signer ssh.Signer
		var err error
		keyBytes := []byte(credential.PrivateKey)
		if credential.Passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte(credential.Passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey(keyBytes)
		}
		if err != nil {
			return nil, fmt.Errorf("SSH私钥解析失败")
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	default:
		return nil, errors.New("不支持的SSH认证方式")
	}
	return config, nil
}

func runRemoteCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var stdout bytes.Buffer
	session.Stdout = &stdout
	runErr := session.Run(command)
	// 即便命令以非零退出码结束（如某些发行版没有 hostname -I），也保留已捕获的 stdout，
	// 让上层根据内容自行判断是否能解析出有效结果。
	return stdout.String(), runErr
}

func sanitizeCommandOutput(value string) string {
	return strings.TrimSpace(strings.ReplaceAll(value, "\x00", ""))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func resolvePlatformInfo(kernelFamily, osName, platform, platformVersion, legacyRelease string) (string, string, string) {
	resolvedOS := firstNonEmpty(osName, platform, kernelFamily)
	resolvedPlatform := firstNonEmpty(normalizePlatformID(platform), normalizePlatformID(osName), kernelFamily)
	resolvedPlatformVersion := firstNonEmpty(platformVersion, osName, platform, legacyRelease, kernelFamily)

	if strings.TrimSpace(osName) != "" && strings.TrimSpace(platformVersion) != "" {
		return resolvedOS, resolvedPlatform, resolvedPlatformVersion
	}

	if normalizedOS, normalizedPlatform, normalizedVersion := parseLegacyRelease(legacyRelease); normalizedOS != "" || normalizedPlatform != "" || normalizedVersion != "" {
		return firstNonEmpty(osName, normalizedOS, kernelFamily),
			firstNonEmpty(normalizePlatformID(platform), normalizedPlatform, normalizePlatformID(osName), kernelFamily),
			firstNonEmpty(platformVersion, normalizedVersion, osName, platform, kernelFamily)
	}

	return resolvedOS, resolvedPlatform, resolvedPlatformVersion
}

func parseLegacyRelease(value string) (string, string, string) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", "", ""
	}
	lower := strings.ToLower(trimmed)
	platform := normalizePlatformID(trimmed)
	osName := mapLegacyOSName(lower, trimmed)
	version := trimmed
	if matched := legacyReleaseVersionRegexp.FindStringSubmatch(trimmed); len(matched) > 1 {
		version = firstNonEmpty(matched[1], matched[2], trimmed)
	}
	return osName, platform, version
}

func normalizePlatformID(value string) string {
	trimmed := strings.TrimSpace(strings.Trim(value, `"'`))
	if trimmed == "" {
		return ""
	}
	lower := strings.ToLower(trimmed)
	switch {
	case strings.Contains(lower, "centos"):
		return "centos"
	case strings.Contains(lower, "red hat"), strings.Contains(lower, "rhel"):
		return "rhel"
	case strings.Contains(lower, "rocky"):
		return "rocky"
	case strings.Contains(lower, "alma"):
		return "almalinux"
	case strings.Contains(lower, "ubuntu"):
		return "ubuntu"
	case strings.Contains(lower, "debian"):
		return "debian"
	case strings.Contains(lower, "sles"), strings.Contains(lower, "opensuse"), strings.Contains(lower, "suse"):
		return "suse"
	case strings.Contains(lower, "alpine"):
		return "alpine"
	case strings.Contains(lower, "openeuler"):
		return "openeuler"
	case strings.Contains(lower, "anolis"):
		return "anolis"
	case strings.Contains(lower, "kylin"):
		return "kylin"
	default:
		return lower
	}
}

func mapLegacyOSName(lower, original string) string {
	switch {
	case strings.Contains(lower, "centos"):
		return "CentOS Linux"
	case strings.Contains(lower, "red hat enterprise linux"), strings.Contains(lower, "red hat"), strings.Contains(lower, "rhel"):
		return "Red Hat Enterprise Linux"
	case strings.Contains(lower, "rocky"):
		return "Rocky Linux"
	case strings.Contains(lower, "alma"):
		return "AlmaLinux"
	case strings.Contains(lower, "ubuntu"):
		return "Ubuntu"
	case strings.Contains(lower, "debian"):
		return "Debian"
	case strings.Contains(lower, "opensuse"):
		return "openSUSE"
	case strings.Contains(lower, "suse"):
		return "SUSE Linux"
	case strings.Contains(lower, "alpine"):
		return "Alpine Linux"
	case strings.Contains(lower, "openeuler"):
		return "openEuler"
	case strings.Contains(lower, "anolis"):
		return "Anolis OS"
	case strings.Contains(lower, "kylin"):
		return "Kylin Linux"
	default:
		return strings.TrimSpace(original)
	}
}

var legacyReleaseVersionRegexp = regexp.MustCompile(`(?i)release\s+([0-9][0-9A-Za-z\.\-_]*)|([0-9]+(?:\.[0-9A-Za-z\-_]+)+)`)
