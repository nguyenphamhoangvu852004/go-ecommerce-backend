package initialize

import (
	"go-ecommerce-backend-api/global"

	"github.com/gin-gonic/gin"
)

func Run() *gin.Engine {
	LoadConfig()
	InitLogger()
	global.Logger.Info("Load Config Success")
	InitMysqlC()
	InitServiceInterface()
	InitRedis()
	InitKafka()
	r := InitRouter()
	return r
}
