package dao

import "gbaseadmin/app/system/internal/dao/internal"

type cronDao struct {
	*internal.CronDao
}

var (
	Cron       = cronDao{internal.NewCronDao()}
	SystemCron = Cron
)
