package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
	"net/http/httputil"
	"sync"
)

var (
	npmProxy     *httputil.ReverseProxy
	npmProxyOnce sync.Once
)

func Npm(w http.ResponseWriter, r *http.Request) {
	npmProxyOnce.Do(func() {
		npmProxy = utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("npm"), "npm")
	})
	npmProxy.ServeHTTP(w, r)
}
