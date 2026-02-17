package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
)

func Aosp(w http.ResponseWriter, r *http.Request) {
	aospProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("aosp"))
	aospProxy.ServeHTTP(w, r)
}