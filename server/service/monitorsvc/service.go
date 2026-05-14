package monitorsvc

import "time"

type MonitorService struct{}

var Default = &MonitorService{}

var processStartedAt = time.Now()
