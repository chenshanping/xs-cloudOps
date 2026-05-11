package service

import (
	"server/model"
	"server/service/ai"
	"server/service/apisvc"
	"server/service/auth"
	"server/service/captcha"
	"server/service/configsvc"
	"server/service/core"
	"server/service/dept"
	"server/service/dict"
	"server/service/echart"
	"server/service/email"
	"server/service/file"
	"server/service/logsvc"
	"server/service/menu"
	"server/service/role"
	"server/service/storagesvc"
	"server/service/user"
)

// Service singletons
var (
	User    = user.Default
	Dept    = dept.Default
	Role    = role.Default
	Menu    = menu.Default
	Dict    = dict.Default
	File    = file.Default
	AI      = ai.Default
	Config  = configsvc.Default
	Storage = storagesvc.Default
	Log     = logsvc.Default
	Email   = email.Default
	Cache   = core.Default
	Echart  = echart.Default
	Api     = apisvc.Default
	Captcha = captcha.Default
)

// Auth types
type AuthLoginInput = auth.AuthLoginInput
type AuthFlowError = auth.AuthFlowError
type AuthFlowErrorKind = auth.AuthFlowErrorKind
type AuthLoginResult = auth.AuthLoginResult
type UserInfoCache = core.UserInfoCache

const (
	AuthFlowErrorKindBadRequest   = auth.AuthFlowErrorKindBadRequest
	AuthFlowErrorKindFail         = auth.AuthFlowErrorKindFail
	AuthFlowErrorKindUnauthorized = auth.AuthFlowErrorKindUnauthorized
)

const (
	FileDeleteModeConfigKey      = file.FileDeleteModeConfigKey
	FileDeleteModeLogical        = file.FileDeleteModeLogical
	FileDeleteModePhysical       = file.FileDeleteModePhysical
	PublicConfigKeysConfigKey    = configsvc.PublicConfigKeysConfigKey
	SysLogoFileIDConfigKey       = configsvc.SysLogoFileIDConfigKey
	RegisterLogoFileIDConfigKey  = configsvc.RegisterLogoFileIDConfigKey
	LoginBGImageFileIDConfigKey  = configsvc.LoginBGImageFileIDConfigKey
	StorageTypeConfigKey         = configsvc.StorageTypeConfigKey
	LegacyStorageConfigConfigKey = configsvc.LegacyStorageConfigConfigKey
)

func DefaultPublicConfigKeys() []string {
	return configsvc.DefaultPublicConfigKeys()
}

func DefaultPublicConfigKeysValue() string {
	return configsvc.DefaultPublicConfigKeysValue()
}

func NewAuthFlowService() *auth.AuthFlowService {
	return auth.NewAuthFlowService()
}

// AI types
type ConversationListInput = ai.ConversationListInput
type CreateConversationInput = ai.CreateConversationInput
type AIChatInput = ai.AIChatInput
type AIProviderModelFetchError = ai.AIProviderModelFetchError
type AIStreamEvent = ai.AIStreamEvent
type CursorInput = ai.CursorInput
type AdminConversationListInput = ai.AdminConversationListInput
type AdminConversationItem = ai.AdminConversationItem
type AdminAIUserListInput = ai.AdminAIUserListInput
type AdminAIUserItem = ai.AdminAIUserItem

func NewAIStreamAccumulator() *ai.AIStreamAccumulator {
	return ai.NewAIStreamAccumulator()
}

func StorageConfigKey(storageType model.StorageType) string {
	return storagesvc.StorageConfigKey(storageType)
}

func EnsureDeptManageable(operatorID, deptID uint) error {
	return core.EnsureDeptManageable(operatorID, deptID)
}
