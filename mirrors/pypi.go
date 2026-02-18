package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

var (
	pypiIndexProxy     *httputil.ReverseProxy
	pypiIndexProxyOnce sync.Once
	pypiFilesProxy     *httputil.ReverseProxy
	pypiFilesProxyOnce sync.Once
)

func Pypi(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/simple") {
		pypiSimple(w, r)
	} else {
		pypiFiles(w, r)
	}
}

func pypiSimple(w http.ResponseWriter, r *http.Request) {
	pypiIndexProxyOnce.Do(func() {
		pypiIndexProxy = utils.GetContentReplaceReverseProxy(
			config.ServerConfig.GetSourceUrl("pypi_index"),
			config.ServerConfig.Source["pypi_files"],
			config.ServerConfig.GetMyUrl().String(),
			"pypi_index",
		)
	})
	pypiIndexProxy.ServeHTTP(w, r)
}

func pypiFiles(w http.ResponseWriter, r *http.Request) {
	pypiFilesProxyOnce.Do(func() {
		pypiFilesProxy = utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("pypi_files"), "pypi_files")
	})
	pypiFilesProxy.ServeHTTP(w, r)
}
