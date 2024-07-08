package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"io"
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
	pypiIndexProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("pypi_index"))
	pypiIndexProxy.ModifyResponse = func(resp *http.Response) error {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		modifiedBody := string(bodyBytes)
		modifiedBody = strings.Replace(modifiedBody, config.ServerConfig.Source["pypi_files"], config.ServerConfig.GetMyUrl().String(), -1)
		resp.Body = io.NopCloser(strings.NewReader(modifiedBody))
		resp.ContentLength = int64(len(modifiedBody))
		return nil
	}
	pypiIndexProxy.ServeHTTP(w, r)
}
func pypiFiles(w http.ResponseWriter, r *http.Request) {
	pypiFilesProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("pypi_files"))
	pypiFilesProxy.ServeHTTP(w, r)
}
