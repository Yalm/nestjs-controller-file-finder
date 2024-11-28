package main

import (
	"bufio"
	"fmt"

	"github.com/Yalm/nestjs-controller-file-finder/src/config"
	"github.com/Yalm/nestjs-controller-file-finder/src/utils"

	"context"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

type Route struct {
	Method string
	Path   string
}

func main() {
	config := config.GetAppConfig()

	if config.RestApiId == "" {
		log.Fatalf("RestApiId is required")
	}

	if config.BackendUrl == "" {
		log.Fatalf("BackendUrl is required")
	}

	ctx := context.TODO()

	ignoreDirs := map[string]bool{
		"node_modules": true,
		"dist":         true,
		"build":        true,
	}

	controllerRegex := regexp.MustCompile(`@Controller\(['"]?([^'"]*)['"]?\)`)
	methodRegex := regexp.MustCompile(`@(Get|Post|Put|Delete|Patch)\(['"]?([^'"]*)['"]?\)`)

	var extractedRoutes []Route

	log.Println("Searching for routes in", config.RootDir)

	err := filepath.Walk(config.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && ignoreDirs[info.Name()] {
			return filepath.SkipDir
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".controller.ts") {
			extractRoutesFromFile(path, controllerRegex, methodRegex, &extractedRoutes)
		}

		return nil
	})

	if err != nil {
		log.Fatal("Error traversing directories:", err)
	}

	log.Println("Total routes:", len(extractedRoutes))

	sess := utils.NewSession(ctx)
	clientApiGateway := apigateway.New(*sess)

	resources, err := clientApiGateway.GetResources(ctx, &apigateway.GetResourcesInput{
		RestApiId: &config.RestApiId,
	})

	if err != nil {
		log.Fatal("Error getting resources:", err)
	}

	for _, route := range extractedRoutes {
		log.Println("Method:", route.Method, "Path:", route)
		createResources(ctx, clientApiGateway, config, route, &resources.Items)
	}

	_, err = clientApiGateway.CreateDeployment(ctx, &apigateway.CreateDeploymentInput{
		RestApiId: &config.RestApiId,
		StageName: &config.StageName,
	})

	if err != nil {
		log.Fatal("Error creating deployment:", err)
	}
}

func findResource(resources *[]types.Resource, path string) *types.Resource {
	for _, resource := range *resources {
		if resource.PathPart == nil && path != "" {
			continue
		}

		if resource.PathPart == nil && path == "" {
			return &resource
		}

		if *resource.PathPart == path {
			return &resource
		}
	}
	return &types.Resource{}
}

func createCORSMethod(ctx context.Context, clientApiGateway *apigateway.Client, config *config.AppConfig, resourceId *string) {
	log.Println("Creating OPTIONS method for", *resourceId)
	_, err := clientApiGateway.PutMethod(ctx, &apigateway.PutMethodInput{
		RestApiId:         &config.RestApiId,
		ResourceId:        resourceId,
		HttpMethod:        aws.String("OPTIONS"),
		AuthorizationType: aws.String("NONE"),
	})
	if err != nil {
		log.Println("Error creating method:", err)
		return
	}
	_, err = clientApiGateway.PutIntegration(ctx, &apigateway.PutIntegrationInput{
		HttpMethod:            aws.String("OPTIONS"),
		ResourceId:            resourceId,
		RestApiId:             &config.RestApiId,
		Type:                  types.IntegrationTypeMock,
		RequestTemplates:      map[string]string{"application/json": `{"statusCode": 200}`},
		IntegrationHttpMethod: aws.String("OPTIONS"),
	})
	if err != nil {
		log.Println("Error creating integration:", err)
		return
	}

	_, err = clientApiGateway.PutMethodResponse(ctx, &apigateway.PutMethodResponseInput{
		HttpMethod: aws.String("OPTIONS"),
		ResourceId: resourceId,
		RestApiId:  &config.RestApiId,
		StatusCode: aws.String("200"),
		ResponseParameters: map[string]bool{
			"method.response.header.Access-Control-Allow-Headers": true,
			"method.response.header.Access-Control-Allow-Methods": true,
			"method.response.header.Access-Control-Allow-Origin":  true,
		},
		ResponseModels: map[string]string{
			"application/json": "Empty",
		},
	})
	if err != nil {
		log.Println("Error creating method response:", err)
		return
	}

	_, err = clientApiGateway.PutIntegrationResponse(ctx, &apigateway.PutIntegrationResponseInput{
		HttpMethod: aws.String("OPTIONS"),
		ResourceId: resourceId,
		RestApiId:  &config.RestApiId,
		StatusCode: aws.String("200"),
		ResponseParameters: map[string]string{
			"method.response.header.Access-Control-Allow-Headers": fmt.Sprintf("'%s'", config.AccessControlAllowHeaders),
			"method.response.header.Access-Control-Allow-Methods": fmt.Sprintf("'%s'", config.AccessControlAllowMethods),
			"method.response.header.Access-Control-Allow-Origin":  fmt.Sprintf("'%s'", config.AccessControlAllowOrigin),
		},
	})
	if err != nil {
		log.Println("Error creating integration response:", err)
		return
	}
}

func addMethodToResourceById(resources []types.Resource, resourceId string, method string) {
	for index, resource := range resources {
		if *resource.Id == resourceId {
			if resources[index].ResourceMethods == nil {
				resources[index].ResourceMethods = make(map[string]types.Method)
			}
			resources[index].ResourceMethods[method] = types.Method{
				HttpMethod: aws.String(method),
			}
		}
	}
}

func createResources(
	ctx context.Context,
	clientApiGateway *apigateway.Client,
	config *config.AppConfig,
	route Route,
	resources *[]types.Resource) {
	splitEndpoint := strings.Split(route.Path, "/")
	lastResource := findResource(resources, "")

	for _, resource := range splitEndpoint {
		if resource == "" {
			continue
		}
		apigatewayResource := findResource(resources, resource)
		if apigatewayResource.Id != nil {
			lastResource = apigatewayResource
			continue
		}
		createResourceOutput, err := clientApiGateway.CreateResource(ctx, &apigateway.CreateResourceInput{
			ParentId:  lastResource.Id,
			RestApiId: &config.RestApiId,
			PathPart:  &resource,
		})
		if err != nil {
			log.Println("Error creating resource:", err)
			return
		}

		newResource := types.Resource{
			Id:              createResourceOutput.Id,
			PathPart:        createResourceOutput.PathPart,
			ParentId:        createResourceOutput.ParentId,
			ResourceMethods: createResourceOutput.ResourceMethods,
		}
		*resources = append(*resources, newResource)
		lastResource = &newResource
	}

	if _, exists := lastResource.ResourceMethods[route.Method]; !exists {
		log.Println("Creating method for", route.Path)
		requestParameters := utils.ExtracParamNames(route.Path)
		_, err := clientApiGateway.PutMethod(ctx, &apigateway.PutMethodInput{
			ApiKeyRequired:    true,
			RestApiId:         &config.RestApiId,
			ResourceId:        lastResource.Id,
			HttpMethod:        &route.Method,
			AuthorizationType: aws.String("NONE"),
			RequestParameters: utils.ConvertParamNamesToMappingWithPrefix(requestParameters, "method.request.path."),
		})
		if err != nil {
			log.Println("Error creating method:", err)
			return
		}

		addMethodToResourceById(*resources, *lastResource.Id, route.Method)

		var sb strings.Builder
		sb.WriteString(config.BackendUrl)
		sb.WriteString(route.Path)

		log.Println("Creating integration for", route.Path)

		_, err = clientApiGateway.PutIntegration(ctx, &apigateway.PutIntegrationInput{
			HttpMethod:            &route.Method,
			ResourceId:            lastResource.Id,
			RestApiId:             &config.RestApiId,
			Type:                  types.IntegrationTypeHttpProxy,
			Uri:                   aws.String(sb.String()),
			IntegrationHttpMethod: &route.Method,
			RequestParameters:     utils.ConvertParamNamesToMapping(requestParameters),
		})
		if err != nil {
			log.Println("Error creating integration:", err)
			return
		}
	}

	if _, exists := lastResource.ResourceMethods["OPTIONS"]; !exists && config.EnableCors {
		createCORSMethod(ctx, clientApiGateway, config, lastResource.Id)
		addMethodToResourceById(*resources, *lastResource.Id, "OPTIONS")
	}
}

func extractRoutesFromFile(filePath string, controllerRegex, methodRegex *regexp.Regexp, routes *[]Route) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file", filePath, ":", err)
		return
	}
	defer file.Close()

	var baseRoute string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if matches := controllerRegex.FindStringSubmatch(line); len(matches) > 1 {
			baseRoute = matches[1]
		}

		if matches := methodRegex.FindStringSubmatch(line); len(matches) > 2 {
			methodRoute := matches[2]
			fullRoute := filepath.Join("/", baseRoute, methodRoute)
			fullRoute = strings.ReplaceAll(fullRoute, "\\", "/")
			fullRoute = regexp.MustCompile(`:(\w+)`).ReplaceAllString(fullRoute, `{$1}`)
			*routes = append(*routes, Route{
				Method: strings.ToUpper(matches[1]),
				Path:   fullRoute,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading file", filePath, ":", err)
	}
}
