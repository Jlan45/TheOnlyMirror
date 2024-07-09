package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
	"strings"
)

func Pypi(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/simple") {
		pypiSimple(w, r)
	} else {
		pypiFiles(w, r)
	}
}
func pypiSimple(w http.ResponseWriter, r *http.Request) {
	pypiIndexProxy := utils.GetContentReplaceReverseProxy(config.ServerConfig.GetSourceUrl("pypi_index"), config.ServerConfig.Source["pypi_files"], config.ServerConfig.GetMyUrl().String())
	pypiIndexProxy.ServeHTTP(w, r)
}
func pypiFiles(w http.ResponseWriter, r *http.Request) {
	pypiFilesProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("pypi_files"))
	pypiFilesProxy.ServeHTTP(w, r)
}
