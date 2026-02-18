package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
	"net/http/httputil"
	"sync"
)

var (
	alpineProxy     *httputil.ReverseProxy
	alpineProxyOnce sync.Once
)

func Alpine(w http.ResponseWriter, r *http.Request) {
	alpineProxyOnce.Do(func() {
		alpineProxy = utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("alpine"), "alpine")
	})
	alpineProxy.ServeHTTP(w, r)
}
