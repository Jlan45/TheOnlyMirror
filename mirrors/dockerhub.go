package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"io"
	"net/http"
)

func Dockerhub(w http.ResponseWriter, r *http.Request) {
	dockerhubProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("dockerhub"))
	dockerhubProxy.ModifyResponse = func(resp *http.Response) error {
		println("dockerhub")
		bodyBytes, _ := io.ReadAll(resp.Body)
		println(string(bodyBytes))
		return nil
	}
	dockerhubProxy.ServeHTTP(w, r)
}
