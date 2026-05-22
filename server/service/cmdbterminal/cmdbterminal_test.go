package cmdbterminal

import (
	"context"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"server/global"
	"server/model"
	"server/model/request"
)

func setupCmdbTerminalServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(
		&model.CmdbHostGroup{},
		&model.CmdbSshCredential{},
		&model.CmdbHost{},
		&model.CmdbTerminalSession{},
		&model.CmdbTerminalLog{},
		&model.CmdbHostSshFingerprint{},
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	redisServer, err := miniredis.Run()
	if err != nil {
		t.Fatalf("start miniredis: %v", err)
	}
	redisClient := redis.NewClient(&redis.Options{Addr: redisServer.Addr()})
	previousDB := global.DB
	previousRedis := global.Redis
	global.DB = db
	global.Redis = redisClient
	t.Cleanup(func() {
		global.DB = previousDB
		global.Redis = previousRedis
		redisClient.Close()
		redisServer.Close()
	})
	return db
}

func seedTerminalHost(t *testing.T) (*model.CmdbHost, *model.CmdbSshCredential) {
	t.Helper()
	group := model.CmdbHostGroup{Name: "默认分组", Status: 1}
	if err := global.DB.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	credential := model.CmdbSshCredential{Name: "默认凭据", AuthType: model.CmdbCredentialAuthTypePassword, Username: "root", Password: "secret"}
	if err := global.DB.Create(&credential).Error; err != nil {
		t.Fatalf("seed credential: %v", err)
	}
	host := model.CmdbHost{Name: "term-01", GroupID: group.ID, SshHost: "127.0.0.1", SshPort: 22, CredentialID: credential.ID}
	if err := global.DB.Create(&host).Error; err != nil {
		t.Fatalf("seed host: %v", err)
	}
	return &host, &credential
}

func TestTerminalService_CreateSession_IssuesWSToken(t *testing.T) {
	setupCmdbTerminalServiceTestDB(t)
	host, credential := seedTerminalHost(t)
	svc := NewService()

	item, err := svc.CreateSession(&request.CreateCmdbTerminalSessionRequest{HostID: host.ID}, 7, "ops", "10.0.0.1")
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if item.SessionID == 0 || item.WSToken == "" {
		t.Fatalf("unexpected session payload: %#v", item)
	}
	var session model.CmdbTerminalSession
	if err := global.DB.First(&session, item.SessionID).Error; err != nil {
		t.Fatalf("load session: %v", err)
	}
	if session.CredentialIDSnapshot != credential.ID {
		t.Fatalf("credential snapshot = %d, want %d", session.CredentialIDSnapshot, credential.ID)
	}
	if _, err := global.Redis.Get(context.Background(), wsTokenKeyPrefix+item.WSToken).Result(); err != nil {
		t.Fatalf("ws token not stored: %v", err)
	}
}

func TestTerminalService_CreateSession_RejectsConcurrentLimit(t *testing.T) {
	setupCmdbTerminalServiceTestDB(t)
	host, credential := seedTerminalHost(t)
	_ = credential
	if err := global.DB.Create(&model.CmdbTerminalSession{
		HostID:               host.ID,
		UserID:               9,
		UsernameSnapshot:     "ops",
		CredentialIDSnapshot: host.CredentialID,
		Status:               model.CmdbTerminalSessionStatusActive,
		IdleTimeoutSeconds:   defaultIdleTimeoutSeconds,
	}).Error; err != nil {
		t.Fatalf("seed session 1: %v", err)
	}
	if err := global.DB.Create(&model.CmdbTerminalSession{
		HostID:               host.ID,
		UserID:               9,
		UsernameSnapshot:     "ops",
		CredentialIDSnapshot: host.CredentialID,
		Status:               model.CmdbTerminalSessionStatusActive,
		IdleTimeoutSeconds:   defaultIdleTimeoutSeconds,
	}).Error; err != nil {
		t.Fatalf("seed session 2: %v", err)
	}
	svc := NewService()
	_, err := svc.CreateSession(&request.CreateCmdbTerminalSessionRequest{HostID: host.ID}, 9, "ops", "10.0.0.1")
	if err == nil || err.Error() != "当前在线终端会话数已达上限" {
		t.Fatalf("expected limit error, got %v", err)
	}
}

func TestTerminalService_ConsumeWSToken_SingleUse(t *testing.T) {
	setupCmdbTerminalServiceTestDB(t)
	svc := NewService()
	token, err := svc.issueWSToken(11, 7)
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}
	payload, err := svc.consumeWSToken(token)
	if err != nil {
		t.Fatalf("consume token: %v", err)
	}
	if payload.SessionID != 11 || payload.UserID != 7 {
		t.Fatalf("payload = %#v", payload)
	}
	if _, err := svc.consumeWSToken(token); err == nil {
		t.Fatalf("expected token reuse to fail")
	}
}

func TestTerminalService_DisconnectSession_RejectsForceForNonAdmin(t *testing.T) {
	setupCmdbTerminalServiceTestDB(t)
	host, _ := seedTerminalHost(t)
	session := model.CmdbTerminalSession{
		HostID:               host.ID,
		UserID:               8,
		UsernameSnapshot:     "owner",
		CredentialIDSnapshot: host.CredentialID,
		Status:               model.CmdbTerminalSessionStatusActive,
		IdleTimeoutSeconds:   defaultIdleTimeoutSeconds,
	}
	if err := global.DB.Create(&session).Error; err != nil {
		t.Fatalf("seed session: %v", err)
	}
	svc := NewService()
	err := svc.DisconnectSession(session.ID, 7, []string{"ops"}, true)
	if err == nil || err.Error() != "无权强制断开会话" {
		t.Fatalf("expected force permission error, got %v", err)
	}
}

func TestTerminalService_DisconnectSession_AllowsOwner(t *testing.T) {
	setupCmdbTerminalServiceTestDB(t)
	host, _ := seedTerminalHost(t)
	session := model.CmdbTerminalSession{
		HostID:               host.ID,
		UserID:               8,
		UsernameSnapshot:     "owner",
		CredentialIDSnapshot: host.CredentialID,
		Status:               model.CmdbTerminalSessionStatusActive,
		IdleTimeoutSeconds:   defaultIdleTimeoutSeconds,
	}
	if err := global.DB.Create(&session).Error; err != nil {
		t.Fatalf("seed session: %v", err)
	}
	svc := NewService()
	if err := svc.DisconnectSession(session.ID, 8, []string{"ops"}, false); err != nil {
		t.Fatalf("disconnect own session: %v", err)
	}
	var updated model.CmdbTerminalSession
	if err := global.DB.First(&updated, session.ID).Error; err != nil {
		t.Fatalf("reload session: %v", err)
	}
	if updated.Status != model.CmdbTerminalSessionStatusClosed {
		t.Fatalf("status = %s, want %s", updated.Status, model.CmdbTerminalSessionStatusClosed)
	}
}
