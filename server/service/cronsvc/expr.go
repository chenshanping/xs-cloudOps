package cronsvc

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

var cronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

func ValidateCronExpr(expr string) (*time.Time, error) {
	schedule, err := cronParser.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("cron_expr_invalid")
	}
	next := schedule.Next(time.Now())
	return &next, nil
}

func nextRunAt(expr string) *time.Time {
	next, err := ValidateCronExpr(expr)
	if err != nil {
		return nil
	}
	return next
}
