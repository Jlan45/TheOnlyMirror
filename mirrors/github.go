package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
)

func Github(w http.ResponseWriter, r *http.Request) {
	githubProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("github"))
	githubProxy.ServeHTTP(w, r)
}
