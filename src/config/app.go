package config

import (
	"os"

	"github.com/Yalm/nestjs-controller-file-finder/src/utils"
)

type AppConfig struct {
	EnableCors                bool
	BackendUrl                string
	RestApiId                 string
	AccessControlAllowOrigin  string
	AccessControlAllowMethods string
	AccessControlAllowHeaders string
	StageName                 string
	RootDir                   string
}

func GetAppConfig() *AppConfig {
	return &AppConfig{
		RootDir:                   utils.Getenv("ROOT_DIR", "./src"),
		EnableCors:                utils.GetBoolenv("ENABLE_CORS", "true"),
		BackendUrl:                os.Getenv("BACKEND_URL"),
		RestApiId:                 os.Getenv("REST_API_ID"),
		StageName:                 utils.Getenv("STAGE_NAME", "V1"),
		AccessControlAllowOrigin:  utils.Getenv("ACCESS_CONTROL_ALLOW_ORIGIN", "*"),
		AccessControlAllowMethods: utils.Getenv("ACCESS_CONTROL_ALLOW_METHODS", "GET,POST,PUT,DELETE,PATCH"),
		AccessControlAllowHeaders: utils.Getenv("ACCESS_CONTROL_ALLOW_HEADERS", "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token"),
	}
}
