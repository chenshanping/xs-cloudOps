package cmdbterminal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"

	"server/global"
	"server/model"
	"server/model/request"
)

const (
	defaultIdleTimeoutSeconds = 1800
	defaultUserSessionLimit   = 2
	wsTokenTTL                = time.Minute
	wsTokenKeyPrefix          = "cmdb:terminal:ws:"
	wsPath                    = "/api/v1/cmdb/terminal/ws"
)

type Service struct {
	liveMu       sync.RWMutex
	liveSessions map[uint]*liveSession
}

type wsTokenPayload struct {
	SessionID uint `json:"session_id"`
	UserID    uint `json:"user_id"`
}

type liveSession struct {
	service        *Service
	sessionID      uint
	userID         uint
	sshClient      *ssh.Client
	sshSession     *ssh.Session
	stdin          io.WriteCloser
	wsConn         *websocket.Conn
	writeMu        sync.Mutex
	closeMu        sync.Mutex
	closeOnce      sync.Once
	closeReason    string
	forcedByUserID uint
	idleTimer      *time.Timer
	idleTimeout    time.Duration
	seq            uint64
}

type wsClientMessage struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
	Cols int    `json:"cols,omitempty"`
	Rows int    `json:"rows,omitempty"`
}

type wsServerMessage struct {
	Type      string `json:"type"`
	Data      string `json:"data,omitempty"`
	Status    string `json:"status,omitempty"`
	Message   string `json:"message,omitempty"`
	SessionID uint   `json:"session_id,omitempty"`
}

var Default = NewService()

func NewService() *Service {
	return &Service{
		liveSessions: make(map[uint]*liveSession),
	}
}

func (s *Service) CreateSession(req *request.CreateCmdbTerminalSessionRequest, userID uint, username, clientIP string) (*SessionConnectPayload, error) {
	host, credential, err := s.loadHostAndCredential(req.HostID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureActiveSessionLimit(userID); err != nil {
		return nil, err
	}
	record := model.CmdbTerminalSession{
		HostID:               host.ID,
		UserID:               userID,
		UsernameSnapshot:     username,
		CredentialIDSnapshot: credential.ID,
		ClientIP:             clientIP,
		Status:               model.CmdbTerminalSessionStatusPrepared,
		IdleTimeoutSeconds:   defaultIdleTimeoutSeconds,
	}
	if err := global.DB.Create(&record).Error; err != nil {
		return nil, err
	}
	token, err := s.issueWSToken(record.ID, userID)
	if err != nil {
		return nil, err
	}
	return &SessionConnectPayload{
		SessionID:          record.ID,
		WSToken:            token,
		WSURL:              wsPath + "?token=" + token,
		IdleTimeoutSeconds: defaultIdleTimeoutSeconds,
	}, nil
}

func (s *Service) ListSessions(req *request.CmdbTerminalSessionListRequest, userID uint, roleCodes []string) ([]SessionListItem, int64, error) {
	var list []model.CmdbTerminalSession
	var total int64
	db := global.DB.Model(&model.CmdbTerminalSession{})
	if !hasTerminalAdminPrivileges(roleCodes) {
		db = db.Where("user_id = ?", userID)
	} else if req.UserID != nil {
		db = db.Where("user_id = ?", *req.UserID)
	}
	if req.HostID != nil {
		db = db.Where("host_id = ?", *req.HostID)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("created_at DESC").Offset(req.GetOffset()).Limit(req.PageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	items, err := s.buildSessionItems(list)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *Service) GetSession(id, userID uint, roleCodes []string) (*SessionListItem, error) {
	item, err := s.loadSession(id)
	if err != nil {
		return nil, err
	}
	if !canAccessSession(item.UserID, userID, roleCodes) {
		return nil, errors.New("无权查看该终端会话")
	}
	list, err := s.buildSessionItems([]model.CmdbTerminalSession{*item})
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("终端会话不存在")
	}
	return &list[0], nil
}

func (s *Service) ListLogs(sessionID, userID uint, roleCodes []string, req *request.CmdbTerminalLogListRequest) ([]SessionLogItem, int64, error) {
	session, err := s.loadSession(sessionID)
	if err != nil {
		return nil, 0, err
	}
	if !canAccessSession(session.UserID, userID, roleCodes) {
		return nil, 0, errors.New("无权查看该终端会话日志")
	}
	var list []model.CmdbTerminalLog
	var total int64
	db := global.DB.Model(&model.CmdbTerminalLog{}).Where("session_id = ?", sessionID)
	if req.StreamType != "" {
		db = db.Where("stream_type = ?", req.StreamType)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("seq ASC").Offset(req.GetOffset()).Limit(req.PageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	result := make([]SessionLogItem, 0, len(list))
	for _, item := range list {
		result = append(result, SessionLogItem{
			ID:         item.ID,
			SessionID:  item.SessionID,
			Seq:        item.Seq,
			StreamType: item.StreamType,
			Content:    item.Content,
			CreatedAt:  item.CreatedAt,
		})
	}
	return result, total, nil
}

func (s *Service) DisconnectSession(id, operatorID uint, roleCodes []string, force bool) error {
	session, err := s.loadSession(id)
	if err != nil {
		return err
	}
	if force {
		if !hasTerminalAdminPrivileges(roleCodes) {
			return errors.New("无权强制断开会话")
		}
	} else if session.UserID != operatorID {
		return errors.New("只能断开自己的会话")
	}
	reason := "user_disconnect"
	if force {
		reason = "force_disconnect"
	}
	if s.terminateLiveSession(id, reason, operatorID) {
		return nil
	}
	now := time.Now()
	updates := map[string]interface{}{
		"status":            model.CmdbTerminalSessionStatusClosed,
		"disconnect_reason": reason,
		"end_time":          now,
		"last_activity_at":  now,
	}
	if force {
		updates["forced_by_user_id"] = operatorID
	}
	if err := global.DB.Model(&model.CmdbTerminalSession{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return s.appendStaticLog(id, model.CmdbTerminalStreamTypeSystem, "会话已断开："+reason)
}

func (s *Service) ServeWebSocket(conn *websocket.Conn, token string) error {
	payload, err := s.consumeWSToken(token)
	if err != nil {
		return err
	}
	session, err := s.loadSession(payload.SessionID)
	if err != nil {
		return err
	}
	if session.UserID != payload.UserID {
		return errors.New("终端会话归属异常")
	}
	if session.Status == model.CmdbTerminalSessionStatusClosed {
		return errors.New("终端会话已关闭")
	}
	if err := s.ensureActiveSessionLimit(payload.UserID); err != nil {
		return err
	}
	host, credential, err := s.loadHostAndCredential(session.HostID)
	if err != nil {
		return err
	}
	var hostKeyErr error
	clientConfig, err := buildTerminalSshClientConfig(*credential, func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		fingerprint, callbackErr := s.ensureHostFingerprint(host.ID, key)
		if callbackErr != nil {
			hostKeyErr = callbackErr
		}
		if callbackErr == nil {
			session.HostKeyFingerprint = fingerprint
		}
		return callbackErr
	})
	if err != nil {
		return err
	}
	address := net.JoinHostPort(host.SshHost, strconv.Itoa(host.SshPort))
	sshClient, err := ssh.Dial("tcp", address, clientConfig)
	if err != nil {
		if hostKeyErr != nil {
			_ = s.failPreparedSession(session.ID, "host_key_mismatch", hostKeyErr.Error())
			return hostKeyErr
		}
		_ = s.failPreparedSession(session.ID, "connect_failed", "SSH连接失败")
		return errors.New("SSH连接失败")
	}
	sshSession, err := sshClient.NewSession()
	if err != nil {
		sshClient.Close()
		_ = s.failPreparedSession(session.ID, "session_create_failed", "创建SSH会话失败")
		return errors.New("创建SSH会话失败")
	}
	if err := sshSession.RequestPty("xterm-256color", 24, 80, ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}); err != nil {
		sshSession.Close()
		sshClient.Close()
		_ = s.failPreparedSession(session.ID, "pty_failed", "初始化终端失败")
		return errors.New("初始化终端失败")
	}
	stdin, err := sshSession.StdinPipe()
	if err != nil {
		sshSession.Close()
		sshClient.Close()
		_ = s.failPreparedSession(session.ID, "stdin_failed", "初始化终端输入失败")
		return errors.New("初始化终端输入失败")
	}
	stdout, err := sshSession.StdoutPipe()
	if err != nil {
		stdin.Close()
		sshSession.Close()
		sshClient.Close()
		_ = s.failPreparedSession(session.ID, "stdout_failed", "初始化终端输出失败")
		return errors.New("初始化终端输出失败")
	}
	stderr, err := sshSession.StderrPipe()
	if err != nil {
		stdin.Close()
		sshSession.Close()
		sshClient.Close()
		_ = s.failPreparedSession(session.ID, "stderr_failed", "初始化终端错误输出失败")
		return errors.New("初始化终端错误输出失败")
	}
	if err := sshSession.Shell(); err != nil {
		stdin.Close()
		sshSession.Close()
		sshClient.Close()
		_ = s.failPreparedSession(session.ID, "shell_failed", "启动远程Shell失败")
		return errors.New("启动远程Shell失败")
	}
	now := time.Now()
	session.Status = model.CmdbTerminalSessionStatusActive
	session.StartTime = &now
	session.LastActivityAt = &now
	if err := global.DB.Model(&model.CmdbTerminalSession{}).Where("id = ?", session.ID).Updates(map[string]interface{}{
		"status":               model.CmdbTerminalSessionStatusActive,
		"start_time":           now,
		"last_activity_at":     now,
		"host_key_fingerprint": session.HostKeyFingerprint,
		"disconnect_reason":    "",
		"forced_by_user_id":    0,
	}).Error; err != nil {
		stdin.Close()
		sshSession.Close()
		sshClient.Close()
		return err
	}
	live := &liveSession{
		service:     s,
		sessionID:   session.ID,
		userID:      session.UserID,
		sshClient:   sshClient,
		sshSession:  sshSession,
		stdin:       stdin,
		wsConn:      conn,
		idleTimeout: time.Duration(session.IdleTimeoutSeconds) * time.Second,
	}
	live.idleTimer = time.AfterFunc(live.idleTimeout, func() {
		live.overrideClose("idle_timeout", 0)
		live.cleanup()
	})
	s.registerLiveSession(live)
	if err := live.send(wsServerMessage{Type: "status", Status: "connected", SessionID: session.ID, Message: "终端连接成功"}); err != nil {
		live.cleanup()
		return err
	}
	_ = live.appendLog(model.CmdbTerminalStreamTypeSystem, "终端连接成功")
	go live.pipeOutput(stdout)
	go live.pipeOutput(stderr)
	live.readInputLoop()
	return nil
}

func (s *Service) buildSessionItems(list []model.CmdbTerminalSession) ([]SessionListItem, error) {
	result := make([]SessionListItem, 0, len(list))
	for _, item := range list {
		var host model.CmdbHost
		if err := global.DB.Select("id", "name").First(&host, item.HostID).Error; err != nil {
			return nil, err
		}
		result = append(result, SessionListItem{
			ID:                   item.ID,
			HostID:               item.HostID,
			HostName:             host.Name,
			UserID:               item.UserID,
			UsernameSnapshot:     item.UsernameSnapshot,
			CredentialIDSnapshot: item.CredentialIDSnapshot,
			ClientIP:             item.ClientIP,
			Status:               item.Status,
			StartTime:            item.StartTime,
			EndTime:              item.EndTime,
			IdleTimeoutSeconds:   item.IdleTimeoutSeconds,
			DisconnectReason:     item.DisconnectReason,
			ForcedByUserID:       item.ForcedByUserID,
			HostKeyFingerprint:   item.HostKeyFingerprint,
			LastActivityAt:       item.LastActivityAt,
			CreatedAt:            item.CreatedAt,
			UpdatedAt:            item.UpdatedAt,
		})
	}
	return result, nil
}

func (s *Service) loadHostAndCredential(hostID uint) (*model.CmdbHost, *model.CmdbSshCredential, error) {
	var host model.CmdbHost
	if err := global.DB.First(&host, hostID).Error; err != nil {
		return nil, nil, errors.New("主机不存在")
	}
	var credential model.CmdbSshCredential
	if err := global.DB.First(&credential, host.CredentialID).Error; err != nil {
		return nil, nil, errors.New("主机绑定的SSH凭据不存在")
	}
	return &host, &credential, nil
}

func (s *Service) loadSession(id uint) (*model.CmdbTerminalSession, error) {
	var item model.CmdbTerminalSession
	if err := global.DB.First(&item, id).Error; err != nil {
		return nil, errors.New("终端会话不存在")
	}
	return &item, nil
}

func (s *Service) ensureActiveSessionLimit(userID uint) error {
	var count int64
	if err := global.DB.Model(&model.CmdbTerminalSession{}).
		Where("user_id = ? AND status = ?", userID, model.CmdbTerminalSessionStatusActive).
		Count(&count).Error; err != nil {
		return err
	}
	if count >= defaultUserSessionLimit {
		return errors.New("当前在线终端会话数已达上限")
	}
	return nil
}

func (s *Service) issueWSToken(sessionID, userID uint) (string, error) {
	token := uuid.NewString()
	payload := wsTokenPayload{SessionID: sessionID, UserID: userID}
	buf, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	if err := global.Redis.Set(context.Background(), wsTokenKeyPrefix+token, string(buf), wsTokenTTL).Err(); err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) consumeWSToken(token string) (*wsTokenPayload, error) {
	if token == "" {
		return nil, errors.New("终端连接令牌不能为空")
	}
	key := wsTokenKeyPrefix + token
	value, err := global.Redis.Get(context.Background(), key).Result()
	if err != nil {
		return nil, errors.New("终端连接令牌无效或已过期")
	}
	_ = global.Redis.Del(context.Background(), key).Err()
	var payload wsTokenPayload
	if err := json.Unmarshal([]byte(value), &payload); err != nil {
		return nil, errors.New("终端连接令牌无效")
	}
	return &payload, nil
}

func (s *Service) ensureHostFingerprint(hostID uint, key ssh.PublicKey) (string, error) {
	fingerprint := ssh.FingerprintSHA256(key)
	algorithm := key.Type()
	var item model.CmdbHostSshFingerprint
	now := time.Now()
	if err := global.DB.Where("host_id = ?", hostID).First(&item).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
		record := model.CmdbHostSshFingerprint{
			HostID:         hostID,
			Algorithm:      algorithm,
			Fingerprint:    fingerprint,
			FirstSeenAt:    &now,
			LastVerifiedAt: &now,
		}
		if createErr := global.DB.Create(&record).Error; createErr != nil {
			return "", createErr
		}
		return fingerprint, nil
	}
	if item.Fingerprint != fingerprint || item.Algorithm != algorithm {
		return "", errors.New("主机指纹异常")
	}
	if err := global.DB.Model(&model.CmdbHostSshFingerprint{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
		"last_verified_at": now,
		"updated_at":       now,
	}).Error; err != nil {
		return "", err
	}
	return fingerprint, nil
}

func (s *Service) failPreparedSession(sessionID uint, reason, message string) error {
	now := time.Now()
	if err := global.DB.Model(&model.CmdbTerminalSession{}).Where("id = ?", sessionID).Updates(map[string]interface{}{
		"status":            model.CmdbTerminalSessionStatusFailed,
		"disconnect_reason": reason,
		"end_time":          now,
		"last_activity_at":  now,
	}).Error; err != nil {
		return err
	}
	return s.appendStaticLog(sessionID, model.CmdbTerminalStreamTypeSystem, message)
}

func (s *Service) appendStaticLog(sessionID uint, streamType, content string) error {
	var maxSeq uint64
	_ = global.DB.Model(&model.CmdbTerminalLog{}).Where("session_id = ?", sessionID).Select("COALESCE(MAX(seq),0)").Scan(&maxSeq).Error
	return global.DB.Create(&model.CmdbTerminalLog{
		SessionID:  sessionID,
		Seq:        maxSeq + 1,
		StreamType: streamType,
		Content:    content,
	}).Error
}

func (s *Service) registerLiveSession(live *liveSession) {
	s.liveMu.Lock()
	defer s.liveMu.Unlock()
	s.liveSessions[live.sessionID] = live
}

func (s *Service) unregisterLiveSession(sessionID uint) {
	s.liveMu.Lock()
	defer s.liveMu.Unlock()
	delete(s.liveSessions, sessionID)
}

func (s *Service) terminateLiveSession(sessionID uint, reason string, forcedBy uint) bool {
	s.liveMu.RLock()
	live, ok := s.liveSessions[sessionID]
	s.liveMu.RUnlock()
	if !ok {
		return false
	}
	live.overrideClose(reason, forcedBy)
	live.cleanup()
	return true
}

func (l *liveSession) readInputLoop() {
	for {
		_, data, err := l.wsConn.ReadMessage()
		if err != nil {
			l.overrideClose("ws_closed", 0)
			l.cleanup()
			return
		}
		var msg wsClientMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		switch msg.Type {
		case "input":
			if msg.Data == "" {
				continue
			}
			if _, err := io.WriteString(l.stdin, msg.Data); err != nil {
				l.overrideClose("stdin_write_failed", 0)
				l.cleanup()
				return
			}
			l.touch()
			_ = l.appendLog(model.CmdbTerminalStreamTypeInput, msg.Data)
		case "resize":
			if msg.Rows > 0 && msg.Cols > 0 {
				_ = l.sshSession.WindowChange(msg.Rows, msg.Cols)
				l.touch()
			}
		case "disconnect":
			l.overrideClose("user_disconnect", 0)
			l.cleanup()
			return
		}
	}
}

func (l *liveSession) pipeOutput(reader io.Reader) {
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			chunk := string(buf[:n])
			if sendErr := l.send(wsServerMessage{Type: "output", Data: chunk}); sendErr != nil {
				l.overrideClose("ws_write_failed", 0)
				l.cleanup()
				return
			}
			l.touch()
			_ = l.appendLog(model.CmdbTerminalStreamTypeOutput, chunk)
		}
		if err != nil {
			if err == io.EOF {
				l.overrideClose("ssh_closed", 0)
			} else {
				l.overrideClose("ssh_read_failed", 0)
			}
			l.cleanup()
			return
		}
	}
}

func (l *liveSession) send(msg wsServerMessage) error {
	l.writeMu.Lock()
	defer l.writeMu.Unlock()
	return l.wsConn.WriteJSON(msg)
}

func (l *liveSession) appendLog(streamType, content string) error {
	seq := atomic.AddUint64(&l.seq, 1)
	return global.DB.Create(&model.CmdbTerminalLog{
		SessionID:  l.sessionID,
		Seq:        seq,
		StreamType: streamType,
		Content:    content,
	}).Error
}

func (l *liveSession) touch() {
	now := time.Now()
	if l.idleTimer != nil {
		l.idleTimer.Reset(l.idleTimeout)
	}
	_ = global.DB.Model(&model.CmdbTerminalSession{}).Where("id = ?", l.sessionID).Updates(map[string]interface{}{
		"last_activity_at": now,
	}).Error
}

func (l *liveSession) overrideClose(reason string, forcedBy uint) {
	l.closeMu.Lock()
	defer l.closeMu.Unlock()
	if reason != "" {
		l.closeReason = reason
	}
	if forcedBy > 0 {
		l.forcedByUserID = forcedBy
	}
}

func (l *liveSession) resolveClose() (string, uint) {
	l.closeMu.Lock()
	defer l.closeMu.Unlock()
	if l.closeReason == "" {
		l.closeReason = "closed"
	}
	return l.closeReason, l.forcedByUserID
}

func (l *liveSession) cleanup() {
	l.closeOnce.Do(func() {
		reason, forcedBy := l.resolveClose()
		l.service.unregisterLiveSession(l.sessionID)
		if l.idleTimer != nil {
			l.idleTimer.Stop()
		}
		if l.stdin != nil {
			_ = l.stdin.Close()
		}
		if l.sshSession != nil {
			_ = l.sshSession.Close()
		}
		if l.sshClient != nil {
			_ = l.sshClient.Close()
		}
		if l.wsConn != nil {
			_ = l.wsConn.Close()
		}
		now := time.Now()
		updates := map[string]interface{}{
			"status":            model.CmdbTerminalSessionStatusClosed,
			"disconnect_reason": reason,
			"end_time":          now,
			"last_activity_at":  now,
		}
		if forcedBy > 0 {
			updates["forced_by_user_id"] = forcedBy
		}
		_ = global.DB.Model(&model.CmdbTerminalSession{}).Where("id = ?", l.sessionID).Updates(updates).Error
		_ = l.appendLog(model.CmdbTerminalStreamTypeSystem, "会话已断开："+reason)
	})
}

func buildTerminalSshClientConfig(credential model.CmdbSshCredential, hostKeyCallback ssh.HostKeyCallback) (*ssh.ClientConfig, error) {
	config := &ssh.ClientConfig{
		User:            credential.Username,
		HostKeyCallback: hostKeyCallback,
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

func canAccessSession(sessionUserID, operatorID uint, roleCodes []string) bool {
	return sessionUserID == operatorID || hasTerminalAdminPrivileges(roleCodes)
}

func hasTerminalAdminPrivileges(roleCodes []string) bool {
	for _, code := range roleCodes {
		if code == "admin" || code == "system_admin" {
			return true
		}
	}
	return false
}
