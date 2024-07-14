package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
)

func Npm(w http.ResponseWriter, r *http.Request) {
	npmProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("npm"))
	npmProxy.ServeHTTP(w, r)
}
