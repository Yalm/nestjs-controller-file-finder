package utils

import "strings"

type Route struct {
	Method    string
	Path      string
	Params    []string
	Resources []string
}

func NewRoute(method string, path string) Route {
	return Route{Method: method, Path: path, Params: ExtracParamNames(path)}
}

func (r *Route) SplitPath() []string {
	return strings.Split(r.Path, "/")
}

func (r *Route) GetUri(basePath string) string {
	var sb strings.Builder
	sb.WriteString(basePath)
	sb.WriteString(r.Path)

	return sb.String()
}
