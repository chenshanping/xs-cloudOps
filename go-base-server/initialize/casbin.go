package initialize

import (
	"fmt"
	"go-base-server/global"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

func InitCasbin() {
	adapter, err := gormadapter.NewAdapterByDB(global.DB)
	if err != nil {
		panic(fmt.Errorf("创建Casbin适配器失败: %w", err))
	}

	enforcer, err := casbin.NewEnforcer(global.Config.Casbin.ModelPath, adapter)
	if err != nil {
		panic(fmt.Errorf("创建Casbin执行器失败: %w", err))
	}

	// 注册URL匹配函数
	enforcer.AddFunction("keyMatch2", util.KeyMatch2Func)

	if err := enforcer.LoadPolicy(); err != nil {
		panic(fmt.Errorf("加载Casbin策略失败: %w", err))
	}

	global.Enforcer = enforcer
	global.Log.Info("Casbin初始化成功")
}
