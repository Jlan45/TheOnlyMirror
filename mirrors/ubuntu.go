package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
	"net/http/httputil"
	"sync"
)

var (
	ubuntuProxy      *httputil.ReverseProxy
	ubuntuProxyOnce  sync.Once
	ubuntuPProxy     *httputil.ReverseProxy
	ubuntuPProxyOnce sync.Once
)

func Ubuntu(w http.ResponseWriter, r *http.Request) {
	ubuntuProxyOnce.Do(func() {
		ubuntuProxy = utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("ubuntu"), "ubuntu")
	})
	ubuntuProxy.ServeHTTP(w, r)
}

func UbuntuPorts(w http.ResponseWriter, r *http.Request) {
	ubuntuPProxyOnce.Do(func() {
		ubuntuPProxy = utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("ubuntu_ports"), "ubuntu_ports")
	})
	ubuntuPProxy.ServeHTTP(w, r)
}