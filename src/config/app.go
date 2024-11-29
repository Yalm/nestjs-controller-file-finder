package config

import (
	"errors"
	"os"
	"strings"

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
	IgnoredPaths              map[string]bool
}

func GetAppConfig() (*AppConfig, error) {
	ignoredPaths := make(map[string]bool)

	for _, path := range strings.Split(utils.Getenv("IGNORED_PATHS", "/health"), ",") {
		if path != "" {
			ignoredPaths[path] = true
		}
	}

	appConfig := &AppConfig{
		IgnoredPaths:              ignoredPaths,
		RootDir:                   utils.Getenv("ROOT_DIR", "./src"),
		EnableCors:                utils.GetBoolenv("ENABLE_CORS", "true"),
		BackendUrl:                os.Getenv("BACKEND_URL"),
		RestApiId:                 os.Getenv("REST_API_ID"),
		StageName:                 utils.Getenv("STAGE_NAME", "V1"),
		AccessControlAllowOrigin:  utils.Getenv("ACCESS_CONTROL_ALLOW_ORIGIN", "*"),
		AccessControlAllowMethods: utils.Getenv("ACCESS_CONTROL_ALLOW_METHODS", "GET,POST,PUT,DELETE,PATCH"),
		AccessControlAllowHeaders: utils.Getenv("ACCESS_CONTROL_ALLOW_HEADERS", "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token"),
	}

	if appConfig.BackendUrl == "" {
		return nil, errors.New("BackendUrl is required")
	}

	if appConfig.RestApiId == "" {
		return nil, errors.New("RestApiId is required")
	}

	return appConfig, nil
}
