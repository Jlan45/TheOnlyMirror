package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
)

func Ubuntu(w http.ResponseWriter, r *http.Request) {
	ubuntuProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("ubuntu"))
	ubuntuProxy.ServeHTTP(w, r)
}
func UbuntuPorts(w http.ResponseWriter, r *http.Request) {
	ubuntuProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("ubuntu_ports"))
	ubuntuProxy.ServeHTTP(w, r)
}
