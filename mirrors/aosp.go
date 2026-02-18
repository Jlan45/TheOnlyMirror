package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
	"net/http/httputil"
	"sync"
)

var (
	aospProxy     *httputil.ReverseProxy
	aospProxyOnce sync.Once
)

func Aosp(w http.ResponseWriter, r *http.Request) {
	aospProxyOnce.Do(func() {
		aospProxy = utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("aosp"), "aosp")
	})
	aospProxy.ServeHTTP(w, r)
}
