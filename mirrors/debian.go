package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
)

func Debian(w http.ResponseWriter, r *http.Request) {
	debianProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("debian"))
	debianProxy.ServeHTTP(w, r)
}
