package config

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/Yalm/nestjs-controller-file-finder/src/utils"
)

type AppConfig struct {
	EnableCors                bool
	VpcLinkId                 string
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

	restApiId := flag.String("rest-api-id", os.Getenv("REST_API_ID"), "Rest API ID")
	vpcLinkId := flag.String("vpc-link-id", os.Getenv("VPC_LINK_ID"), "VPC Link ID")
	backendUrl := flag.String("backend-url", os.Getenv("BACKEND_URL"), "Backend URL")
	flag.Parse()

	appConfig := &AppConfig{
		IgnoredPaths:              ignoredPaths,
		RootDir:                   utils.Getenv("ROOT_DIR", "./src"),
		EnableCors:                utils.GetBoolenv("ENABLE_CORS", "true"),
		BackendUrl:                *backendUrl,
		VpcLinkId:                 *vpcLinkId,
		RestApiId:                 *restApiId,
		StageName:                 utils.Getenv("STAGE_NAME", "V1"),
		AccessControlAllowOrigin:  utils.Getenv("ACCESS_CONTROL_ALLOW_ORIGIN", "*"),
		AccessControlAllowMethods: utils.Getenv("ACCESS_CONTROL_ALLOW_METHODS", "GET,POST,PUT,DELETE,PATCH,OPTIONS"),
		AccessControlAllowHeaders: utils.Getenv("ACCESS_CONTROL_ALLOW_HEADERS", "*"),
	}

	if appConfig.BackendUrl == "" {
		return nil, errors.New("BackendUrl is required")
	}

	if appConfig.RestApiId == "" {
		return nil, errors.New("RestApiId is required")
	}

	return appConfig, nil
}
