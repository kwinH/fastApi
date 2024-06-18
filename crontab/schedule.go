package crontab

func Schedule() {
	for _, cron := range CronList {
		AddCron(cron)
	}
}
