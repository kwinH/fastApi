package core

import (
	"fastApi/core/logger"
)

func CortInit() {

	ViperInit()

	logger.InitLogger()

	Database()

	RedisInit()

	NsqProducerInit()
}
