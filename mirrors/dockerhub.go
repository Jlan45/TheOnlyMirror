package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
)

func Dockerhub(w http.ResponseWriter, r *http.Request) {
	dockerhubProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("dockerhub"))
	dockerhubProxy.ServeHTTP(w, r)
}
