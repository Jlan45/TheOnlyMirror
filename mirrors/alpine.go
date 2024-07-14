package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
)

func Alpine(w http.ResponseWriter, r *http.Request) {
	alpineProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("alpine"))
	alpineProxy.ServeHTTP(w, r)
}
