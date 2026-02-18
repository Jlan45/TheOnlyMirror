package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
	"net/http/httputil"
	"sync"
)

var (
	debianProxy     *httputil.ReverseProxy
	debianProxyOnce sync.Once
)

func Debian(w http.ResponseWriter, r *http.Request) {
	debianProxyOnce.Do(func() {
		debianProxy = utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("debian"), "debian")
	})
	debianProxy.ServeHTTP(w, r)
}
