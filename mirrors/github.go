package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
	"net/http/httputil"
	"sync"
)

var (
	githubProxy     *httputil.ReverseProxy
	githubProxyOnce sync.Once
)

func Github(w http.ResponseWriter, r *http.Request) {
	githubProxyOnce.Do(func() {
		githubProxy = utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("github"), "github")
	})
	githubProxy.ServeHTTP(w, r)
}