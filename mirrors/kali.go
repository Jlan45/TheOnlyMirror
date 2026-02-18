package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
	"net/http/httputil"
	"sync"
)

var (
	kaliProxy     *httputil.ReverseProxy
	kaliProxyOnce sync.Once
)

func Kali(w http.ResponseWriter, r *http.Request) {
	kaliProxyOnce.Do(func() {
		kaliProxy = utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("kali"), "kali")
	})
	kaliProxy.ServeHTTP(w, r)
}
